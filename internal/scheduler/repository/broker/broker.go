package broker

import (
	"github.com/keenywheels/go-spy/internal/pkg/producer/kafka"
)

// Topics represents available topics
type Topics struct {
	ScraperData string
}

// Broker represents broker instance
type Broker struct {
	topics Topics
	kafka  *kafka.Kafka
}

// New creates new broker instance
func New(kafka *kafka.Kafka, topics Topics) *Broker {
	return &Broker{
		topics: topics,
		kafka:  kafka,
	}
}
