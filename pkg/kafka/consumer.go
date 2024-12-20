package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	SaleGroup           = "sale-group"
	PaymentCreatedTopic = "payment-created"
)

type ConsumerServer struct {
	Consumer *kafka.Consumer
}

func NewKafkaConsumer(brokerURL string) (*ConsumerServer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokerURL,
		"group.id":          SaleGroup,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	err = c.SubscribeTopics([]string{
		PaymentCreatedTopic,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topics: %v", err)
	}

	return &ConsumerServer{
		Consumer: c,
	}, nil
}

func (c *ConsumerServer) Close() error {
	return c.Consumer.Close()
}
