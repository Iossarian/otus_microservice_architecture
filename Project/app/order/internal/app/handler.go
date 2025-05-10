package app

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	db                       *sql.DB
	orchestrator             Orchestrator
	orderCreatedProducer     Producer
	idempotencyKeyRepository IdempotencyKeyRepository
}

func NewHandler(
	db *sql.DB,
	orchestrator Orchestrator,
	orderCreatedProducer Producer,
	idempotencyKeyRepository IdempotencyKeyRepository,
) *Handler {
	return &Handler{
		db:                       db,
		orchestrator:             orchestrator,
		orderCreatedProducer:     orderCreatedProducer,
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
