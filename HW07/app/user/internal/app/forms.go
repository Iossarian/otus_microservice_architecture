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
