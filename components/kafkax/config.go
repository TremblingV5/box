package kafkax

type Config struct {
	ProducerList map[string]*ProducerConfig `json:"producerList" yaml:"producerList"`
	ConsumerList map[string]*ConsumerConfig `yaml:"consumerList" json:"consumerList"`
}

type ProducerConfig struct {
	Brokers []string `yaml:"brokers" json:"brokers"`
	Topics  []string `yaml:"topics" json:"topics"`
}

type ConsumerConfig struct {
	Brokers []string `yaml:"brokers" json:"brokers"`
	GroupId string   `json:"groupId" yaml:"groupId"`
	Topics  []string `yaml:"topics" json:"topics"`
}

func (c *Config) SetDefault() {
	return
}
