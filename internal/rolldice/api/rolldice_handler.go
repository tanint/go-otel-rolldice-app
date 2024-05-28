package api

import (
	"net/http"

	"github.com/demo/rolldice/internal/rolldice/services"
	"github.com/labstack/echo/v4"
)

type RolldiceHandler struct {
	rolldiceService *services.RollDiceService
}

func InitRolldiceHandler(e *echo.Echo, rolldiceService *services.RollDiceService) {
	handler := &RolldiceHandler{
		rolldiceService,
	}

	e.Group("/")
	e.GET("/roll", handler.Roll)
}

func (h *RolldiceHandler) Roll(c echo.Context) error {
	result := h.rolldiceService.Dice(c.Request().Context())

	return c.JSON(http.StatusOK, map[string]int{"result": result})
}
