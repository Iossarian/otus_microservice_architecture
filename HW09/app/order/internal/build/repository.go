package build

import (
	"order/internal/infrastructure"

	"github.com/pkg/errors"
)

func (b *Builder) orderRepository() (*infrastructure.OrderRepository, error) {
	db, err := b.postgres()
	if err != nil {
		return nil, errors.Wrap(err, "build postgres")
	}

	return infrastructure.NewOrderRepository(db), nil
}

func (b *Builder) idempotencyKeyRepository() (*infrastructure.IdempotencyKeyRepository, error) {
	db, err := b.postgres()
	if err != nil {
		return nil, errors.Wrap(err, "build postgres")
	}

	return infrastructure.NewIdempotencyKeyRepository(db), nil
}
