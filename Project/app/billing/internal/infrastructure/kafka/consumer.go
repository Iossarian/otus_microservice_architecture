package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"billing/internal/app"
	"billing/internal/domain"

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

		var u domain.User

		err = json.Unmarshal(msg.Value, &u)
		if err != nil {
			log.Errorf("unmarshal message %s", err.Error())

			continue
		}

		fmt.Println("User created: ", u)

		err = p.handler.CreateAccount(ctx, u)
		if err != nil {
			log.Errorf("create account for user %d, err %s", u.ID, err.Error())

			continue
		}
	}
}
