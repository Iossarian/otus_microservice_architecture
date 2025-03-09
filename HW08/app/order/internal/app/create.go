package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Request struct {
	GoodQty      int     `json:"quantity"`
	Price        float64 `json:"price"`
	DeliverySlot string  `json:"slot"`
}

func (h *Handler) Create(ctx echo.Context) error {
	type formRequest struct {
		GoodQty      int     `form:"quantity"`
		Price        float64 `form:"price"`
		DeliverySlot string  `form:"slot"`
	}

	var req formRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request",
		})
	}

	err := h.orchestrator.Exec(ctx, Request{
		GoodQty:      req.GoodQty,
		Price:        req.Price,
		DeliverySlot: req.DeliverySlot,
	})
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "order was not created",
		})
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message": "order created",
	})
}
