package main

import (
	"context"

	authService "sakucita/internal/app/auth/service"
	"sakucita/internal/domain"
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
	databases := databaseProvider(cfg, log)

	queries := repository.New(databases.postgres)

	security := securityProvider(cfg, log)

	services := serviceProvider(cfg, log, databases, queries, security)

	middleware := middlewareProvider(log, security, services)

	serverHttp := ServerHTTPProvider(cfg, log, services, middleware)

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
		panic(err)
	}

	return security
}

// service provider
type services struct {
	authService domain.AuthService
}

func serviceProvider(config config.App, log zerolog.Logger, databases *databases, queries *repository.Queries, security *security.Security) *services {
	return &services{
		authService: authService.NewService(databases.postgres, databases.redis, queries, config, security, log),
	}
}

// database provider
type databases struct {
	postgres *pgxpool.Pool
	redis    *redis.Client
}

func databaseProvider(cfg config.App, log zerolog.Logger) *databases {
	pg, err := postgres.NewDB(context.Background(), cfg, log)
	if err != nil {
		panic(err)
	}

	redis, err := redisClient.NewRedisClient(cfg, log)
	if err != nil {
		panic(err)
	}
	return &databases{
		postgres: pg,
		redis:    redis,
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
func ServerHTTPProvider(cfg config.App, log zerolog.Logger, services *services, middleware *middleware.Middleware) *server.Server {
	return server.NewServer(
		cfg, log, services.authService, middleware,
	)
}
