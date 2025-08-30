package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ Cache = (*cache)(nil)

//go:generate mockgen -source=cache.go -destination=cache_mock.go -package=cache

// Cache interface for plain cache
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
	SetObject(ctx context.Context, key string, val interface{}, duration time.Duration) error
	Delete(ctx context.Context, key string) error
	SetWithDuration(ctx context.Context, key string, value []byte, duration time.Duration) error
	SetExpireTime(ctx context.Context, key string, seconds int64) error
	Exists(ctx context.Context, keys ...string) int64

	LSet(ctx context.Context, key string, vals []byte) error
	LLen(ctx context.Context, key string) (int64, error)
	LGet(ctx context.Context, key string) ([]byte, error)
	LList(ctx context.Context, key string) ([]string, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	DecrBy(ctx context.Context, key string, value int64) (int64, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)

	SetVal(ctx context.Context, key string, value string) error
	GetVal(ctx context.Context, key string) (string, error)
	LRange(ctx context.Context, key string, from int, to int) ([]string, error)
	ZAdd(ctx context.Context, key string, score float64, member string) error
	ZRange(ctx context.Context, key string, start int64, stop int64) ([]string, error)
	ZRemRangeByRank(ctx context.Context, key string, start int64, stop int64) error
	ZIncrBy(ctx context.Context, key string, increment float64, member string) error
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error)
	ZRevRank(ctx context.Context, key string, member string) (int64, error)
	ZScore(ctx context.Context, key string, member string) (float64, error)

	GetSMembers(ctx context.Context, key string) ([]string, error)
	SetSAdd(ctx context.Context, key string, members ...interface{}) error
	SetNX(ctx context.Context, key string, seconds int64, data interface{}) (bool, error)
}

// cache structure
type cache struct {
	rc        *redis.Client
	cacheTime time.Duration
}

// New initializes Cache
func New(rc *redis.Client, cacheTime time.Duration) Cache {
	return &cache{
		rc:        rc,
		cacheTime: cacheTime,
	}
}

// Get reads value by key
func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := c.rc.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("get key=%s failed err=%w", key, err)
	}
	return value, nil
}

// Set sets value by key
func (c *cache) Set(ctx context.Context, key string, value []byte) error {
	err := c.rc.Set(ctx, key, value, c.cacheTime).Err()
	if err != nil {
		return fmt.Errorf("set value=%v with key=%s to redis failed err=%w", value, key, err)
	}
	return nil
}

func (c *cache) SetObject(ctx context.Context, key string, val interface{}, duration time.Duration) error {
	dataBytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return c.rc.Set(ctx, key, dataBytes, duration).Err()
}

// Delete ...
func (c *cache) Delete(ctx context.Context, key string) error {
	err := c.rc.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("delete with key=%s to redis failed err=%w", key, err)
	}
	return nil
}

// SetWithDuration ...
func (c *cache) SetWithDuration(ctx context.Context, key string, value []byte, duration time.Duration) error {
	err := c.rc.Set(ctx, key, value, duration).Err()
	if err != nil {
		return fmt.Errorf("set value=%v with key=%s and duration=%d to redis failed err=%w", value, key, duration, err)
	}
	return nil
}

// LGet ...
func (c *cache) LGet(ctx context.Context, key string) ([]byte, error) {
	value, err := c.rc.LPop(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("get cache key=%s from redis failed err=%w", key, err)
	}
	return value, nil
}

// LSet ...
func (c *cache) LSet(ctx context.Context, key string, val []byte) error {
	return c.rc.LPush(ctx, key, val).Err()
}

// LLen ...
func (c *cache) LLen(ctx context.Context, key string) (int64, error) {
	val, err := c.rc.LLen(ctx, key).Result()
	if err != nil {
		return -1, fmt.Errorf("get len of key=%s from redis failed err=%w", key, err)
	}
	return val, nil
}

// LList ...
func (c *cache) LList(ctx context.Context, key string) ([]string, error) {
	vals, err := c.rc.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("get datas of key=%s from redis failed err=%w", key, err)
	}
	return vals, nil
}

// Decr ...
func (c *cache) Decr(ctx context.Context, key string) (int64, error) {
	val, err := c.rc.Decr(ctx, key).Result()
	if err != nil {
		return -1, fmt.Errorf("decrby key=%s from redis failed err=%w", key, err)
	}
	return val, err
}

// Incr ...
func (c *cache) Incr(ctx context.Context, key string) (int64, error) {
	val, err := c.rc.Incr(ctx, key).Result()
	if err != nil {
		return -1, fmt.Errorf("incrBy key=%s from redis failed err=%w", key, err)
	}
	return val, err
}

