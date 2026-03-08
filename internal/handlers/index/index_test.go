package index_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/UiP9AV6Y/fake-secrets/internal/assert"
	"github.com/UiP9AV6Y/fake-secrets/internal/handlers/index"
)

func TestHandlerServeHTTP(t *testing.T) {
	testCases := map[string]struct {
		Want assert.Assertions[*http.Response]
	}{
		"OK": {
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusNotFound),
			},
		},
	}

	for name, test := range testCases {
		scenario := func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.test/", nil)
			w := httptest.NewRecorder()

			index.ServeHTTP(w, req)

			assert.Assert(t, w.Result(), test.Want)
		}

		t.Run(name, scenario)
	}
}
