package user

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, h *UserHandler) {
	userGroup := e.Group("/users")
	{
		userGroup.POST("/setIsActive", h.SetIsActive)
		userGroup.GET("/getReview", h.GetReview)
	}
}
