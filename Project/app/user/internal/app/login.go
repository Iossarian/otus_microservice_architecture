package app

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Login(ctx echo.Context) error {
	type User struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}

	var u User

	if err := ctx.Bind(&u); err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to bind user",
			Code:     http.StatusBadRequest,
		}
	}

	if u.Name == "" || u.Password == "" {
		return ctx.String(http.StatusBadRequest, "username or password is empty")
	}

	var userID int
	var storedHashedPassword string

	err := h.db.QueryRow("SELECT id, password FROM users WHERE name = $1", u.Name).
		Scan(&userID, &storedHashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.String(http.StatusUnauthorized, "user not found")
		}

		return internalError(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(u.Password))
	if err != nil {
		return &echo.HTTPError{
			Message: "invalid password",
			Code:    http.StatusUnauthorized,
		}
	}

	token, err := h.jwtService.Token(strconv.Itoa(userID), u.Name)
	if err != nil {
		return internalError(err)
	}

	type Resp struct {
		Token string `json:"token"`
	}

	resp := Resp{
		Token: token,
	}

	return ctx.JSON(http.StatusOK, resp)
}
