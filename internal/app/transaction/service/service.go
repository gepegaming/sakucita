package service

import (
	"context"
	"sakucita/internal/domain"
	"sakucita/internal/infra/postgres/repository"
	"sakucita/internal/shared/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type TransactionService interface {
	CreateWithTx(ctx context.Context, qtx repository.Querier, cmd CreateTransactionCommand) (domain.Transaction, error)
	UpdateExternalReferenceAndStatus(ctx context.Context, id uuid.UUID, externalRef string, status domain.TransactionStatus) error
	GetByExternalReference(ctx context.Context, ref string) (domain.Transaction, error)
	MarkAsSuccess(ctx context.Context, id uuid.UUID) error
	MarkAsFailed(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TransactionStatus) error
}

type transaction struct {
	q   *repository.Queries
	log zerolog.Logger
}

func NewService(q *repository.Queries, log zerolog.Logger) TransactionService {
	return &transaction{q, log}
}

func (t *transaction) CreateWithTx(ctx context.Context, qtx repository.Querier, cmd CreateTransactionCommand) (domain.Transaction, error) {
	transactionID, err := utils.GenerateUUIDV7()
	if err != nil {
		t.log.Err(err).Msg("failed to generate transaction id")
		return domain.Transaction{}, domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}

	res, err := qtx.CreateTransaction(ctx, repository.CreateTransactionParams{
		ID:                       transactionID,
		DonationMessageID:        cmd.DonationMessageID,
		PaymentChannelID:         cmd.PaymentChannelID,
		PayerUserID:              utils.UUIDPtrToPgTypeUUID(cmd.PayerUserID),
		GrossPaidAmount:          cmd.GrossPaidAmount,
		GatewayFeeFixed:          cmd.GatewayFeeFixed,
		GatewayFeePercentageBps:  cmd.GatewayFeePercentageBps,
		GatewayFeeAmount:         cmd.GatewayFeeAmount,
		PlatformFeeFixed:         cmd.PlatformFeeFixed,
		PlatformFeePercentageBps: cmd.PlatformFeePercentageBps,
		PlatformFeeAmount:        cmd.PlatformFeeAmount,
		FeeFixed:                 cmd.FeeFixed,
		FeePercentageBps:         cmd.FeePercentageBps,
		FeeAmount:                cmd.FeeAmount,
		NetAmount:                cmd.NetAmount,
		Currency:                 cmd.Currency,
		Status:                   repository.TransactionStatus(cmd.Status),
		ExternalReference:        cmd.ExternalReference,
	})
	if err != nil {
		return domain.Transaction{}, err
	}

	return transactionRepoToDomain(res), nil
}

func (t *transaction) UpdateExternalReferenceAndStatus(ctx context.Context, id uuid.UUID, externalRef string, status domain.TransactionStatus) error {
	if err := t.q.UpdateTransactionExternalReferenceAndStatus(ctx, repository.UpdateTransactionExternalReferenceAndStatusParams{
		ID:                id,
		ExternalReference: utils.StringToPgTypeText(externalRef),
		Status:            repository.TransactionStatus(status),
	}); err != nil {
		t.log.Err(err).Msg("failed to update transaction external reference")
		return domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}
	return nil
}

func (t *transaction) GetByExternalReference(ctx context.Context, ref string) (domain.Transaction, error) {
	res, err := t.q.GetTransactionByExternalReference(ctx, utils.StringToPgTypeText(ref))
	if err != nil {
		if utils.IsNotFoundError(err) {
			return domain.Transaction{}, domain.NewAppError(fiber.StatusNotFound, domain.ErrMsgTransactionNotFound, domain.ErrNotfound)
		}
	}
	return transactionRepoToDomain(res), nil
}

func (t *transaction) MarkAsFailed(ctx context.Context, id uuid.UUID) error {
	if err := t.q.UpdateTransactionStatus(ctx, repository.UpdateTransactionStatusParams{
		ID:     id,
		Status: repository.TransactionStatusFAILED,
	}); err != nil {
		t.log.Err(err).Msg("failed to mark transaction as failed")
		return domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}
	return nil
}

func (t *transaction) MarkAsSuccess(ctx context.Context, id uuid.UUID) error {
	if err := t.q.UpdateTransactionStatus(ctx, repository.UpdateTransactionStatusParams{
		ID:     id,
		Status: repository.TransactionStatusPAID,
	}); err != nil {
		t.log.Err(err).Msg("failed to mark transaction as success")
		return domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}
	return nil
}

func (t *transaction) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TransactionStatus) error {
	if err := t.q.UpdateTransactionStatus(ctx, repository.UpdateTransactionStatusParams{
		ID:     id,
		Status: repository.TransactionStatus(status),
	}); err != nil {
		t.log.Err(err).Msg("failed to update transaction status")
		return domain.NewAppError(
			fiber.StatusInternalServerError,
			domain.ErrMsgInternalServerError,
			domain.ErrInternalServerError,
		)
	}
	return nil
}
