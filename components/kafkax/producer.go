package kafkax

import (
	"errors"
	"github.com/IBM/sarama"
)

type Producer struct {
	Topic    string
	Producer *sarama.SyncProducer
}

func GetProducer(storeKey, topic string) (*Producer, error) {
	clients, err := getKafkaClients(storeKey)
	if err != nil {
		return nil, err
	}

	if _, ok := clients.producerTopics.Load(topic); !ok {
		return nil, errors.New("topic not found")
	}

	producerValue, ok := clients.producers.Load(storeKey)
	if !ok {
		return nil, errors.New("producer not found")
	}

	producer, ok := producerValue.(*sarama.SyncProducer)
	if !ok {
		return nil, errors.New("producer not found")
	}

	return &Producer{
		Topic:    topic,
		Producer: producer,
	}, nil
}
