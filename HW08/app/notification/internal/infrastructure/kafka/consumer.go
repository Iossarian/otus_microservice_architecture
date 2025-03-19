package kafka

import (
	"context"
	"encoding/json"
	"io"

	"notification/internal/app"
	"notification/internal/domain"

	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	handler *app.Handler
}

func NewConsumer(
	reader *kafka.Reader,
	handler *app.Handler,
) *Consumer {
	return &Consumer{
		reader:  reader,
		handler: handler,
	}
}

func (p *Consumer) Consume(ctx context.Context) error {
	for {
		msg, err := p.reader.ReadMessage(ctx)
		if err != nil && !errors.Is(io.EOF, err) {
			log.Errorf("read message %s", err.Error())

			continue
		}

		var m domain.Message

		err = json.Unmarshal(msg.Value, &m)
		if err != nil {
			log.Errorf("unmarshal message %s", err.Error())

			continue
		}

		err = p.handler.Notify(ctx, m)
		if err != nil {
			log.Errorf("notify user %s, err %s", m.UserID, err.Error())

			continue
		}
	}
}
