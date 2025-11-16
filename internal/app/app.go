package app

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	config "github.com/ssokov/pr-reviewer-service/cfg"
	"github.com/ssokov/pr-reviewer-service/internal/http"
	"github.com/ssokov/pr-reviewer-service/internal/repository/postgres"
	"github.com/ssokov/pr-reviewer-service/internal/service"
	"github.com/vmkteam/embedlog"
)

type App struct {
	sl      embedlog.Logger
	appName string
	config  *config.Config
	db      *pgxpool.Pool
	echo    *echo.Echo

	userService  service.UserService
	prService    service.PRService
	teamService  service.TeamService
	statsService service.StatsService
}

func New(appName string, slogger embedlog.Logger, c *config.Config, db *pgxpool.Pool) *App {
	a := &App{
		appName: appName,
		config:  c,
		db:      db,
		sl:      slogger,
	}
	a.initDependencies()

	a.echo = http.NewServer(
		a.sl,
		a.userService,
		a.prService,
		a.teamService,
		a.statsService,
	)
	return a
}

func (a *App) initDependencies() {
	// init repositories
	userRepo := postgres.NewUserRepository(a.db)
	teamRepo := postgres.NewTeamRepository(a.db)
	prRepo := postgres.NewPRRepository(a.db)
	statsRepo := postgres.NewStatsRepository(a.db)

	// init services
	a.prService = service.NewPRService(prRepo, userRepo, teamRepo, a.sl)
	a.teamService = service.NewTeamService(teamRepo, userRepo, prRepo, a.sl)
	a.userService = service.NewUserService(userRepo, teamRepo, a.sl)
	a.statsService = service.NewStatsService(statsRepo, a.sl)
}

func (a *App) Run(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)
	a.sl.Print(ctx, "starting server", "addr", addr)

	serverErr := make(chan error, 1)
	go func() {
		if err := a.echo.Start(addr); err != nil {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return a.echo.Shutdown(shutdownCtx)
	}
}
