package redis

import (
	"github.com/idoyudha/eshop-order/config"
	"github.com/redis/go-redis/v9"
)

func RedisOptions(cfg config.Redis) *redis.Options {
	return &redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
	}
}
