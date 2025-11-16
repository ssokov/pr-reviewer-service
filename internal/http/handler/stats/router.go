package stats

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, handler *Handler) {
	e.GET("/stats", handler.GetStats)
}
