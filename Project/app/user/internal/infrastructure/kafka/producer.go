package kafka

import (
	"context"
	"encoding/json"

	"user/internal/domain"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(writer *kafka.Writer) *Producer {
	return &Producer{
		writer: writer,
	}
}

func (p *Producer) Produce(ctx context.Context, e domain.UserCreatedEvent) error {
	str, err := json.Marshal(e)
	if err != nil {
		return errors.Wrap(err, "marshal event")
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Value: str,
	})
	if err != nil {
		return errors.Wrap(err, "write message")
	}

	return nil
}
