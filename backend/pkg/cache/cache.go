package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/base-go/backend/pkg/config"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Ping(ctx context.Context) error
	Close() error
	Exists(ctx context.Context, key string) *redis.IntCmd

	// Queue operations (Redis List)
	LPush(ctx context.Context, key string, values ...interface{}) error
	BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error)
	LLen(ctx context.Context, key string) (int64, error)
}

type cache struct {
	client *redis.Client
}

type errCache error

var (
	CacheNil errCache = errors.New("cache is nil")
)

func NewCache() Cache {
	cfg := config.GetConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Cache.Host, cfg.Cache.Port),
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.Db,
	})

	return &cache{
		client: rdb,
	}
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", CacheNil
		}
		return "", err
	}

	return val, nil
}

func (c *cache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *cache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *cache) Close() error {
	return c.client.Close()
}

func (c *cache) Exists(ctx context.Context, key string) *redis.IntCmd {
	return c.client.Exists(ctx, key)
}

func (c *cache) LPush(ctx context.Context, key string, values ...interface{}) error {
	return c.client.LPush(ctx, key, values...).Err()
}

func (c *cache) BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	result, err := c.client.BRPop(ctx, timeout, keys...).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (c *cache) LLen(ctx context.Context, key string) (int64, error) {
	return c.client.LLen(ctx, key).Result()
}
