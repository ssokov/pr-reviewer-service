package team

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, handler *TeamHandler) {
	e.POST("/team/add", handler.AddTeam)
	e.GET("/team/get", handler.GetTeam)
	e.POST("/team/deactivate", handler.DeactivateTeam)
}
