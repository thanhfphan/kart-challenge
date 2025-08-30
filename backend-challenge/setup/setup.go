package setup

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/thanhfphan/kart-challenge/config"
	"github.com/thanhfphan/kart-challenge/env"
	"github.com/thanhfphan/kart-challenge/pkg/cache"
	"github.com/thanhfphan/kart-challenge/pkg/infras/mysql"
	"github.com/thanhfphan/kart-challenge/pkg/infras/redis"
	"github.com/thanhfphan/kart-challenge/pkg/logging"

	"github.com/sethvargo/go-envconfig"
)

func LoadFromEnv(ctx context.Context) (*config.Config, *env.Env, error) {
	log := logging.FromContext(ctx)
	log.Infof("starting load config from ENV ...")

	var envOption []env.Option
	cfg := &config.Config{}
	if err := envconfig.ProcessWith(ctx, &envconfig.Config{
		Target:   cfg,
		Lookuper: envconfig.OsLookuper(),
	}); err != nil {
		return nil, nil, fmt.Errorf("envconfig.ProcessWith has err=%w", err)
	}

	// DB
	db, _, err := mysql.InitConn(cfg.DB)
	if err != nil {
		return nil, nil, fmt.Errorf("init mysql failed err=%w", err)
	}
	envOption = append(envOption, env.WithDatabase(db))

	// redis
	redisClient, err := redis.New(cfg.Redis)
	if err != nil {
		return nil, nil, fmt.Errorf("init redis failed err=%w", err)
	}
	envOption = append(envOption, env.WithRedisClient(redisClient))

	// cache
	cache := cache.New(redisClient, 24*time.Hour)
	envOption = append(envOption, env.WithCache(cache))

	// // prometheus
	// xm, err := xmetric.New()
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("init xmetric err=%w", err)
	// }
	// envOption = append(envOption, env.WithXMetric(xm))

	log.Info("setup", slog.Any("config", cfg))

	return cfg, env.NewEnv(envOption...), nil
}
