package service

import (
	"context"
	"fmt"
	userService "sakucita/internal/app/user/service"
	"sakucita/internal/domain"
	"sakucita/internal/infra/postgres/repository"
	"sakucita/internal/server/security"
	"sakucita/internal/shared/utils"
	"sakucita/pkg/config"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type service struct {
	userService userService.UserService
	db          *pgxpool.Pool
	rdb         *redis.Client
	q           *repository.Queries
	config      config.App
	security    *security.Security
	log         zerolog.Logger
}

type AuthService interface {
	RegisterLocal(ctx context.Context, req RegisterCommand) error
	LoginLocal(ctx context.Context, req LoginLocalCommand) (LoginResponse, error)
	Me(ctx context.Context, userID uuid.UUID) (domain.UserWithRoles, error)
	RefreshToken(ctx context.Context, req RefreshCommand) (RefreshResponse, error)

	// auth policy
	CheckLoginBan(ctx context.Context, id string) (time.Duration, error)
	OnLoginFail(ctx context.Context, id string) (time.Duration, error)
	OnLoginSuccess(ctx context.Context, id string) error
}

func NewService(
	userService userService.UserService,
	db *pgxpool.Pool,
	rdb *redis.Client,
	q *repository.Queries,
	config config.App,
	security *security.Security,
	log zerolog.Logger,
) AuthService {
	return &service{userService, db, rdb, q, config, security, log}
}

func (s *service) RefreshToken(ctx context.Context, req RefreshCommand) (RefreshResponse, error) {
	// get session
	session, err := s.q.GetActiveSessionByTokenID(ctx, utils.StringToPgTypeUUID(req.Claims.RegisteredClaims.ID))
	if err != nil {
		if utils.IsNotFoundError(err) {
			return RefreshResponse{}, domain.NewAppError(fiber.StatusUnauthorized, domain.ErrMsgSessionNotFound, domain.ErrUnauthorized)
		}
		s.log.Err(err).Msg("failed to get session for refresh token")
		return RefreshResponse{}, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	// cek kalo device nya beda dari hash device_id maka error
	deviceID := security.GenerateDeviceID(req.Claims.UserID, req.ClientInfo)
	if deviceID != session.DeviceID {
		s.log.Warn().
			Str("device_id", deviceID).
			Str("session_device_id", session.DeviceID).
			Msg("device id mismatch")
		return RefreshResponse{}, domain.NewAppError(fiber.StatusUnauthorized, domain.ErrMsgDeviceIdMissmatch, domain.ErrUnauthorized)
	}

	// generate token
	accessTokenID := uuid.New()
	accessToken, _, err := s.security.GenerateToken(req.Claims.UserID, accessTokenID, req.Claims.Role, s.config.JWT.AccessTokenExpiresIn)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to generate access token")
		return RefreshResponse{}, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}
	// generate refresh token
	refreshTokenID := uuid.New()
	refreshToken, rtClaims, err := s.security.GenerateToken(req.Claims.UserID, refreshTokenID, req.Claims.Role, s.config.JWT.RefreshTokenExpiresIn)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to generate refresh token")
		return RefreshResponse{}, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	// create session
	_, err = s.q.UpsertSession(ctx, repository.UpsertSessionParams{
		UserID:   req.Claims.UserID,
		DeviceID: deviceID,
		RefreshTokenID: pgtype.UUID{
			Bytes: refreshTokenID,
			Valid: true,
		},
		ExpiresAt: pgtype.Timestamptz{
			Time:  rtClaims.ExpiresAt.Time,
			Valid: true,
		},
		Meta: map[string]any{
			"ip":          req.ClientInfo.IP,
			"user_agent":  req.ClientInfo.UserAgent,
			"device_name": req.ClientInfo.DeviceName,
		},
	})
	if err != nil {
		s.log.Error().Err(err).Msg("failed to create session")
		return RefreshResponse{}, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	return RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) Me(ctx context.Context, userID uuid.UUID) (domain.UserWithRoles, error) {
	userWithRolesRow, err := s.userService.GetByIDWithRoles(ctx, userID)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get user with roles")
		return domain.UserWithRoles{}, domain.NewAppError(fiber.StatusNotFound, "user not found", domain.ErrNotfound)
	}

	// userResponse := &domain.UserWithRoles{
	// 	User: domain.User{
	// 		ID:            userWithRolesRow.ID,
	// 		Email:         userWithRolesRow.Email,
	// 		EmailVerified: userWithRolesRow.EmailVerified,
	// 		Phone:         userWithRolesRow.Phone,
	// 		Name:          userWithRolesRow.Name,
	// 		Nickname:      userWithRolesRow.Nickname,
	// 		ImageUrl:      userWithRolesRow.ImageUrl,
	// 		SingleSession: userWithRolesRow.SingleSession,
	// 		Meta:          userWithRolesRow.Meta,
	// 		CreatedAt:     userWithRolesRow.CreatedAt,
	// 		UpdatedAt:     userWithRolesRow.UpdatedAt,
	// 		DeletedAt:     userWithRolesRow.DeletedAt,
	// 	},
	// 	Roles: userWithRolesRow.Roles,
	// }

	return userWithRolesRow, nil
}

func (s *service) LoginLocal(ctx context.Context, req LoginLocalCommand) (LoginResponse, error) {
	// check ban user attemp, walau udah pake middleware tp best practice nya gini, karna kedepanya mungkin pake grpc
	ttl, err := s.CheckLoginBan(ctx, req.Email)
	if err != nil {
		return LoginResponse{}, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	if ttl > 0 {
		return LoginResponse{}, domain.NewAppError(
			fiber.StatusTooManyRequests,
			fmt.Sprintf("too many attempts. please wait after %s", ttl),
			domain.ErrTooManyRequests,
		)
	}

	// get user identity
	authIdentity, err := s.q.GetAuthIdentityByEmail(ctx, req.Email)
	if err != nil {
		s.log.Error().Err(err).Msg("user not found ges")
		return LoginResponse{}, domain.NewAppError(fiber.StatusNotFound, domain.ErrMsgUserNotFound, domain.ErrNotfound)
	}
	// compare password
	if !utils.CheckPassword(req.Password, authIdentity.PasswordHash.String) {
		// kalo gagal login, tambahin counter attempt
		_, _ = s.OnLoginFail(ctx, req.Email)
		return LoginResponse{}, domain.NewAppError(fiber.StatusUnauthorized, domain.ErrMsgInvalidCredentials, domain.ErrUnauthorized)
	}
	// get user
	userWithRolesRow, err := s.userService.GetByIDWithRoles(ctx, authIdentity.UserID)
	if err != nil {
		s.log.Warn().Err(err).Msg("failed to get user but auth identity found")
		return LoginResponse{}, domain.NewAppError(fiber.StatusNotFound, domain.ErrMsgUserNotFound, domain.ErrNotfound)
	}

	userResponse := domain.UserWithRoles{
		User: domain.User{
			ID:            userWithRolesRow.ID,
			Email:         userWithRolesRow.Email,
			EmailVerified: userWithRolesRow.EmailVerified,
			Phone:         userWithRolesRow.Phone,
			Name:          userWithRolesRow.Name,
			Nickname:      userWithRolesRow.Nickname,
			ImageUrl:      userWithRolesRow.ImageUrl,
			SingleSession: userWithRolesRow.SingleSession,
			Meta:          userWithRolesRow.Meta,
			CreatedAt:     userWithRolesRow.CreatedAt,
			UpdatedAt:     userWithRolesRow.UpdatedAt,
			DeletedAt:     userWithRolesRow.DeletedAt,
		},
		Roles: userWithRolesRow.Roles,
	}

	// generate token
	accessTokenID := uuid.New()
	accessToken, _, err := s.security.GenerateToken(userResponse.ID, accessTokenID, userResponse.Roles, s.config.JWT.AccessTokenExpiresIn)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to generate access token")
		return LoginResponse{}, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	// generate refresh token
	refreshTokenID := uuid.New()
	refreshToken, rtClaims, err := s.security.GenerateToken(userResponse.ID, refreshTokenID, userResponse.Roles, s.config.JWT.RefreshTokenExpiresIn)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to generate refresh token")
		return LoginResponse{}, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	// delete all session if user activate single session
	if userResponse.SingleSession {
		err := s.q.RevokeAllSessionsByUserID(ctx, userResponse.ID)
		if err != nil {
			s.log.Error().Err(err).Msg("failed to revoke all sessions")
			return LoginResponse{}, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
		}
	}

	// device id
	deviceID := security.GenerateDeviceID(userResponse.ID, security.ClientInfo{
		IP:         req.ClientInfo.IP,
		UserAgent:  req.ClientInfo.UserAgent,
		DeviceName: req.ClientInfo.DeviceName,
	})

	// create session
	sessionID, err := utils.GenerateUUIDV7()
	if err != nil {
		s.log.Err(err).Msg("failed to generate uuid v7 for sessionID")
		return LoginResponse{}, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}
	_, err = s.q.UpsertSession(ctx, repository.UpsertSessionParams{
		ID:       sessionID,
		UserID:   userResponse.ID,
		DeviceID: deviceID,
		RefreshTokenID: pgtype.UUID{
			Bytes: refreshTokenID,
			Valid: true,
		},
		ExpiresAt: pgtype.Timestamptz{
			Time:  rtClaims.ExpiresAt.Time,
			Valid: true,
		},
		Meta: map[string]any{
			"ip":          req.ClientInfo.IP,
			"user_agent":  req.ClientInfo.UserAgent,
			"device_name": req.ClientInfo.DeviceName,
		},
	})
	if err != nil {
		s.log.Error().
			Err(err).
			Str("user_id", userResponse.ID.String()).
			Str("device_id", deviceID).
			Msg("failed to create session")
		return LoginResponse{}, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	return LoginResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) RegisterLocal(ctx context.Context, req RegisterCommand) error {
	// setup tx
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qtx := s.q.WithTx(tx)

	// create user with tx
	user, err := s.userService.CreateUserWithTx(ctx, qtx, req.Email, req.Phone, req.Name, req.Nickname)
	if err != nil {
		return err
	}

	// create user role with tx
	err = qtx.CreateUserRole(ctx, repository.CreateUserRoleParams{
		UserID: user.ID,
		RoleID: domain.CREATOR.ID, // creator
	})
	if err != nil {
		return domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	// hashing password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.log.Err(err).Msg("failed to hash password")
		return domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	// create auth identity
	authIdentityID, err := utils.GenerateUUIDV7()
	if err != nil {
		s.log.Err(err).Msg("failed to generate uuid v7 for auth identity")
		return domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}
	err = qtx.CreateAuthIdentityLocal(ctx, repository.CreateAuthIdentityLocalParams{
		ID:           authIdentityID,
		UserID:       user.ID,
		Provider:     domain.PROVIDERLOCAL,
		ProviderID:   req.Email,
		PasswordHash: utils.StringToPgTypeText(hashedPassword),
	})
	if err != nil {
		s.log.Err(err).Msg("failed to create auth identity")
		return domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	// success
	if err := tx.Commit(ctx); err != nil {
		s.log.Err(err).Msg("failed to commit transaction")
		return domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}
	return nil
}
