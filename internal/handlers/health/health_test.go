package health_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/UiP9AV6Y/fake-secrets/internal/assert"
	"github.com/UiP9AV6Y/fake-secrets/internal/handlers/health"
	status "github.com/UiP9AV6Y/fake-secrets/internal/health"
)

func TestHandlerServeHTTP(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	start := time.Unix(0, 0).UTC()
	testCases := map[string]struct {
		Want assert.Assertions[*http.Response]
	}{
		"OK": {
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertStatusStatus(
						assert.StringEqual(status.StatusOK),
					),
					assertStatusUptime(
						assert.IntGreaterThan(0),
					),
				),
			},
		},
	}

	for name, test := range testCases {
		scenario := func(t *testing.T) {
			subject := health.NewHandler(start, logger)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.test/health", nil)
			w := httptest.NewRecorder()

			subject.ServeHTTP(w, req)

			assert.Assert(t, w.Result(), test.Want)
		}

		t.Run(name, scenario)
	}
}
