package http

import (
	"sakucita/internal/domain"
	"sakucita/internal/server/middleware"
	"sakucita/pkg/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Handler struct {
	config    config.App
	log       zerolog.Logger
	validator *validator.Validate
	mw        *middleware.Middleware
}

func NewHandler(
	config config.App,
	log zerolog.Logger,
	validator *validator.Validate,
	mw *middleware.Middleware,
) *Handler {
	return &Handler{
		config,
		log,
		validator,
		mw,
	}
}

func (h *Handler) Routes(r fiber.Router) {
	r.Route("/donations", func(router fiber.Router) {
		router.Post("/", h.CreateDonation)
	})
}

func (h *Handler) CreateDonation(c *fiber.Ctx) error {
	var req domain.CreateDonationMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := h.validator.Struct(req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(domain.Response{
		Message: "success",
		Data:    req,
	})
}
