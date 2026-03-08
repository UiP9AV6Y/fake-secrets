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

func TestGeneratorHandlerServeToken(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	testCases := map[string]struct {
		HaveRequest []requestOption
		Want        assert.Assertions[*http.Response]
	}{
		"random_default": {
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.StringLongerThan(35),
						assert.StringShorterThan(37),
					),
				),
			},
		},
		"seeded_default": {
			HaveRequest: []requestOption{
				WithRequestPathValue("seed", "super-secret-value"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.StringEqual("60a37f27-3aac-3eff-b2ef-652d193db18e"),
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
				WithRequestPath("tokens"),
			}

			if len(test.HaveRequest) > 0 {
				reqopt = append(reqopt, test.HaveRequest...)
			}

			req := newRequest(t.Context(), reqopt...)
			w := httptest.NewRecorder()

			subject.ServeToken(w, req)

			assert.Assert(t, w.Result(), test.Want)
		}

		t.Run(name, scenario)
	}
}
