package app

import (
	"context"

	"order/internal/domain"

	"github.com/labstack/echo/v4"
)

type Orchestrator interface {
	Exec(ctx echo.Context, request Request) error
}

type Producer interface {
	Produce(ctx context.Context, event domain.OrderCreatedEvent) error
}
