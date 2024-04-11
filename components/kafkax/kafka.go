package kafkax

import (
	"errors"
	"github.com/IBM/sarama"
	"github.com/TremblingV5/box/components"
	"sync"
	"time"
)

var (
	globalKafkaClientsMap = sync.Map{}
	globalConfigMap       = make(map[string]*Config)
)

type KafkaClients struct {
	producers      sync.Map
	producerTopics sync.Map
	consumers      sync.Map
	consumerTopics sync.Map
}

func GetConfig() components.ConfigMap[*Config] {
	return globalConfigMap
}

func Init(cm components.ConfigMap[*Config]) error {
	globalConfigMap = cm

	for k, v := range cm {
		clients := &KafkaClients{}

		clients.ConnectProducers(v.ProducerList)

		clients.ConnectConsumers(v.ConsumerList)

		globalKafkaClientsMap.Store(k, clients)
	}

	return nil
}

func (k *KafkaClients) ConnectProducers(configs map[string]*ProducerConfig) error {
	for name, config := range configs {
		conf := sarama.NewConfig()
		conf.Consumer.Offsets.AutoCommit.Enable = true
		conf.Consumer.Offsets.AutoCommit.Interval = time.Second * 1
		conf.Producer.Retry.Max = 1
		conf.Producer.RequiredAcks = sarama.WaitForAll
		conf.Producer.Return.Successes = true
		conf.Metadata.Full = true
		conf.Version = sarama.V0_10_2_0
		conf.Metadata.Full = true

		for _, topic := range config.Topics {
			k.producerTopics.Store(topic, true)
		}

		producer, err := sarama.NewSyncProducer(config.Brokers, conf)
		if err != nil {
			panic(err)
		}

		k.producers.Store(name, producer)
	}

	return nil
}

func (k *KafkaClients) ConnectConsumers(configs map[string]*ConsumerConfig) error {
	for name, config := range configs {
		conf := sarama.NewConfig()
		conf.Consumer.Offsets.AutoCommit.Enable = true
		conf.Consumer.Offsets.AutoCommit.Interval = time.Second * 1
		conf.Producer.Retry.Max = 1
		conf.Producer.RequiredAcks = sarama.WaitForAll
		conf.Producer.Return.Successes = true
		conf.Metadata.Full = true
		conf.Version = sarama.V0_10_2_0
		conf.Metadata.Full = true

		for _, topic := range config.Topics {
			k.consumerTopics.Store(topic, true)
		}

		consumer, err := sarama.NewConsumerGroup(config.Brokers, config.GroupId, conf)
		if err != nil {
			panic(err)
		}

		k.producers.Store(name, consumer)
	}

	return nil
}

func getKafkaClients(storeKey string) (*KafkaClients, error) {
	KafkaClientsValue, ok := globalKafkaClientsMap.Load(storeKey)
	if !ok {
		return nil, errors.New("kafka clients not found")
	}

	kafkaClients, ok := KafkaClientsValue.(*KafkaClients)
	if !ok {
		return nil, errors.New("kafka clients not found")
	}

	return kafkaClients, nil
}
