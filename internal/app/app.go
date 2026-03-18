package app

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	config "github.com/ssokov/pr-reviewer-service/cfg"
	"github.com/ssokov/pr-reviewer-service/internal/db"
	"github.com/ssokov/pr-reviewer-service/internal/pr"
	"github.com/ssokov/pr-reviewer-service/internal/rest"
	"github.com/vmkteam/embedlog"
)

type App struct {
	sl     embedlog.Logger
	config *config.Config
	db     *pg.DB
	echo   *echo.Echo

	prService *pr.PrService
}

func New(slogger embedlog.Logger, c *config.Config, db *pg.DB) *App {
	a := &App{
		config: c,
		db:     db,
		sl:     slogger,
	}
	a.initDependencies()

	a.echo = rest.NewServer(
		a.sl,
		a.prService,
	)
	return a
}

func (a *App) initDependencies() {
	repo := db.NewPrReviewerServiceRepo(a.db)
	a.prService = pr.NewPrService(repo, a.sl)
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
