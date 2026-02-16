package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oklog/run"

	"github.com/UiP9AV6Y/fake-secrets/internal/config"
	"github.com/UiP9AV6Y/fake-secrets/internal/handlers"
	"github.com/UiP9AV6Y/fake-secrets/internal/version"
)

func Run(ctx context.Context, cfg *config.Config, logger *slog.Logger) error {
	var g run.Group

	ln, err := net.Listen("tcp", cfg.Listen())
	if err != nil {
		return err
	}

	handler, err := handlers.NewRouter(cfg, logger)
	if err != nil {
		return err
	}

	server := http.Server{
		Addr:         ln.Addr().String(),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), cfg.LogVerbosity()),
	}
	serverStart := func() error {
		logger.Info("Listening on", "address", ln.Addr().String())

		if err := server.Serve(ln); err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}
	serverStop := func(_ error) {
		graceCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		logger.Info("Server is shutting down")

		if err := server.Shutdown(graceCtx); err != nil {
			logger.Error("server shutdown interrupted", "err", err)
		}
	}
	g.Add(serverStart, serverStop)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	sigCtx, sigCancel := context.WithCancel(ctx)
	sigStart := func() error {
		select {
		case s := <-sig:
			logger.Info("received shutdown signal", "signal", s)
			return nil
		case <-sigCtx.Done():
			return sigCtx.Err()
		}
	}
	sigStop := func(_ error) {
		logger.Debug("dismantling shutdown handler")
		sigCancel()
	}
	g.Add(sigStart, sigStop)

	logger.Info("starting fake-secrets service", "version", version.Version())

	return g.Run()
}
