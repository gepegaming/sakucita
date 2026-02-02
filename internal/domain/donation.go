package domain

import (
	"time"

	"github.com/google/uuid"
)

type DonationService interface {
	CreateDonation()
}

type DonationMessage struct {
	ID                   uuid.UUID              `json:"id"`
	PayeeUserID          uuid.UUID              `json:"payee_user_id"`
	PayerUserID          *string                `json:"payer_user_id"`
	PayerName            string                 `json:"payer_name"`
	Message              string                 `json:"message"`
	Email                string                 `json:"email"`
	MediaType            string                 `json:"media_type"`
	TTSLanguage          *string                `json:"tts_language"`
	TTSVoice             *string                `json:"tts_voice"`
	MediaProvider        *string                `json:"media_provider"`
	MediaVideoID         *string                `json:"media_video_id"`
	MediaStartSeconds    *int                   `json:"media_start_seconds"`
	MediaEndSeconds      *int                   `json:"media_end_seconds"`
	MediaDurationSeconds *int                   `json:"media_duration_seconds"`
	PlayedAt             time.Time              `json:"played_at"`
	Status               string                 `json:"status"`
	Meta                 map[string]interface{} `json:"meta"`
	CreatedAt            time.Time              `json:"created_at"`
}

type CreateDonationMessageRequest struct {
	PayeeUserID string `json:"payee_user_id" form:"payee_user_id" validate:"required,uuid"`

	PayerUserID *string `json:"payer_user_id,omitempty" form:"payer_user_id" validate:"omitempty,uuid"`

	PayerName string  `json:"payer_name" form:"payer_name" validate:"required,min=1,max=20"`
	Email     *string `json:"email,omitempty" form:"email" validate:"required,email"`
	Message   *string `json:"message,omitempty" form:"message" validate:"omitempty,min=1,max=300"`

	MediaType string `json:"media_type" form:"media_type" validate:"required,oneof=TEXT YOUTUBE GIF"`

	// media input dari user
	MediaURL          *string `json:"media_url,omitempty" form:"media_url" validate:"omitempty,url"`
	MediaStartSeconds int     `json:"media_start_seconds,omitempty" form:"media_start_seconds" validate:"min=0"`
}
