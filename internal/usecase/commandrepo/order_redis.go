package commandrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	rClient "github.com/idoyudha/eshop-order/pkg/redis"
)

const orderKey = "order"

type OrderRedisRepo struct {
	*rClient.RedisClient
}

func NewOrderRedisRepo(client *rClient.RedisClient) *OrderRedisRepo {
	return &OrderRedisRepo{
		client,
	}
}

func getOrderKey(orderID uuid.UUID) string {
	return fmt.Sprintf("order:%s", orderID.String())
}

func (r *OrderRedisRepo) Set(ctx context.Context, orderID uuid.UUID, value string, ttl time.Duration) error {
	key := getOrderKey(orderID)
	return r.RedisClient.Client.Set(ctx, key, value, ttl).Err()
}

func (r *OrderRedisRepo) Delete(ctx context.Context, orderID uuid.UUID) error {
	key := getOrderKey(orderID)
	return r.RedisClient.Client.Del(ctx, key).Err()
}

func (r *OrderRedisRepo) GetTTL(ctx context.Context, orderID uuid.UUID) (time.Duration, error) {
	key := getOrderKey(orderID)
	return r.RedisClient.Client.TTL(ctx, key).Result()
}
