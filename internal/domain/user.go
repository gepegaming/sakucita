package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID      `json:"id"`
	Email         string         `json:"email"`
	EmailVerified bool           `json:"email_verified"`
	Phone         *string        `json:"phone"`
	Name          string         `json:"name"`
	Nickname      string         `json:"nickname"`
	ImageUrl      *string        `json:"image_url"`
	SingleSession bool           `json:"single_session"`
	Meta          map[string]any `json:"meta,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     time.Time      `json:"deleted_at"`
}

type UserWithRoles struct {
	User
	Roles []Role `json:"roles"`
}
