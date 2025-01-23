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
