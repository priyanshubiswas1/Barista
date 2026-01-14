package session

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type IRedisClient interface {
	SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetValue(ctx context.Context, key string) (string, error)
	DeleteKey(ctx context.Context, key string) error
}

type redisClient struct {
	client *redis.Client
}

func NewRedisClient(redisUrl string) (IRedisClient, error) {
	var opt *redis.Options
	var err error

	if strings.HasPrefix(redisUrl, "redis://") || strings.HasPrefix(redisUrl, "rediss://") {
		opt, err = redis.ParseURL(redisUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse redis url: %v", err)
		}
	} else {
		opt = &redis.Options{
			Addr: redisUrl,
		}
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis: %v", err)
	}

	return &redisClient{
		client: client,
	}, nil
}

func (c *redisClient) SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *redisClient) GetValue(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *redisClient) DeleteKey(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
