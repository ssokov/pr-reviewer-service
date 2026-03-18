package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/ssokov/pr-reviewer-service/docs"
	"github.com/ssokov/pr-reviewer-service/internal/pr"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/vmkteam/embedlog"
)

func NewServer(
	logger embedlog.Logger,
	prService *pr.PrService,
) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("logger", logger)
			return next(c)
		}
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	handler := NewHandler(prService, logger)
	RegisterRoutes(e, handler)

	return e
}
