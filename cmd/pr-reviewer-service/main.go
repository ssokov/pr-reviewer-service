// @title PR Reviewer Service API
// @version 1.0
// @description API for managing pull request reviews and team assignments
// @BasePath /
// @host localhost:8080
// @schemes rest
package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-pg/pg/v10"
	config "github.com/ssokov/pr-reviewer-service/cfg"
	"github.com/ssokov/pr-reviewer-service/internal/app"
	"github.com/vmkteam/embedlog"
)

var (
	flVerbose = flag.Bool("verbose", false, "print verbose output")
	flJSON    = flag.Bool("json", false, "print output as JSON")
	flDev     = flag.Bool("dev", true, "uses development mode")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	sl := embedlog.NewLogger(*flVerbose, *flJSON)
	if *flDev {
		sl = embedlog.NewDevLogger()
	}
	slog.SetDefault(sl.Log())

	cfg, err := config.Load("cfg/config.toml")
	if err != nil {
		sl.Errorf("Failed to load config. error: %v", err)
		exitOnError(err)
	}

	dsn := cfg.Database.DSN()
	sl.Print(ctx, "connecting to database", "host", cfg.Database.Host, "database", cfg.Database.Database, "dsn", dsn)

	opt, err := pg.ParseURL(dsn)
	if err != nil {
		sl.Errorf("failed to parse pg config: %v", err)
		exitOnError(err)
	}

	db := pg.Connect(opt)
	defer func() {
		_ = db.Close()
	}()

	if err := db.Ping(ctx); err != nil {
		sl.Errorf("db ping failed: %v", err)
		exitOnError(err)
	}

	var version string
	if _, err := db.QueryOneContext(ctx, pg.Scan(&version), "select version()"); err != nil {
		sl.Errorf("failed to get version: %v", err)
		exitOnError(err)
	}
	sl.Print(ctx, "connected to db", "version", version)

	application := app.New(sl, cfg, db)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Run
	go func() {
		if err := application.Run(ctx); err != nil {
			exitOnError(err)
		}
	}()
	<-quit

	sl.Print(ctx, "Application finished")
}

func exitOnError(err error) {
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
