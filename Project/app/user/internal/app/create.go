package app

import (
	"net/http"
	"strconv"

	"user/internal/domain"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Create(ctx echo.Context) error {
	type NewUser struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}

	var u NewUser

	if err := ctx.Bind(&u); err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to bind user",
			Code:     http.StatusBadRequest,
		}
	}

	if u.Name == "" || u.Password == "" {
		return &echo.HTTPError{
			Message: "invalid name and/or password",
			Code:    http.StatusBadRequest,
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return internalError(err)
	}

	_, err = h.db.ExecContext(
		ctx.Request().Context(),
		"INSERT INTO users (name, password) VALUES ($1, $2)", u.Name, string(hashedPassword))
	if err != nil {
		return internalError(err)
	}

	var userID int
	err = h.db.QueryRowContext(ctx.Request().Context(), "SELECT id FROM users WHERE name = $1", u.Name).Scan(&userID)
	if err != nil {
		return internalError(err)
	}

	err = h.userCreatedProducer.Produce(ctx.Request().Context(), domain.UserCreatedEvent{
		ID:   strconv.Itoa(userID),
		Name: u.Name,
	})

	return ctx.JSON(http.StatusCreated, u)
}
