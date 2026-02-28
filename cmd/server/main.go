package main

import (
	"context"

	authService "sakucita/internal/app/auth/service"
	donationService "sakucita/internal/app/donation/service"
	feeService "sakucita/internal/app/fee/service"
	paymentChannelService "sakucita/internal/app/paymentChannel/service"
	transactionService "sakucita/internal/app/transaction/service"
	userService "sakucita/internal/app/user/service"
	"sakucita/internal/infra/midtrans"
	"sakucita/internal/infra/postgres"
	"sakucita/internal/infra/postgres/repository"
	redisClient "sakucita/internal/infra/redis"
	"sakucita/internal/server"
	"sakucita/internal/server/middleware"
	"sakucita/internal/server/security"
	"sakucita/pkg/config"
	"sakucita/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func main() {
	cfg := configProvider()
	log := loggerProvider(cfg)
	infras := infrasProvider(cfg, log)

	queries := repository.New(infras.postgres)

	security := securityProvider(cfg, log)

	services := serviceProvider(cfg, log, infras, queries, security)

	middleware := middlewareProvider(log, security, services)

	serverHttp := ServerHTTPProvider(cfg, log, services, middleware, infras)

	serverHttp.Start()
}

// middleware provider
func middlewareProvider(log zerolog.Logger, security *security.Security, serservices *services) *middleware.Middleware {
	return middleware.NewMiddleware(log, security, serservices.authService)
}

// security provider
func securityProvider(cfg config.App, log zerolog.Logger) *security.Security {
	security := security.NewSecurity(cfg, log)
	if err := security.LoadRSAKeys(cfg.JWT.KeyDirPath); err != nil {
		log.Error().Err(err).Msg("failed to load RSA")
		panic(err)
	}

	return security
}

// service provider
type services struct {
	authService           authService.AuthService
	donationService       donationService.DonationService
	userService           userService.UserService
	feeService            feeService.FeeService
	paymentChannelService paymentChannelService.PaymentChannelService
	transactionService    transactionService.TransactionService
}

func serviceProvider(config config.App, log zerolog.Logger, infras *infras, queries *repository.Queries, security *security.Security) *services {
	userService := userService.NewService(infras.postgres, infras.redis, queries, config, log)
	authService := authService.NewService(userService, infras.postgres, infras.redis, queries, config, security, log)
	feeService := feeService.NewService(queries, log)
	paymentChannelService := paymentChannelService.NewService(queries, log)
	transactionService := transactionService.NewService(queries, log)
	donationService := donationService.NewService(infras.postgres, queries, log, infras.midtransClient, transactionService, feeService, userService, paymentChannelService)

	return &services{
		authService:           authService,
		userService:           userService,
		donationService:       donationService,
		feeService:            feeService,
		paymentChannelService: paymentChannelService,
		transactionService:    transactionService,
	}
}

// infras provider
type infras struct {
	postgres       *pgxpool.Pool
	redis          *redis.Client
	midtransClient midtrans.MidtransClient
}

func infrasProvider(cfg config.App, log zerolog.Logger) *infras {
	pg, err := postgres.NewDB(context.Background(), cfg, log)
	if err != nil {
		panic(err)
	}

	redis, err := redisClient.NewRedisClient(cfg, log)
	if err != nil {
		panic(err)
	}

	midtransClient := midtrans.NewMidtransClient(cfg, log)

	return &infras{
		postgres:       pg,
		redis:          redis,
		midtransClient: midtransClient,
	}
}

// config provider
func configProvider() config.App {
	cfg, err := config.New("./config.yaml")
	if err != nil {
		panic(err)
	}
	return cfg
}

// logger provider
func loggerProvider(cfg config.App) zerolog.Logger {
	return logger.New("sakucita", cfg)
}

// server provider
func ServerHTTPProvider(cfg config.App, log zerolog.Logger, services *services, middleware *middleware.Middleware, infras *infras) *server.Server {
	return server.NewServer(
		cfg,
		log,
		services.authService,
		services.donationService,
		middleware,
	)
}
