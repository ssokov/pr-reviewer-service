package pr

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo, p *PRHandler) {
	prGroup := e.Group("/pullRequest")
	{
		prGroup.POST("/create", p.CreatePR)
		prGroup.POST("/merge", p.MergePR)
		prGroup.POST("/reassign", p.ReassignReviewer)
	}
}
