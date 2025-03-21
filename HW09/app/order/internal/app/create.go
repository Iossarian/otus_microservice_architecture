package app

import (
	"net/http"

	"order/internal/infrastructure"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Request struct {
	Price          float64 `form:"price"`
	UserID         int     `form:"user_id"`
	IdempotencyKey string  `form:"idempotency_key"`
}

func (h *Handler) Create(ctx echo.Context) error {
	var req Request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request",
		})
	}

	idempotencyKey, err := uuid.Parse(req.IdempotencyKey)
	if err != nil {
		return internalError(err)
	}

	order, err := h.orderRepository.Create(
		ctx.Request().Context(),
		idempotencyKey,
		uuid.New(),
		req.Price,
		req.UserID,
	)
	if err != nil {
		if errors.Is(err, infrastructure.ErrDuplicateRequest) {
			return ctx.JSON(http.StatusConflict, map[string]interface{}{
				"error": "duplicated request",
			})
		}

		return internalError(err)
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id": order.ID.String(),
	})
}
