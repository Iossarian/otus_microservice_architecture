package app

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) DeleteKey(ctx echo.Context) error {
	keyStr := ctx.Param("key")
	if keyStr == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request",
		})
	}

	key, err := uuid.Parse(keyStr)
	if err != nil {
		return internalError(err)
	}

	err = h.idempotencyKeyRepository.Delete(ctx.Request().Context(), key)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "key deleted " + key.String(),
	})

}
