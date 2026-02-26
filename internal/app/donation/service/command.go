package service

import "github.com/google/uuid"

type CreateDonationCommand struct {
	PayeeUserID uuid.UUID
	PayerUserID uuid.UUID

	PayerName string
	Email     string
	Message   string

	MediaType string

	// Media input dari user
	MediaURL          *string
	MediaStartSeconds *int32

	// Transaction
	Amount         int64
	PaymentChannel string
}

type CreateDonationResult struct {
	TransactionID string `json:"transaction_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	QrString      string `json:"qr_string"`
}
