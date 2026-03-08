package fake_test

import (
	"log/slog"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/UiP9AV6Y/fake-secrets/internal/assert"
	"github.com/UiP9AV6Y/fake-secrets/internal/handlers/fake"
)

func TestGeneratorHandlerServePassword(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	testCases := map[string]struct {
		HaveRequest []requestOption
		Want        assert.Assertions[*http.Response]
	}{
		"default": {
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.StringLongerThan(11),
						assert.StringShorterThan(13),
					),
				),
			},
		},
		"length": {
			HaveRequest: []requestOption{
				WithRequestQuery("length", "16"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.StringLongerThan(15),
						assert.StringShorterThan(17),
					),
				),
			},
		},
	}

	for name, test := range testCases {
		scenario := func(t *testing.T) {
			seed := rand.NewSource(0)
			rnd := rand.New(seed)
			subject := fake.NewGeneratorHandler(rnd, logger)
			reqopt := []requestOption{
				WithRequestPath("passwords"),
			}

			if len(test.HaveRequest) > 0 {
				reqopt = append(reqopt, test.HaveRequest...)
			}

			req := newRequest(t.Context(), reqopt...)
			w := httptest.NewRecorder()

			subject.ServePassword(w, req)

			assert.Assert(t, w.Result(), test.Want)
		}

		t.Run(name, scenario)
	}
}
