package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/idoyudha/eshop-order/config"
)

type ProducerServer struct {
	Producer *kafka.Producer
}

func NewKafkaProducer(kafkaCfg config.Kafka) (*ProducerServer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaCfg.Broker,
		"acks":               "all",
		"enable.idempotence": true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	return &ProducerServer{
		Producer: p,
	}, nil
}

func (s *ProducerServer) Close() {
	s.Producer.Close()
}

func (s *ProducerServer) Publish(topic string, key []byte, message interface{}) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal kafka message: %w", err)
	}
	return s.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          messageBytes,
	}, nil)
}
