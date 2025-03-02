package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) OrderForm(c echo.Context) error {
	return c.Render(http.StatusOK, "order.html", nil)
}
