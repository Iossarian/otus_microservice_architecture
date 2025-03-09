package app

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(
	db *sql.DB,
) *Handler {
	return &Handler{
		db: db,
	}
}

func internalError(err error) error {
	return &echo.HTTPError{
		Internal: err,
		Message:  "Internal error",
		Code:     http.StatusInternalServerError,
	}
}
