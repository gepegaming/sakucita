package service

import (
	"context"
	"encoding/json"
	"fmt"

	"sakucita/internal/database/repository"
	"sakucita/internal/domain"
	"sakucita/internal/server/security"
	"sakucita/internal/shared/utils"
	"sakucita/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type service struct {
	db       *pgxpool.Pool
	rdb      *redis.Client
	q        *repository.Queries
	config   config.App
	security *security.Security
	log      zerolog.Logger
}

func NewService(
	db *pgxpool.Pool,
	rdb *redis.Client,
	q *repository.Queries,
	config config.App,
	security *security.Security,
	log zerolog.Logger,
) domain.AuthService {
	return &service{db, rdb, q, config, security, log}
}

func (s *service) RefreshToken(ctx context.Context, req domain.RefreshRequest) (*domain.RefreshResponse, error) {
	// get session
	session, err := s.q.GetActiveSessionByTokenID(ctx, utils.SringToPgTypeUUID(req.Claims.RegisteredClaims.ID))
	if err != nil {
		if utils.IsNotFoundError(err) {
			return nil, domain.NewAppError(fiber.StatusUnauthorized, domain.ErrMsgSessionNotFound, domain.ErrUnauthorized)
		}
		s.log.Err(err).Msg("failed to get session for refresh token")
		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	// cek kalo device nya beda dari hash device_id maka error
	deviceID := utils.GenerateDeviceID(req.Claims.UserID, req.ClientInfo)
	if deviceID != session.DeviceID {
		s.log.Warn().Msgf("device id missmatch: %v != %v", deviceID, session.DeviceID)
		return nil, domain.NewAppError(fiber.StatusUnauthorized, domain.ErrMsgDeviceIdMissmatch, domain.ErrUnauthorized)
	}

	// generate token
	accessTokenID := uuid.New()
	accessToken, _, err := s.security.GenerateToken(req.Claims.UserID, accessTokenID, req.Claims.Role, s.config.JWT.AccessTokenExpiresIn)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to generate access token")
		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}
	// generate refresh token
	refreshTokenID := uuid.New()
	refreshToken, rtClaims, err := s.security.GenerateToken(req.Claims.UserID, refreshTokenID, req.Claims.Role, s.config.JWT.RefreshTokenExpiresIn)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to generate refresh token")
		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
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
		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	return &domain.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) Me(ctx context.Context, userID uuid.UUID) (*domain.UserWithRoles, error) {
	userWithRolesRow, err := s.q.GetUserByIDWithRoles(ctx, userID)
	if err != nil {
		s.log.Warn().Err(err).Msg("failed to get user with roles")
		return nil, domain.NewAppError(fiber.StatusNotFound, "user not found", domain.ErrNotfound)
	}
	// unmarshar role
	var roles []domain.Role
	if len(userWithRolesRow.Roles) > 0 {
		if err := json.Unmarshal(userWithRolesRow.Roles, &roles); err != nil {
			s.log.Error().Err(err).Msg("failed to unmarshal roles array")
			return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
		}
	}

	userResponse := &domain.UserWithRoles{
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
		Roles: roles,
	}

	return userResponse, nil
}

func (s *service) LoginLocal(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	// check ban user attemp, walau udah pake middleware tp best practice nya gini, karna kedepanya mungkin pake grpc
	ttl, err := s.CheckLoginBan(ctx, req.Email)
	if err != nil {
		return nil, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	if ttl > 0 {
		return nil, domain.NewAppError(
			fiber.StatusTooManyRequests,
			fmt.Sprintf("too many attempts. please wait after %s", ttl),
			domain.ErrTooManyRequests,
		)
	}

	// get user identity
	authIdentity, err := s.q.GetAuthIdentityByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.NewAppError(fiber.StatusNotFound, domain.ErrMsgUserNotFound, domain.ErrNotfound)
	}
	// compare password
	if !utils.CheckPassword(req.Password, authIdentity.PasswordHash.String) {
		// kalo gagal login, tambahin counter attempt
		_, _ = s.OnLoginFail(ctx, req.Email)
		return nil, domain.NewAppError(fiber.StatusUnauthorized, domain.ErrMsgInvalidCredentials, domain.ErrUnauthorized)
	}
	// get user
	userWithRolesRow, err := s.q.GetUserByIDWithRoles(ctx, authIdentity.UserID)
	if err != nil {
		s.log.Warn().Err(err).Msg("failed to get user but auth identity found")
		return nil, domain.NewAppError(fiber.StatusNotFound, domain.ErrMsgUserNotFound, domain.ErrNotfound)
	}

	// unmarshar role
	var roles []domain.Role
	if len(userWithRolesRow.Roles) > 0 {
		if err := json.Unmarshal(userWithRolesRow.Roles, &roles); err != nil {
			s.log.Error().Err(err).Msg("failed to unmarshal roles array")
			return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
		}
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
		Roles: roles,
	}

	// generate token
	accessTokenID := uuid.New()
	accessToken, _, err := s.security.GenerateToken(userResponse.ID, accessTokenID, roles, s.config.JWT.AccessTokenExpiresIn)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to generate access token")
		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	// generate refresh token
	refreshTokenID := uuid.New()
	refreshToken, rtClaims, err := s.security.GenerateToken(userResponse.ID, refreshTokenID, roles, s.config.JWT.RefreshTokenExpiresIn)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to generate refresh token")
		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	// delete all session if user activate single session
	if userResponse.SingleSession {
		err := s.q.RevokeAllSessionsByUserID(ctx, userResponse.ID)
		if err != nil {
			s.log.Error().Err(err).Msg("failed to revoke all sessions")
			return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
		}
	}

	// device id
	deviceID := utils.GenerateDeviceID(userResponse.ID, req.ClientInfo)

	// create session
	_, err = s.q.UpsertSession(ctx, repository.UpsertSessionParams{
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
		s.log.Error().Err(err).Msg("failed to create session")
		return nil, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	return &domain.LoginResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) RegisterLocal(ctx context.Context, req domain.RegisterRequest) error {
	// setup tx
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qtx := s.q.WithTx(tx)

	// create user
	user, err := qtx.CreateUser(ctx, repository.CreateUserParams{
		Email:    req.Email,
		Phone:    utils.StringToPgTypeText(req.Phone),
		Name:     req.Name,
		Nickname: req.Nickname,
	})
	if err != nil {
		if utils.IsDuplicateUniqueViolation(err) {
			switch utils.PgConstraint(err) {
			case "users_email_key":
				return domain.NewAppError(fiber.StatusConflict, domain.ErrMsgEmailAlreadyExists, domain.ErrConflict)
			case "users_phone_key":
				return domain.NewAppError(fiber.StatusConflict, domain.ErrMsgPhoneAlreadyExists, domain.ErrConflict)
			case "users_nickname_key":
				return domain.NewAppError(fiber.StatusConflict, domain.ErrMsgNicknameAlreadyExists, domain.ErrConflict)
			}
		}
		s.log.Err(err).Msg("failed to create user")
		return domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	// create user role
	err = qtx.CreateUserRole(ctx, repository.CreateUserRoleParams{
		UserID: user.ID,
		RoleID: domain.CREATOR.ID, // creator
	})
	if err != nil {
		s.log.Err(err).Msg("failed to create user role")
		return domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	// hashing password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.log.Err(err).Msg("failed to hash password")
		return domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	// create auth identity
	err = qtx.CreateAuthIdentityLocal(ctx, repository.CreateAuthIdentityLocalParams{
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
