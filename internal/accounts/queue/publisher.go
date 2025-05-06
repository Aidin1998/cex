package queue

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"

	"cex/pkg/apiutil"
)

type Config struct {
	Accounts struct {
		Port string
		DSN  string
	}
	Queue struct {
		URL    string
		Topics []string
	}
}
type Publisher struct {
	writer *kafka.Writer
}

// NewPublisher returns a Kafka-based event publisher.
// brokers: []string{"localhost:9092"}, topic must be non-empty.
func NewPublisher(brokers []string, topic string) *Publisher {
	return &Publisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// PublishAccountCreated sends an AccountCreatedEvent.
func (p *Publisher) PublishAccountCreated(ctx context.Context, e apiutil.AccountCreatedEvent) error {
	key := e.EventID.String()
	msg, _ := json.Marshal(e)
	return p.writer.WriteMessages(ctx, kafka.Message{Key: []byte(key), Value: msg})
}

// PublishBalanceUpdated sends a BalanceUpdatedEvent.
func (p *Publisher) PublishBalanceUpdated(ctx context.Context, e apiutil.BalanceUpdatedEvent) error {
	key := e.EventID.String()
	msg, _ := json.Marshal(e)
	return p.writer.WriteMessages(ctx, kafka.Message{Key: []byte(key), Value: msg})
}

// Close closes the Kafka writer.
func (p *Publisher) Close() error {
	return p.writer.Close()
}
