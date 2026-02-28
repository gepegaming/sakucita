package service

import (
	"sakucita/internal/infra/postgres/repository"

	"github.com/google/uuid"
)

type CreateDonationCommand struct {
	PayeeUserID uuid.UUID
	PayerUserID *uuid.UUID

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
	Actions       []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"actions,omitempty"`
}

type CreateDonationMessageCommand struct {
	PayeeUserID       uuid.UUID
	PayerUserID       *uuid.UUID
	PayerName         string
	Email             string
	Message           string
	MediaType         repository.DonationMediaType
	MediaUrl          *string
	MediaStartSeconds *int32
	PricePerSecond    int64
	GrossPaidAmount   int64
	Amount            int64
	Currency          string
	Meta              []byte
}
