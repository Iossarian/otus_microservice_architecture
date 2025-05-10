package app

import (
	"net/http"

	"orchestrator/internal/domain/saga/createorder"

	"github.com/labstack/echo/v4"
)

type Order interface {
	Create(ctx echo.Context, request createorder.Request) error
	Cancel(ctx echo.Context, request createorder.Request) error
}
type Handler struct {
	warehouse Warehouse
	billing   Billing
	delivery  Delivery
	order     Order
}

func NewHandler(
	warehouse Warehouse,
	billing Billing,
	delivery Delivery,
) *Handler {
	return &Handler{
		warehouse: warehouse,
		billing:   billing,
		delivery:  delivery,
	}
}

func internalError(err error) error {
	return &echo.HTTPError{
		Internal: err,
		Message:  err.Error(),
		Code:     http.StatusInternalServerError,
	}
}
