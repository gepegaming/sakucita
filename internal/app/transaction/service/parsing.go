package service

import (
	"sakucita/internal/domain"
	"sakucita/internal/infra/postgres/repository"
	"sakucita/internal/shared/utils"
)

func transactionRepoToDomain(t repository.Transaction) domain.Transaction {
	return domain.Transaction{
		ID:                       t.ID,
		DonationMessageID:        t.DonationMessageID,
		PaymentChannelID:         t.PaymentChannelID,
		PayerUserID:              utils.PgTypeUUIDToUUIDPtr(t.PayerUserID),
		PayeeUserID:              t.PayeeUserID,
		GrossPaidAmount:          t.GrossPaidAmount,
		GatewayFeeFixed:          t.GatewayFeeFixed,
		GatewayFeePercentageBPS:  int32(t.GatewayFeePercentageBps),
		GatewayFeeAmount:         t.GatewayFeeAmount,
		PlatformFeeFixed:         t.PlatformFeeFixed,
		PlatformFeePercentageBPS: int32(t.PlatformFeePercentageBps),
		PlatformFeeAmount:        t.PlatformFeeAmount,
		FeeFixed:                 t.FeeFixed,
		FeePercentageBPS:         int32(t.FeePercentageBps),
		FeeAmount:                t.FeeAmount,
		NetAmount:                t.NetAmount,
		Currency:                 t.Currency,
		Status:                   domain.TransactionStatus(t.Status),
		CreatedAt:                t.CreatedAt.Time,
	}
}
