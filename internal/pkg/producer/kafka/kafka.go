package kafka

import "github.com/IBM/sarama"

// Kafka represents kafka broker instance
type Kafka struct {
	p sarama.SyncProducer
}

// New creates new kafka broker instance
func New(brokers []string, kafkaConfig Config) (*Kafka, error) {
	// create cfg
	cfg := sarama.NewConfig()

	// basic settings
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll

	// config settings
	cfg.Producer.Retry.Max = 5
	if kafkaConfig.MaxRetry != 0 {
		cfg.Producer.Retry.Max = kafkaConfig.MaxRetry
	}

	// create producer
	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	return &Kafka{
		p: producer,
	}, nil
}

func (k *Kafka) Close() error {
	return k.p.Close()
}
