package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cex/internal/accounts/model"
	"cex/pkg/apiutil"

	"github.com/segmentio/kafka-go"
	"github.com/sony/gobreaker"
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

var writer *kafka.Writer // Initialize in app.go

type Publisher struct {
	writer  *kafka.Writer
	breaker *gobreaker.CircuitBreaker
}

// NewPublisher returns a Kafka-based event publisher with circuit breaker and retry logic.
// brokers: []string{"localhost:9092"}, topic must be non-empty.
func NewPublisher(brokers []string, topic string) *Publisher {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "AccountPublisher",
		MaxRequests: 5,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
	})
	return &Publisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
		breaker: cb,
	}
}

// PublishAccountCreated sends an AccountCreatedEvent with retry logic and circuit breaker.
func (p *Publisher) PublishAccountCreated(ctx context.Context, e apiutil.AccountCreatedEvent) error {
	key := e.EventID.String()
	msgBytes, _ := json.Marshal(e)

	_, err := p.breaker.Execute(func() (interface{}, error) {
		for i, backoff := 0, time.Millisecond*100; i < 3; i, backoff = i+1, backoff*2 {
			if err := p.writer.WriteMessages(ctx, kafka.Message{Key: []byte(key), Value: msgBytes}); err != nil {
				time.Sleep(backoff)
				continue
			}
			return nil, nil
		}
		return nil, fmt.Errorf("publish AccountCreatedEvent failed after retries")
	})
	return err
}

// PublishBalanceUpdated sends a BalanceUpdatedEvent.
func (p *Publisher) PublishBalanceUpdated(ctx context.Context, e apiutil.BalanceUpdatedEvent) error {
	key := e.EventID.String()
	msg, _ := json.Marshal(e)
	return p.writer.WriteMessages(ctx, kafka.Message{Key: []byte(key), Value: msg})
}

// PublishAccountCreated publishes an account creation event.
func PublishAccountCreated(ctx context.Context, acct model.Account) error {
	msg := apiutil.Event{
		Type:    "account.created",
		Payload: acct,
	}
	return writer.WriteMessages(ctx, kafka.Message{
		Key: []byte(acct.ID.String()),
		Value: func() []byte {
			data, _ := json.Marshal(msg)
			return data
		}(),
	})
}

// Close closes the Kafka writer.
func (p *Publisher) Close() error {
	return p.writer.Close()
}
