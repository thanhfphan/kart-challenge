package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/thanhfphan/kart-challenge/config"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

// New create a redis from config
func New(cfg *config.Redis) (*redis.Client, error) {

	opts, err := redis.ParseURL(cfg.ConnectionURL)
	if err != nil {
		return nil, fmt.Errorf("parseURl failed err=%w", err)
	}

	opts.PoolSize = cfg.PoolSize
	opts.DialTimeout = time.Duration(cfg.DialTimeoutSeconds) * time.Second
	opts.ReadTimeout = time.Duration(cfg.ReadTimeoutSeconds) * time.Second
	opts.WriteTimeout = time.Duration(cfg.WriteTimeoutSeconds) * time.Second

	redisClient := redis.NewClient(opts)

	cmd := redisClient.Ping(context.Background())
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		return nil, fmt.Errorf("instrument tracing redis got err=%w", err)
	}
	if err := redisotel.InstrumentMetrics(redisClient); err != nil {
		return nil, fmt.Errorf("instrument metrics redis got err=%w", err)
	}

	return redisClient, nil
}
