package broker

import (
	"github.com/keenywheels/go-spy/internal/pkg/producer/kafka"
	"github.com/keenywheels/go-spy/internal/scheduler/models"
)

// SendScraperData sends scraper data to the specified topic
func (b *Broker) SendScraperData(event models.ScraperEvent) error {
	kafkaMsg := kafka.Message{
		Topic: b.topics.ScraperData,
		Value: event,
	}

	return b.kafka.ProduceJSON(kafkaMsg)
}
