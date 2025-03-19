package build

import (
	"order/internal/app"

	"github.com/pkg/errors"
)

func (b *Builder) handler() (*app.Handler, error) {
	orderRepository, err := b.orderRepository()
	if err != nil {
		return nil, errors.Wrap(err, "build order repository")
	}

	idempotencyKeyRepository, err := b.idempotencyKeyRepository()
	if err != nil {
		return nil, errors.Wrap(err, "build idempotency key repository")
	}

	return app.NewHandler(
		orderRepository,
		idempotencyKeyRepository,
	), nil
}
