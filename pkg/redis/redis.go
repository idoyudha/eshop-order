package redis

import (
	"context"
	"fmt"
	"log"
	"math"
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
	options := RedisFailoverOptions(cfg)
	client := &RedisClient{
		Client: redis.NewFailoverClient(options),
	}

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := client.Client.Ping(ctx).Result()
		cancel()

		if err == nil {
			log.Println("connected to redis")
			return client, nil
		}
		backoffDuration := time.Second * time.Duration(math.Pow(2, float64(i)))
		log.Printf("failed to connect to redis (attempt %d/%d): %v. waiting %v before next attempt",
			i+1, maxRetries, err, backoffDuration)
		time.Sleep(backoffDuration)
	}

	return nil, fmt.Errorf("failed to connect to redis after %d attempts", maxRetries)
}
