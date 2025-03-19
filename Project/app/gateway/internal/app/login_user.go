package app

import (
	"net/http"

	"gateway/infrastructure/user"

	"github.com/labstack/echo/v4"
)

func (h *Handler) LoginUser(ctx echo.Context) error {
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

	jwtToken, err := h.userClient.Login(ctx, user.User{
		Name:     u.Name,
		Password: u.Password,
	})
	if err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to login user",
			Code:     http.StatusInternalServerError,
		}
	}

	ctx.SetCookie(&http.Cookie{Name: "Authorization", Value: jwtToken, Path: "/", MaxAge: 3600})
	ctx.Response().Header().Set("Authorization", jwtToken)

	return ctx.JSON(http.StatusOK, map[string]string{"token": jwtToken})
}
