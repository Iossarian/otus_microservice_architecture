package build

import (
	"context"

	kfk "billing/internal/infrastructure/kafka"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func (b *Builder) UserCreatedConsumer(ctx context.Context) (*kfk.Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{b.config.Kafka.Broker},
		GroupID: "billing-service",
		Topic:   b.config.Kafka.UserCreatedTopic,
	})

	b.shutdown.add(func(ctx context.Context) error {
		return reader.Close()
	})

	handler, err := b.handler()
	if err != nil {
		return nil, errors.Wrap(err, "build handler")
	}

	return kfk.NewConsumer(reader, handler), nil
}
