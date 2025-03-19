package build

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func (b *Builder) RestServer() *echo.Echo {
	e := echo.New()

	b.shutdown.add(func(ctx context.Context) error {
		return e.Close()
	})

	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler, err := b.handler()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.POST("/charge", handler.Charge)
	e.POST("/refund", handler.Refund)
	e.GET("/balance", handler.Balance)
	e.POST("/deposit", handler.Deposit)

	return e
}
