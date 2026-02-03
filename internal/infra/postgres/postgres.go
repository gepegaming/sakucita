package postgres

import (
	"context"
	"fmt"

	"sakucita/pkg/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func NewDB(ctx context.Context, config config.App, log zerolog.Logger) (DB *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.Database)
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// pool settings
	poolConfig.MaxConns = config.Database.MaxOpenConns
	poolConfig.MinConns = config.Database.MinOpenConns
	poolConfig.MaxConnLifetime = config.Database.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = config.Database.ConnMaxIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	log.Info().Msg("postgres connected")
	return pool, nil
}
