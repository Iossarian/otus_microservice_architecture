package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	warehouse Warehouse
	billing   Billing
	delivery  Delivery
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
