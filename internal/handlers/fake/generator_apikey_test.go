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

func TestGeneratorHandlerServeAPIKey(t *testing.T) {
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
						assert.StringLongerThan(39),
						assert.StringShorterThan(41),
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
						assert.StringEqual("fsp_Qo5oIGioaBwrvRLdWFkKxCxuJTeC7iCZscXT"),
					),
				),
			},
		},
		"seeded_valid_type": {
			HaveRequest: []requestOption{
				WithRequestPathValue("seed", "super-secret-value"),
				WithRequestQuery("type", "t"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.StringEqual("fst_Qo5oIGioaBwrvRLdWFkKxCxuJTeC7iCZscXT"),
					),
				),
			},
		},
		"seeded_valid_organization": {
			HaveRequest: []requestOption{
				WithRequestPathValue("seed", "super-secret-value"),
				WithRequestQuery("organization", "test"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.StringEqual("testp_Qo5oIGioaBwrvRLdWFkKxCxuJTeC7iCZscXT"),
					),
				),
			},
		},
		"seeded_invalid_type": {
			HaveRequest: []requestOption{
				WithRequestPathValue("seed", "super-secret-value"),
				WithRequestQuery("type", "tttttt"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusBadRequest),
				assert.HTTPResponseBodyJSON(
					assertDTOString("error",
						assert.StringEqual("token type is too long 6/5"),
					),
				),
			},
		},
		"seeded_invalid_organization": {
			HaveRequest: []requestOption{
				WithRequestPathValue("seed", "super-secret-value"),
				WithRequestQuery("organization", "oooooo"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusBadRequest),
				assert.HTTPResponseBodyJSON(
					assertDTOString("error",
						assert.StringEqual("token organization is too long 6/5"),
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
				WithRequestPath("apikeys"),
			}

			if len(test.HaveRequest) > 0 {
				reqopt = append(reqopt, test.HaveRequest...)
			}

			req := newRequest(t.Context(), reqopt...)
			w := httptest.NewRecorder()

			subject.ServeAPIKey(w, req)

			assert.Assert(t, w.Result(), test.Want)
		}

		t.Run(name, scenario)
	}
}
