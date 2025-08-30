package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Item interface define method CacheKey for all objects need add to Cache Storage
type Item interface {
	// Build key to store to Key-Value Cache system.
	CacheKey() string
}

// RedisRepo implement ICacheManager interface
type RedisRepo struct {
	rc        *redis.Client
	cacheTime time.Duration
}

// NewRedisRepo returns a new RedisRepo
func NewRedisRepo(rc *redis.Client, cacheTime time.Duration) *RedisRepo {
	return &RedisRepo{
		rc:        rc,
		cacheTime: cacheTime,
	}
}

// GetByCacheKey fetch CacheItem from Redis.
// Result is returned via parameter pointer result
// Result is object implemented Item with method CacheKey()
func (r *RedisRepo) GetByCacheKey(ctx context.Context, result Item) error {
	key := result.CacheKey()
	data, err := r.rc.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("get cache key=%s from redis failed err=%w", key, err)
	}

	err = json.Unmarshal([]byte(data), result)
	if err != nil {
		return fmt.Errorf("parse cache key=%s value='%s' from redis failed err=%w", key, data, err)
	}

	return nil
}

// GetByCacheKeys ...
func (r *RedisRepo) GetByCacheKeys(ctx context.Context, keys []string) ([]interface{}, error) {
	results, err := r.rc.MGet(ctx, keys...).Result()
	if err != nil {
		return []interface{}{}, fmt.Errorf("get cache keys=%v from redis failed err=%w", keys, err)
	}

	return results, nil
}

// Create stores CacheItem into Redis
func (r *RedisRepo) Create(ctx context.Context, item Item) error {
	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("marshal CacheItem=%v failed err=%w", item, err)
	}

	err = r.rc.Set(ctx, item.CacheKey(), data, r.cacheTime).Err()
	if err != nil {
		return fmt.Errorf("set CacheItem=%v to redis failed err=%w", item, err)
	}

	return nil
}

// UpdateWithMap stores CacheItem into Redis and update its value
func (r *RedisRepo) UpdateWithMap(ctx context.Context, item Item, params map[string]interface{}) error {
	bytes, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("marshal CacheItem=%v failed err=%w", item, err)
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return fmt.Errorf("unmarshal bytes=%v failed err=%w", bytes, err)
	}
	for key, value := range params {
		data[key] = value
	}

	cacheData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal CacheItem=%v failed err=%w", cacheData, err)
	}

	err = r.rc.Set(ctx, item.CacheKey(), cacheData, r.cacheTime).Err()
	if err != nil {
		return fmt.Errorf("set CacheItem=%v to redis failed", item)
	}

	return nil
}

// CreateList ...
func (r *RedisRepo) CreateList(ctx context.Context, items []Item) error {
	args := []interface{}{}
	for _, item := range items {
		key := item.CacheKey()
		data, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("marshal CacheItem=%v failed err=%w", item, err)
		}
		args = append(args, key, string(data))
	}

	err := r.rc.MSet(ctx, args...).Err()
	if err != nil {
		return fmt.Errorf("mset CacheItem=%v to redis failed err=%w", args, err)
	}

	return nil
}

// Delete remove CacheItem out Redis
func (r *RedisRepo) Delete(ctx context.Context, cacheItem Item) error {
	return r.rc.Del(ctx, cacheItem.CacheKey()).Err()
}
