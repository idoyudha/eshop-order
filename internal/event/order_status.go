package event

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/config"
	"github.com/idoyudha/eshop-order/internal/entity"
	"github.com/idoyudha/eshop-order/internal/usecase"
	"github.com/idoyudha/eshop-order/pkg/logger"
	"github.com/idoyudha/eshop-order/pkg/redis"
)

const expiredKeyEventRedis = "__keyevent@0__:expired"

type redisScheduledEvents struct {
	r    *redis.RedisClient
	ucoc usecase.OrderCommand
	l    logger.Interface
	c    config.Constant
}

func NewRedisScheduledEvents(
	r *redis.RedisClient,
	ucoc usecase.OrderCommand,
	l logger.Interface,
	c config.Constant,
) error {
	events := &redisScheduledEvents{
		r:    r,
		ucoc: ucoc,
		l:    l,
		c:    c,
	}

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process
	log.Println("starting redis consumer in order service, consuming message from other producer...")
	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
			return nil
		default:
			pubsub := r.Client.PSubscribe(context.Background(), expiredKeyEventRedis)
			defer pubsub.Close()

			for msg := range pubsub.Channel() {
				if err := events.handleOrderExpired(msg.Payload); err != nil {
					events.l.Error(err, "http - v1 - redisScheduledEvents - handleMessage")
				}
			}
		}
	}
	return nil
}

func (e *redisScheduledEvents) handleOrderExpired(orderID string) error {
	e.l.Info("Order expired", "http - v1 - redisScheduledEvents - handleOrderExpired")

	order := entity.Order{
		ID: uuid.MustParse(orderID),
	}
	err := e.ucoc.UpdateOrderStatus(context.Background(), &order, entity.ORDER_EXPIRED)
	if err != nil {
		return err
	}

	return nil
}
