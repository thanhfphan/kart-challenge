package contxt

import (
	"context"
	"fmt"
	"sync"

	"github.com/thanhfphan/kart-challenge/pkg/utils"
)

type ctxKey int

const (
	AppCtxKey ctxKey = iota
)

type AppContext struct {
	parrentCtx context.Context
	mu         sync.RWMutex

	// Keys is a key/value pair exclusively for the context of each request.
	keys map[string]any
}

func GetAppWrapper(ctx context.Context) (*AppContext, error) {
	value := ctx.Value(AppCtxKey)
	if value == nil {
		return nil, fmt.Errorf("could not get contxt.AppContext from context")
	}

	wrapperCtx, ok := value.(*AppContext)
	if !ok {
		return nil, fmt.Errorf("could not get contxt.AppContext from context")
	}

	return wrapperCtx, nil
}

func ContextWithAppWrapper(ctx context.Context, wrapperCtx *AppContext) context.Context {
	// nolint
	return context.WithValue(ctx, AppCtxKey, wrapperCtx)
}

func (c *AppContext) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.keys == nil {
		c.keys = make(map[string]any)
	}

	c.keys[key] = value
}

func (c *AppContext) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists = c.keys[key]
	return
}

func (c *AppContext) GetInt64(key string) (int64, error) {
	if val, ok := c.Get(key); ok && val != nil {
		return utils.GetInt64(val)
	}
	return 0, fmt.Errorf("key: %s not found", key)
}

func (c *AppContext) GetString(key string) (string, error) {
	if val, ok := c.Get(key); ok && val != nil {
		s, ok := val.(string)
		if !ok {
			return "", fmt.Errorf("parse value: %v to string error", val)
		}
		return s, nil
	}

	return "", fmt.Errorf("key: %s not found", key)
}
