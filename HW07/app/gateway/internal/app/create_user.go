package app

import (
	"net/http"

	"gateway/infrastructure/user"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateUser(ctx echo.Context) error {
	type User struct {
		Name     string `form:"username"`
		Password string `form:"password"`
	}

	var u User

	if err := ctx.Bind(&u); err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to bind user",
			Code:     http.StatusBadRequest,
		}
	}

	err := h.userClient.Create(ctx, user.User{
		Name:     u.Name,
		Password: u.Password,
	})

	if err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to create user",
			Code:     http.StatusInternalServerError,
		}
	}

	return ctx.Redirect(http.StatusMovedPermanently, "/login")
}
