package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/idoyudha/eshop-order/config"
	"github.com/redis/go-redis/v9"
)

// type RedisClient struct {
// 	Client *redis.Client
// }

type RedisClient struct {
	Client redis.UniversalClient
}

func NewRedis(cfg config.Redis) (*RedisClient, error) {
	client := &RedisClient{
		// Client: redis.NewClient(RedisOptions(cfg)),
		Client: redis.NewFailoverClient(RedisFailoverOptions(cfg)),
	}

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		_, err := client.Client.Ping(context.Background()).Result()
		if err == nil {
			log.Println("connected to redis")
			return client, nil
		}
		log.Printf("failed to connect to redis (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to redis after %d attempts", maxRetries)
}
