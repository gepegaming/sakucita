package domain

import (
	"time"

	"github.com/google/uuid"
)

type DonationMessage struct {
	ID          uuid.UUID `json:"id"`
	PayeeUserID uuid.UUID `json:"payee_user_id"`

	PayerUserID *uuid.UUID `json:"payer_user_id,omitempty"`

	PayerName string `json:"payer_name"`
	Email     string `json:"email"`
	Message   string `json:"message"`

	MediaType string `json:"media_type"`

	MediaUrl          *string `json:"media_url,omitempty"`
	MediaStartSeconds *int32  `json:"media_start_seconds,omitempty"`
	MaxPlaySeconds    *int32  `json:"max_play_seconds,omitempty"`
	PricePerSecond    *int64  `json:"price_per_second,omitempty"`

	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`

	Status string `json:"status"`

	Meta map[string]any `json:"meta"`

	PlayedAt  time.Time `json:"played_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
