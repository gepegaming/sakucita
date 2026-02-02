package server

import (
	"context"
	"fmt"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	authHandlerHTTP "sakucita/internal/app/auth/delivery/http"
	donationHandler "sakucita/internal/app/donation/delivery/http"
	"sakucita/internal/domain"
	"sakucita/internal/server/middleware"
	"sakucita/pkg/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
)

type Server struct {
	app       *fiber.App
	log       zerolog.Logger
	config    config.App
	validator *validator.Validate
	handlers
}

type handlers struct {
	authHandler     *authHandlerHTTP.Handler
	donationHandler *donationHandler.Handler
}

func NewServer(
	config config.App,
	log zerolog.Logger,
	authService domain.AuthService,
	middleware *middleware.Middleware,
) *Server {
	// setup validator
	validator := validator.New()
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// ambil tag dari json, jika ga ada dari form, jika ga ada dari field default
		name := fld.Tag.Get("json")
		if name == "" {
			name = fld.Tag.Get("form")
		}
		if name == "" {
			name = fld.Name
		}
		return strings.ToLower(name)
	})

	// setup fiber
	app := fiber.New(fiber.Config{
		AppName:      "Sakucita",
		ErrorHandler: fiberErrorHandler,
	})

	// setup handler
	authHandler := authHandlerHTTP.NewHandler(config, log, validator, authService, middleware)
	donationHandler := donationHandler.NewHandler(config, log, validator, middleware)
	return &Server{
		app:       app,
		log:       log,
		config:    config,
		validator: validator,
		handlers: handlers{
			authHandler:     authHandler,
			donationHandler: donationHandler,
		},
	}
}

func (s *Server) Start() {
	s.setupGlobalMiddlewares()
	s.mountRoutes()

	// start fasthttp server
	// run server di go routine
	go func() {
		addr := fmt.Sprintf("%s:%s", s.config.Server.Addr, s.config.Server.Port)
		if err := s.app.Listen(addr); err != nil {
			s.log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(shutdownCtx); err != nil {
		s.log.Error().Err(err).Msg("failed to shutdown server")
	}

	s.log.Info().Msg("server shutdown gracefully")
}

func (s *Server) setupGlobalMiddlewares() {
	s.app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &s.log,
	}))
	s.app.Use(recover.New(recover.ConfigDefault))
}

func (s *Server) mountRoutes() {
	apiv1 := s.app.Group("/api/v1")

	s.authHandler.Routes(apiv1)
	s.donationHandler.Routes(apiv1)
}
