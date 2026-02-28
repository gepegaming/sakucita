package service

import (
	"context"
	"sakucita/internal/domain"
	"sakucita/internal/infra/postgres/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type FeeService interface {
	GetUserFees(ctx context.Context, userId uuid.UUID, paymentChannelId int32) (domain.UserFee, error)
	GetUserFeesWithTx(ctx context.Context, tx repository.Querier, userId uuid.UUID, paymentChannelId int32) (domain.UserFee, error)
}

type service struct {
	q   *repository.Queries
	log zerolog.Logger
}

func NewService(q *repository.Queries, log zerolog.Logger) FeeService {
	return &service{q, log}
}

func (s *service) GetUserFees(ctx context.Context, userId uuid.UUID, paymentChannelId int32) (domain.UserFee, error) {
	res, err := s.q.GetUserFee(ctx, repository.GetUserFeeParams{Userid: userId, Paymentchannelid: paymentChannelId})
	if err != nil {
		s.log.Err(err).Msg("failed to get user fees")
		return domain.UserFee{}, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}
	return userFeeRowRepoToDomain(res), nil
}

func (s *service) GetUserFeesWithTx(ctx context.Context, tx repository.Querier, userId uuid.UUID, paymentChannelId int32) (domain.UserFee, error) {
	res, err := tx.GetUserFee(ctx, repository.GetUserFeeParams{Userid: userId, Paymentchannelid: paymentChannelId})
	if err != nil {
		s.log.Err(err).Msg("failed to get user fees with tx")
		return domain.UserFee{}, domain.NewAppError(fiber.StatusInternalServerError, domain.ErrMsgInternalServerError, domain.ErrInternalServerError)
	}
	return userFeeRowRepoToDomain(res), nil
}
