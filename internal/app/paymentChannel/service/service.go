package service

import (
	"context"
	"sakucita/internal/domain"
	"sakucita/internal/infra/postgres/repository"
	"sakucita/internal/shared/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

type PaymentChannelService interface {
	GetPaymentChannelByCodeWithTx(ctx context.Context, qtx repository.Querier, code string) (domain.PaymentChannel, error)
	GetPaymentChannelByCode(ctx context.Context, code string) (domain.PaymentChannel, error)
}

type service struct {
	q   *repository.Queries
	log zerolog.Logger
}

func NewService(q *repository.Queries, log zerolog.Logger) PaymentChannelService {
	return &service{q, log}
}

func (s *service) GetPaymentChannelByCodeWithTx(ctx context.Context, qtx repository.Querier, code string) (domain.PaymentChannel, error) {
	paymentChannel, err := qtx.GetPaymentChannelByCode(ctx, code)
	if err != nil {
		s.log.Err(err).Msg("failed to get payment channel by code")
		return domain.PaymentChannel{}, err
	}

	return paymentChannelRepoToDomain(paymentChannel), nil
}

func (s *service) GetPaymentChannelByCode(ctx context.Context, code string) (domain.PaymentChannel, error) {
	paymentChannel, err := s.q.GetPaymentChannelByCode(ctx, code)
	if err != nil {
		if utils.IsNotFoundError(err) {
			return domain.PaymentChannel{}, domain.NewAppError(
				fiber.StatusNotFound,
				domain.ErrMsgPaymentChannelNotFound,
				domain.ErrNotfound,
			)
		}
		s.log.Err(err).Msg("failed to get payment channel")
		return domain.PaymentChannel{}, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	return paymentChannelRepoToDomain(paymentChannel), nil
}
