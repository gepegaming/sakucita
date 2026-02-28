package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	TransactionStatusInitial  TransactionStatus = "INITIAL"
	TransactionStatusPending  TransactionStatus = "PENDING"
	TransactionStatusPaid     TransactionStatus = "PAID"
	TransactionStatusFailed   TransactionStatus = "FAILED"
	TransactionStatusExpired  TransactionStatus = "EXPIRED"
	TransactionStatusRefunded TransactionStatus = "REFUNDED"
)

type Transaction struct {
	ID uuid.UUID `json:"id"`

	DonationMessageID uuid.UUID `json:"donation_message_id"`
	PaymentChannelID  int32     `json:"payment_channel_id"`

	PayerUserID *uuid.UUID `json:"payer_user_id,omitempty"`
	PayeeUserID uuid.UUID  `json:"payee_user_id"`

	// Gross
	GrossPaidAmount int64  `json:"gross_paid_amount"`
	Currency        string `json:"currency"`

	// ------------------------------
	// Gateway Snapshot
	// ------------------------------
	GatewayFeeFixed         int64 `json:"gateway_fee_fixed"`
	GatewayFeePercentageBPS int32 `json:"gateway_fee_percentage_bps"`
	GatewayFeeAmount        int64 `json:"gateway_fee_amount"`

	// ------------------------------
	// Platform Snapshot
	// ------------------------------
	PlatformFeeFixed         int64 `json:"platform_fee_fixed"`
	PlatformFeePercentageBPS int32 `json:"platform_fee_percentage_bps"`
	PlatformFeeAmount        int64 `json:"platform_fee_amount"`

	// ------------------------------
	// Total Fee
	// ------------------------------
	FeeFixed         int64 `json:"fee_fixed"`
	FeePercentageBPS int32 `json:"fee_percentage_bps"`
	FeeAmount        int64 `json:"fee_amount"`

	// Net
	NetAmount int64 `json:"net_amount"`

	Status TransactionStatus `json:"status"`

	ExternalReference *string        `json:"external_reference,omitempty"`
	Meta              map[string]any `json:"meta"`

	CreatedAt time.Time  `json:"created_at"`
	PaidAt    *time.Time `json:"paid_at,omitempty"`
	SettledAt *time.Time `json:"settled_at,omitempty"`
}

func (s TransactionStatus) IsFinal() bool {
	return s == TransactionStatusPaid ||
		s == TransactionStatusFailed ||
		s == TransactionStatusExpired ||
		s == TransactionStatusRefunded
}
