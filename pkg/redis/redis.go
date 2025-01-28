package redis

import (
	"context"
	"fmt"

	"github.com/idoyudha/eshop-order/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedis(cfg config.Redis) (*RedisClient, error) {
	client := &RedisClient{
		Client: redis.NewClient(RedisOptions(cfg)),
	}

	_, err := client.Client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis connection error: %w", err)
	}
	return client, nil
}
