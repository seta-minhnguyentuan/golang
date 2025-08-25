package kafka

import (
	"context"
	"encoding/json"
	"shared/pkg/log"
	"time"
	"user-service/internal/config"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Publish(ctx context.Context, topic string, key []byte, v any) error
	Close() error
}

type producer struct {
	w *kafka.Writer
}

func NewProducer(cfg *config.KafkaConfig) Producer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.KafkaBrokers...),
		BatchBytes:             cfg.KafkaBatchBytes,
		BatchTimeout:           cfg.KafkaBatchTimeout,
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.Hash{},
	}
	return &producer{w: w}
}

func (p *producer) Publish(ctx context.Context, topic string, key []byte, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Time:  time.Now(),
		Key:   key,
		Value: b,
	}
	err = p.w.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   msg.Key,
		Value: msg.Value,
		Time:  msg.Time,
	})
	if err != nil {
		log.Error.Printf("producer write error: %v", err)
	}
	return err
}

func (p *producer) Close() error {
	return p.w.Close()
}
