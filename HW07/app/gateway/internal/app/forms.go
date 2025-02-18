package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterForm(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

func (h *Handler) LoginForm(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func (h *Handler) DepositForm(c echo.Context) error {
	return c.Render(http.StatusOK, "deposit.html", nil)
}

func (h *Handler) WithdrawForm(c echo.Context) error {
	return c.Render(http.StatusOK, "withdraw.html", nil)
}

func (h *Handler) OrderForm(c echo.Context) error {
	return c.Render(http.StatusOK, "order.html", nil)
}
