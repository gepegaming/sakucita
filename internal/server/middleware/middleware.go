package middleware

import (
	"sakucita/internal/domain"
	"sakucita/internal/server/security"

	"github.com/rs/zerolog"
)

type Middleware struct {
	log         zerolog.Logger
	security    *security.Security
	authService domain.AuthService
}

func NewMiddleware(log zerolog.Logger, security *security.Security, authService domain.AuthService) *Middleware {
	return &Middleware{log, security, authService}
}
