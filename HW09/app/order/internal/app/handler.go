package app

import (
	"context"
	"net/http"

	"order/internal/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrderRepository interface {
	Create(
		ctx context.Context,
		key uuid.UUID,
		id uuid.UUID,
		price float64,
		userID int,
	) (*domain.Order, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Order, error)
}

type IdempotencyKeyRepository interface {
	Delete(ctx context.Context, key uuid.UUID) error
}

type Handler struct {
	orderRepository          OrderRepository
	idempotencyKeyRepository IdempotencyKeyRepository
}

func NewHandler(
	orderRepository OrderRepository,
	idempotencyKeyRepository IdempotencyKeyRepository,
) *Handler {
	return &Handler{
		orderRepository:          orderRepository,
		idempotencyKeyRepository: idempotencyKeyRepository,
	}
}

func internalError(err error) error {
	return &echo.HTTPError{
		Internal: err,
		Message:  err.Error(),
		Code:     http.StatusInternalServerError,
	}
}
