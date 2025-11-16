package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/ssokov/pr-reviewer-service/docs"
	"github.com/ssokov/pr-reviewer-service/internal/http/handler/pr"
	"github.com/ssokov/pr-reviewer-service/internal/http/handler/stats"
	"github.com/ssokov/pr-reviewer-service/internal/http/handler/team"
	"github.com/ssokov/pr-reviewer-service/internal/http/handler/user"
	"github.com/ssokov/pr-reviewer-service/internal/service"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/vmkteam/embedlog"
)

func NewServer(
	logger embedlog.Logger,
	userService service.UserService,
	prService service.PRService,
	teamService service.TeamService,
	statsService service.StatsService,
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

	userHandler := user.NewHandler(userService, logger)
	prHandler := pr.NewHandler(prService, logger)
	teamHandler := team.NewHandler(teamService, logger)
	statsHandler := stats.NewHandler(statsService, logger)

	user.RegisterRoutes(e, userHandler)
	pr.RegisterRoutes(e, prHandler)
	team.RegisterRoutes(e, teamHandler)
	stats.RegisterRoutes(e, statsHandler)

	return e
}
