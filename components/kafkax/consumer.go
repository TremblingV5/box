package kafkax

import (
	"errors"
	"github.com/IBM/sarama"
)

type Consumer struct {
	Topic    string
	Consumer *sarama.ConsumerGroup
}

func GetConsumer(storeKey, topic string) (*Consumer, error) {
	clients, err := getKafkaClients(storeKey)
	if err != nil {
		return nil, err
	}

	if _, ok := clients.consumerTopics.Load(topic); !ok {
		return nil, errors.New("topic not found")
	}

	consumerValue, ok := clients.consumers.Load(storeKey)
	if !ok {
		return nil, errors.New("consumer not found")
	}

	consumer, ok := consumerValue.(*sarama.ConsumerGroup)
	if !ok {
		return nil, errors.New("consumer not found")
	}

	return &Consumer{
		Topic:    topic,
		Consumer: consumer,
	}, nil
}
