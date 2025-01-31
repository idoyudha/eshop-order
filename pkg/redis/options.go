package redis

import (
	"strings"
	"time"

	"github.com/idoyudha/eshop-order/config"
	"github.com/redis/go-redis/v9"
)

// func RedisOptions(cfg config.Redis) *redis.Options {
// 	return &redis.Options{
// 		Addr:     cfg.RedisURL,
// 		Password: cfg.RedisPassword,
// 	}
// }

func RedisFailoverOptions(cfg config.Redis) *redis.FailoverOptions {
	sentinelAddrs := strings.Split(cfg.RedisSentinelAddrs, ",")
	return &redis.FailoverOptions{
		MasterName:    cfg.RedisMaster,
		SentinelAddrs: sentinelAddrs,
		Password:      cfg.RedisPassword,
		DB:            0,
		ReadTimeout:   time.Second * 3,
		WriteTimeout:  time.Second * 3,
	}
}
