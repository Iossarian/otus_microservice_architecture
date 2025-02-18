package app

import (
	"context"
	"database/sql"
	"net/http"

	"order/infrastructure/billing"
	"order/internal/domain"

	"github.com/labstack/echo/v4"
)

type Producer interface {
	Produce(ctx context.Context, event domain.OrderCreatedEvent) error
}

type Handler struct {
	db                   *sql.DB
	billingClient        *billing.Client
	orderCreatedProducer Producer
}

func NewHandler(
	db *sql.DB,
	billingClient *billing.Client,
	orderCreatedProducer Producer,
) *Handler {
	return &Handler{
		db:                   db,
		billingClient:        billingClient,
		orderCreatedProducer: orderCreatedProducer,
	}
}

func internalError(err error) error {
	return &echo.HTTPError{
		Internal: err,
		Message:  "Internal error",
		Code:     http.StatusInternalServerError,
	}
}
