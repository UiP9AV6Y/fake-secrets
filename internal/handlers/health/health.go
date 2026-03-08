package health

import (
	"log/slog"
	nethttp "net/http"
	"time"

	"github.com/UiP9AV6Y/fake-secrets/internal/config"
	status "github.com/UiP9AV6Y/fake-secrets/internal/health"
	"github.com/UiP9AV6Y/fake-secrets/internal/http"
)

type Handler struct {
	logger *slog.Logger
	start  time.Time
}

func NewHandler(start time.Time, logger *slog.Logger) *Handler {
	result := &Handler{
		logger: logger,
		start:  start,
	}

	return result
}

func (h *Handler) Route(cfg *config.Config) string {
	return cfg.HandlerPattern(status.URLPath)
}

func (h *Handler) ServeHTTP(w nethttp.ResponseWriter, _ *nethttp.Request) {
	h.logger.Debug("serving health status")

	dto := status.NewStatus(int(time.Since(h.start).Seconds()))

	http.ServeJSON(w, dto)
}
