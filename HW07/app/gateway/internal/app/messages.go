package app

import (
	"net/http"
	"strconv"

	"gateway/infrastructure/notification"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Messages(ctx echo.Context) error {
	tokenString, err := ctx.Request().Cookie("Authorization")
	if err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to get token",
			Code:     http.StatusBadRequest,
		}
	}

	token, err := h.jwtService.Parse(tokenString.Value)
	if err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to parse token",
			Code:     http.StatusBadRequest,
		}
	}

	userID, err := strconv.Atoi(token.ID)
	if err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to convert user ID",
			Code:     http.StatusBadRequest,
		}
	}

	messages, err := h.notificationClient.Message(
		ctx.Request().Context(),
		notification.Request{
			UserID: userID,
		},
	)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusOK, messages)
}
