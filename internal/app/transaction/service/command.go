package service

import (
	"sakucita/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateTransactionCommand struct {
	ID                       uuid.UUID
	DonationMessageID        uuid.UUID
	PaymentChannelID         int32
	PayerUserID              *uuid.UUID
	PayeeUserID              uuid.UUID
	GrossPaidAmount          int64
	GatewayFeeFixed          int64
	GatewayFeePercentageBps  int32
	GatewayFeeAmount         int64
	PlatformFeeFixed         int64
	PlatformFeePercentageBps int32
	PlatformFeeAmount        int64
	FeeFixed                 int64
	FeePercentageBps         int32
	FeeAmount                int64
	NetAmount                int64
	Currency                 string
	Status                   domain.TransactionStatus
	ExternalReference        pgtype.Text
}