// DecrBy ...
func (c *cache) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	val, err := c.rc.DecrBy(ctx, key, value).Result()
	if err != nil {
		return -1, fmt.Errorf("decrby key=%s from redis failed err=%w", key, err)
	}
	return val, err
}

// IncrBy ...
func (c *cache) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	val, err := c.rc.IncrBy(ctx, key, value).Result()
	if err != nil {
		return -1, fmt.Errorf("incrBy key=%s from redis failed err=%w", key, err)
	}
	return val, err
}

// LRange ...
func (c *cache) LRange(ctx context.Context, key string, from int, to int) ([]string, error) {
	result, err := c.rc.LRange(ctx, key, int64(from), int64(to)).Result()
	if err != nil {
		return nil, fmt.Errorf("lrange cache key=%s from redis failed err=%w", key, err)
	}
	return result, nil
}

// ZIncrBy ...
func (c *cache) ZIncrBy(ctx context.Context, key string, increment float64, member string) error {
	err := c.rc.
		ZIncrBy(ctx, key, increment, member).
		Err()
	if err != nil {
		return fmt.Errorf("zincrby member=%v with key=%s and increment=%v to redis failed. Error: %w", member, key, increment, err)
	}

	return nil
}

// ZAdd ...
func (c *cache) ZAdd(ctx context.Context, key string, score float64, member string) error {
	err := c.rc.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
	if err != nil {
		return fmt.Errorf("zadd member=%v with key=%s and score=%v to redis failed, err=%w", member, key, score, err)
	}

	return nil
}

// ZRange ...
func (c *cache) ZRange(ctx context.Context, key string, start int64, stop int64) ([]string, error) {
	result, err := c.rc.ZRange(ctx, key, start, stop).Result()

	if err != nil {
		return nil, fmt.Errorf("zrange with key=%s, start=%v and stop=%v to redis failed, err=%w", key, start, stop, err)
	}

	return result, nil
}

// ZRevRangeWithScores ...
func (c *cache) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	result, err := c.rc.ZRevRangeWithScores(ctx, key, start, stop).Result()

	if err != nil {
		return nil, fmt.Errorf("ZRevRangeWithScores with key=%s, start=%v and stop=%v to redis failed. Error: %w", key, start, stop, err)
	}

	return result, nil
}

// ZRevRank ...
func (c *cache) ZRevRank(ctx context.Context, key string, member string) (int64, error) {
	result, err := c.rc.ZRevRank(ctx, key, member).Result()
	if err != nil {
		return -1, fmt.Errorf("ZRevRank with key=%s, member=%v. Error: %w", key, member, err)
	}

	return result, nil
}

// ZScore ...
func (c *cache) ZScore(ctx context.Context, key string, member string) (float64, error) {
	result, err := c.rc.ZScore(ctx, key, member).Result()
	if errors.Is(err, redis.Nil) {
		return -1, err
	}

	if err != nil {
		return -1, fmt.Errorf("ZRevRank with key=%s, member=%v. Error: %w", key, member, err)
	}

	return result, nil
}

// ZRemRangeByRank ...
func (c *cache) ZRemRangeByRank(ctx context.Context, key string, start int64, stop int64) error {
	err := c.rc.ZRemRangeByRank(ctx, key, start, stop).Err()

	if err != nil {
		return fmt.Errorf("zremrangebyrank with key=%s, start=%v and stop=%v to redis failed, err=%w", key, start, stop, err)
	}

	return nil
}

// SetVal set value by key
func (c *cache) SetVal(ctx context.Context, key string, value string) error {
	err := c.rc.Set(ctx, key, value, c.cacheTime).Err()
	if err != nil {
		return fmt.Errorf("set value=%v with key=%s to redis failed, err=%w", value, key, err)
	}
	return nil
}

// GetVal reads value by key
func (c *cache) GetVal(ctx context.Context, key string) (string, error) {
	return c.rc.Get(ctx, key).Result()
}

func (c *cache) SetExpireTime(ctx context.Context, key string, seconds int64) error {
	return c.rc.Expire(ctx, key, time.Second*time.Duration(seconds)).Err()
}

func (c *cache) Exists(ctx context.Context, keys ...string) int64 {
	return c.rc.Exists(ctx, keys...).Val()
}

func (c *cache) GetSMembers(ctx context.Context, key string) ([]string, error) {
	data, err := c.rc.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *cache) SetSAdd(ctx context.Context, key string, members ...interface{}) error {
	return c.rc.SAdd(ctx, key, members).Err()
}

func (c *cache) SetNX(ctx context.Context, key string, seconds int64, data interface{}) (bool, error) {
	return c.rc.SetNX(ctx, key, data, time.Second*time.Duration(seconds)).Result()
}
