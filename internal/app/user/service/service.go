package service

import (
	"context"
	"sakucita/internal/domain"
	"sakucita/internal/infra/postgres/repository"
	"sakucita/internal/shared/utils"
	"sakucita/pkg/config"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type UserService interface {
	CreateUserWithTx(ctx context.Context, qtx repository.Querier, email, phone, name, nickname string) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByIDWithRoles(ctx context.Context, id uuid.UUID) (domain.UserWithRoles, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)

	// user role
	AddCreatorRoleToUserWithTx(ctx context.Context, qtx repository.Querier, userID uuid.UUID) error
}

type userService struct {
	db     *pgxpool.Pool
	rdb    *redis.Client
	q      *repository.Queries
	config config.App
	log    zerolog.Logger
}

func NewService(db *pgxpool.Pool, rdb *redis.Client, q *repository.Queries, config config.App, log zerolog.Logger) UserService {
	return &userService{db, rdb, q, config, log}
}

func (s *userService) AddCreatorRoleToUserWithTx(ctx context.Context, qtx repository.Querier, userID uuid.UUID) error {
	err := qtx.CreateUserRole(ctx, repository.CreateUserRoleParams{
		UserID: userID,
		RoleID: domain.CREATOR.ID, // creator
	})
	if err != nil {
		s.log.Err(err).Msg("failed to create user role")
		return domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}
	return nil
}

func (s *userService) CreateUserWithTx(ctx context.Context, qtx repository.Querier, email, phone, name, nickname string) (domain.User, error) {
	id, err := utils.GenerateUUIDV7()
	if err != nil {
		s.log.Err(err).Msg("failed to generate uuid v7 for user")
		return domain.User{}, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	user, err := qtx.CreateUser(ctx, repository.CreateUserParams{
		ID:       id,
		Email:    email,
		Phone:    utils.StringToPgTypeText(phone),
		Name:     name,
		Nickname: nickname,
	})
	if err != nil {
		if utils.IsDuplicateUniqueViolation(err) {
			switch utils.PgConstraint(err) {
			case "users_email_key":
				return domain.User{}, domain.NewAppError(fiber.StatusConflict, domain.ErrMsgEmailAlreadyExists, domain.ErrConflict)
			case "users_phone_key":
				return domain.User{}, domain.NewAppError(fiber.StatusConflict, domain.ErrMsgPhoneAlreadyExists, domain.ErrConflict)
			case "users_nickname_key":
				return domain.User{}, domain.NewAppError(fiber.StatusConflict, domain.ErrMsgNicknameAlreadyExists, domain.ErrConflict)
			}
		}
		s.log.Err(err).Msg("failed to create user")
		return domain.User{}, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}

	return userRepoToUserDomain(user), nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.q.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return userRepoToUserDomain(user), nil
}

func (s *userService) GetByIDWithRoles(ctx context.Context, id uuid.UUID) (domain.UserWithRoles, error) {
	userWithRoles, err := s.q.GetUserByIDWithRoles(ctx, id)
	if err != nil {
		return domain.UserWithRoles{}, err
	}
	return userRepoWithRolesToUserDomain(userWithRoles)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return userRepoToUserDomain(user), nil
}
