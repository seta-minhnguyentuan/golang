package kafka

import (
	"context"
	"errors"
	"shared/pkg/log"
	"user-service/internal/config"

	"github.com/segmentio/kafka-go"
)

type HandlerFunc func(ctx context.Context, key []byte, value []byte) error

type Consumer struct {
	r       *kafka.Reader
	handler HandlerFunc
}

func NewConsumer(cfg *config.KafkaConfig, topic string, handler HandlerFunc) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.KafkaBrokers,
		GroupID:        cfg.KafkaGroupID,
		Topic:          topic,
		MinBytes:       cfg.KafkaMinBytes,
		MaxBytes:       cfg.KafkaMaxBytes,
		MaxWait:        cfg.KafkaMaxWait,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 0,
	})
	return &Consumer{r: r, handler: handler}
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		m, err := c.r.FetchMessage(ctx)
		if err != nil {
			// ctx canceled on shutdown
			if errors.Is(err, context.Canceled) {
				return nil
			}
			log.Error.Printf("fetch error: %v", err)
			continue
		}

		if err := c.handler(ctx, m.Key, m.Value); err != nil {
			log.Error.Printf("handler error (no commit): %v", err)
			// Optionally: dead-letter here by producing to a DLQ topic
			// _ = dlq.Publish(ctx, "orders.dlq", m.Key, map[string]any{"raw": string(m.Value), "error": err.Error()})
			continue
		}

		if err := c.r.CommitMessages(ctx, m); err != nil {
			log.Error.Printf("commit error: %v", err)
			continue
		}
	}
}

func (c *Consumer) Close() error {
	return c.r.Close()
}
