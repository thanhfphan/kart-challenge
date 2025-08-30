package env

import (
	"context"

	"github.com/thanhfphan/kart-challenge/pkg/cache"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Env represent latest environment configuration in this application
type Env struct {
	database    *gorm.DB
	redisClient *redis.Client
	cache       cache.Cache
}

type Option func(*Env) *Env

func NewEnv(opts ...Option) *Env {
	env := &Env{}
	for _, opt := range opts {
		env = opt(env)
	}
	return env
}

func (e *Env) Database() *gorm.DB {
	return e.database
}

func (e *Env) RedisClient() *redis.Client {
	return e.redisClient
}

func (e *Env) Cache() cache.Cache {
	return e.cache
}

func (e *Env) Close(ctx context.Context) {
	if e.database != nil {
		ins, _ := e.database.DB()
		ins.Close() //nolint
	}

	if e.redisClient != nil {
		e.redisClient.Close()
	}
}

func WithDatabase(db *gorm.DB) Option {
	return func(env *Env) *Env {
		env.database = db
		return env
	}
}

func WithRedisClient(r *redis.Client) Option {
	return func(env *Env) *Env {
		env.redisClient = r
		return env
	}
}

func WithCache(c cache.Cache) Option {
	return func(env *Env) *Env {
		env.cache = c
		return env
	}
}
