package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/UiP9AV6Y/fake-secrets/internal/config"
	"github.com/UiP9AV6Y/fake-secrets/internal/health"
	"github.com/UiP9AV6Y/fake-secrets/internal/server"
	"github.com/UiP9AV6Y/fake-secrets/internal/version"
)

func main() {
	os.Exit(run(os.Args, os.Stderr, os.Stdout))
}

func run(argv []string, o, e io.Writer) int {
	ctx := context.Background()
	cfg := config.New(filepath.Base(argv[0]))

	logger := slog.New(cfg.LogHandler(e))
	if err := cfg.LoadEnv(); err != nil {
		logger.Error("unable to load config from environment", "err", err)
		return 1
	}

	if usage, err := cfg.LoadArgs(argv[1:]); err != nil {
		logger.Error("unable to load config from commandline arguments", "err", err)
		return 1
	} else if usage != nil {
		usage(o)
		return 0
	}

	logger = slog.New(cfg.LogHandler(e))

	switch cfg.Command {
	case config.CommandHealthCheck:
		return runHealthCheck(ctx, cfg, logger)
	case config.CommandVersion:
		return version.Run(o)
	case config.CommandServe:
		return runServer(ctx, cfg, logger)
	default:
		logger.Error("invalid application command", "command", cfg.Command)
		return 2
	}
}

func runHealthCheck(ctx context.Context, cfg *config.Config, logger *slog.Logger) int {
	endpoint, err := cfg.SelfURL(health.URLPath)
	if err != nil {
		logger.Error("unable to create healthcheck URL", "err", err)
		return 1
	}

	if err := health.Run(ctx, endpoint, logger); err != nil {
		logger.Error("healthcheck request failed", "err", err)
		return 3
	} else {
		logger.Info("healthcheck request succeeded")
		return 0
	}
}

func runServer(ctx context.Context, cfg *config.Config, logger *slog.Logger) int {
	if err := server.Run(ctx, cfg, logger); err != nil {
		logger.Error("unable to run server", "err", err)
		return 1
	}

	return 0
}
