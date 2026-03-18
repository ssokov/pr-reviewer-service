package rest

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, handler *PRHandler) {
	prGroup := e.Group("/pullRequest")
	{
		prGroup.POST("/create", handler.CreatePR)
		prGroup.POST("/merge", handler.MergePR)
		prGroup.POST("/reassign", handler.ReassignReviewer)

		e.GET("/stats", handler.GetStats)

		e.POST("/team/add", handler.AddTeam)
		e.GET("/team/get", handler.GetTeam)
	}

	userGroup := e.Group("/user")
	{
		userGroup.POST("/setIsActive", handler.SetIsActive)
		userGroup.GET("/getReview", handler.GetReview)
	}
}
