package redis

import (
	"context"
	"fmt"
	"time"

	"sakucita/pkg/config"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func NewRedisClient(cfg config.App, log zerolog.Logger) (*redis.Client, error) {
	opt := &redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           0,
		PoolSize:     100,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  30 * time.Second,
	}

	rdb := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	log.Info().Msg("redis connected")
	return rdb, nil
}
