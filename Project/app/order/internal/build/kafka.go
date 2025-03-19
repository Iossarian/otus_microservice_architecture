package build

import (
	"context"

	kfk "order/internal/infrastructure/kafka"

	"github.com/segmentio/kafka-go"
)

func (b *Builder) OrderCreatedProducer() *kfk.Producer {
	writer := kafka.Writer{
		Addr:  kafka.TCP(b.config.Kafka.Broker),
		Topic: b.config.Kafka.OrderCreatedTopic,
	}

	b.shutdown.add(func(ctx context.Context) error {
		return writer.Close()
	})

	return kfk.NewProducer(&writer)
}
