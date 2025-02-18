package app

import (
	"fmt"
	"net/http"
	"strconv"

	"gateway/infrastructure/billing"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Deposit(ctx echo.Context) error {
	type Withdraw struct {
		Amount float64 `form:"amount"`
	}

	var w Withdraw

	if err := ctx.Bind(&w); err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to bind deposit",
			Code:     http.StatusBadRequest,
		}
	}

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

	balance, err := h.billingClient.Deposit(
		ctx,
		billing.Request{
			UserID: userID,
			Amount: w.Amount,
		})
	if err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to deposit",
			Code:     http.StatusInternalServerError,
		}
	}

	return ctx.JSON(http.StatusOK, map[string]string{"balance": fmt.Sprintf("%.2f", balance)})
}
