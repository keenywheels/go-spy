package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

// Message represents a Kafka message with a topic and a value
type Message struct {
	Topic string
	Value any
}

// ProduceJSON produces a JSON message to Kafka
func (k *Kafka) ProduceJSON(msg Message) error {
	jsonMsg, err := k.getJSON(msg.Value)
	if err != nil {
		return fmt.Errorf("failed to get JSON message: %w", err)
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: msg.Topic,
		Value: jsonMsg,
	}

	_, _, err = k.p.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	return nil
}

// getJSON processes the JSON value of the message
func (k *Kafka) getJSON(value any) (sarama.StringEncoder, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return sarama.StringEncoder(""), fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return sarama.StringEncoder(bytes), nil
}
