package session

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type IRedisClient interface {
	SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetValue(ctx context.Context, key string) (string, error)
}

type redisClient struct {
	client *redis.Client
}

func NewRedisClient(redisUrl string) (IRedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})

	pong, err := client.Ping(context.Background()).Result()
	fmt.Println(pong, err)
	return &redisClient{
		client: client,
	}, err
}

func (c *redisClient) SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *redisClient) GetValue(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}
