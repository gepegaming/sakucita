package middleware

import (
	"fmt"

	"sakucita/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) LoginLimiter(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email" form:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid request body",
		)
	}

	if req.Email == "" {
		// biarin validator di handler yang tangani
		return c.Next()
	}

	// 2. Cek apakah email sedang diban
	ttl, err := m.authService.CheckLoginBan(c.Context(), req.Email)
	if err != nil {
		// Redis error → infra error
		return fiber.NewError(
			fiber.StatusInternalServerError,
			"internal server error",
		)
	}

	// 3. Kalau masih diban
	if ttl > 0 {
		return c.Status(fiber.StatusTooManyRequests).JSON(domain.ErrorResponse{
			Message: fiber.ErrTooManyRequests.Message,
			Errors:  fmt.Sprintf("too many attempts. please wait after %s", ttl.String()),
		})
	}

	return c.Next()
}
