package app

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"order/internal/domain"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

const (
	orderStatusSuccess = "success"
	orderStatusFailed  = "failed"
)

type Request struct {
	GoodQty      int     `json:"quantity"`
	Price        float64 `json:"price"`
	DeliverySlot string  `json:"slot"`
	UserID       int     `json:"user_id"`
}

func (h *Handler) Create(ctx echo.Context) error {
	var req Request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request",
		})
	}

	err := h.orchestrator.Exec(ctx, req)
	if err != nil {
		return internalError(err)
	}

	_, err = h.db.Exec("INSERT INTO orders (user_id, price) VALUES ($1, $2)",
		req.UserID,
		req.Price,
	)
	if err != nil {
		return h.produceMessage(ctx.Request().Context(), req.UserID, 0, req.Price, orderStatusFailed)
	}

	var orderID int
	err = h.db.QueryRowContext(
		ctx.Request().Context(),
		"SELECT id FROM orders WHERE user_id = $1 AND price = $2 LIMIT 1",
		req.UserID,
		req.Price,
	).Scan(&orderID)
	if err != nil {
		return h.produceMessage(ctx.Request().Context(), req.UserID, 0, req.Price, orderStatusFailed)
	}

	err = h.produceMessage(ctx.Request().Context(), req.UserID, orderID, req.Price, orderStatusSuccess)

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message": "order created",
	})
}

func (h *Handler) produceMessage(
	ctx context.Context,
	userID int,
	orderID int,
	price float64,
	status string,
) error {
	err := h.orderCreatedProducer.Produce(
		ctx,
		domain.OrderCreatedEvent{
			UserID:     strconv.Itoa(userID),
			OrderID:    strconv.Itoa(orderID),
			Price:      price,
			Status:     status,
			OccurredAt: time.Now(),
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to produce message")
	}

	return nil
}
