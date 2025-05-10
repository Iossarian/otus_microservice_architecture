package app

import (
	"context"

	"order/internal/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Orchestrator interface {
	Exec(ctx echo.Context, request Request) error
}

type Producer interface {
	Produce(ctx context.Context, event domain.OrderCreatedEvent) error
}
type IdempotencyKeyRepository interface {
	Exists(key uuid.UUID) (bool, error)
	Set(key uuid.UUID) error
}
