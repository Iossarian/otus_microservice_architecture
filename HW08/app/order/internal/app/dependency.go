package app

import (
	"github.com/labstack/echo/v4"
)

type Orchestrator interface {
	Exec(ctx echo.Context, request Request) error
}

type Billing interface {
	Balance(ctx echo.Context) (float64, error)
}

type Warehouse interface {
	Stock(ctx echo.Context) (int, error)
}

type Delivery interface {
	Slots(ctx echo.Context) ([]string, error)
}
