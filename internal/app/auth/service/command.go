package service

import (
	"sakucita/internal/domain"
	"sakucita/internal/server/security"
)

type RegisterCommand struct {
	Email    string
	Phone    string
	Name     string
	Nickname string
	Password string
}

type LoginLocalCommand struct {
	Email    string
	Password string

	ClientInfo security.ClientInfo
}

type LoginResponse struct {
	User         domain.UserWithRoles `json:"user"`
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
}

type RefreshCommand struct {
	Claims     security.TokenClaims
	ClientInfo security.ClientInfo
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
