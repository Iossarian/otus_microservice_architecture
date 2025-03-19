package app

import (
	"context"
	"database/sql"
	"net/http"

	"user/internal/domain"
	"user/internal/jwt"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	db                  *sql.DB
	jwtService          *jwt.Service
	userCreatedProducer Producer
}

type Producer interface {
	Produce(ctx context.Context, event domain.UserCreatedEvent) error
}

func NewHandler(
	db *sql.DB,
	jwtService *jwt.Service,
	userCreatedProducer Producer,
) *Handler {
	return &Handler{
		db:                  db,
		jwtService:          jwtService,
		userCreatedProducer: userCreatedProducer,
	}
}

func internalError(err error) error {
	return &echo.HTTPError{
		Internal: err,
		Message:  "Internal error",
		Code:     http.StatusInternalServerError,
	}
}
