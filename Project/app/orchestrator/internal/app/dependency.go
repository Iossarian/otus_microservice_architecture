package app

import (
	"orchestrator/internal/domain/saga/createorder"

	"github.com/labstack/echo/v4"
)

type Billing interface {
	Charge(ctx echo.Context, request createorder.Request) error
	Refund(ctx echo.Context, request createorder.Request) error
}

type Warehouse interface {
	Reserve(ctx echo.Context, request createorder.Request) error
	Release(ctx echo.Context, request createorder.Request) error
}

type Delivery interface {
	Reserve(ctx echo.Context, request createorder.Request) error
	Cancel(ctx echo.Context, request createorder.Request) error
}
