package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/flum1025/tweam/internal/config"
	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

var ctx = context.Background()

func NewRedisClient(
	config *config.Config,
) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Redis.Address,
		DB:   config.Redis.DB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %w", err)
	}

	return &RedisClient{
		client: rdb,
	}, nil
}

func (c *RedisClient) Get(key string) (*string, error) {
	val, err := c.client.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, fmt.Errorf("redis: get key: %w", err)
	}

	return &val, nil
}

func (c *RedisClient) Set(key string, val interface{}, expiration time.Duration) error {
	if err := c.client.Set(ctx, key, val, expiration).Err(); err != nil {
		return fmt.Errorf("redis: set and expire: %w", err)
	}

	return nil
}

func (c *RedisClient) Exists(key string) (bool, error) {
	res, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis: exists: %w", err)
	}

	return res == 1, nil
}
