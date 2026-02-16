package health

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

const URLPath = "health"

func Run(ctx context.Context, endpoint *url.URL, logger *slog.Logger) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return err
	}

	now := time.Now()
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = rsp.Body.Close()
	}()

	logger.Info("server responded", "duration", time.Since(now), "code", rsp.StatusCode)

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("health endpoint responded with status code %d", rsp.StatusCode)
	}

	var status Status
	if err := json.NewDecoder(rsp.Body).Decode(&status); err != nil {
		return err
	}

	logger.Info("server status received", "status", status.Status, "uptime", status.Uptime)

	if status.Status != StatusOK {
		return fmt.Errorf("health endpoint responded with status %q", status.Status)
	}

	return nil
}
