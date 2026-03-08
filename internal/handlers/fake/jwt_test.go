package fake_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/UiP9AV6Y/fake-secrets/internal/assert"
	"github.com/UiP9AV6Y/fake-secrets/internal/handlers/fake"
	"github.com/UiP9AV6Y/fake-secrets/internal/io"
)

func TestJWTHandlerServeToken(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	start := time.Unix(0, 0).UTC()
	testCases := map[string]struct {
		HaveSubject string
		HaveRequest []requestOption
		Want        assert.Assertions[*http.Response]
	}{
		"default": {
			HaveSubject: "default",
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.JWT(
							assert.JWTNoAudience(),
							assert.JWTExpiration(
								assert.TimeEqual(start.Add(2*time.Hour)),
							),
							assert.JWTIssuedAt(
								assert.TimeEqual(start),
							),
							assert.JWTIssuer(
								assert.StringEqual("http://example.test"),
							),
							assert.JWTJwtID(
								assert.StringEqual("6a77742d-7261-4e64-af6d-2d736565646a"),
							),
							assert.JWTNotBefore(
								assert.TimeEqual(start),
							),
							assert.JWTSubject(
								assert.StringEqual("default"),
							),
						),
					),
				),
			},
		},
		"issuer_header_host": {
			HaveSubject: "issuer_header_host",
			HaveRequest: []requestOption{
				WithRequestHeader("X-Forwarded-Host", "jwt-spec.test"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.JWT(
							assert.JWTIssuer(
								assert.StringEqual("http://jwt-spec.test"),
							),
						),
					),
				),
			},
		},
		"issuer_header_proto": {
			HaveSubject: "issuer_header_proto",
			HaveRequest: []requestOption{
				WithRequestHeader("X-Forwarded-Proto", "https"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.JWT(
							assert.JWTIssuer(
								assert.StringEqual("https://example.test"),
							),
						),
					),
				),
			},
		},
		"audience": {
			HaveSubject: "audience",
			HaveRequest: []requestOption{
				WithRequestQuery("audience", "http://example.com:1234"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.JWT(
							assert.JWTAudience(
								assert.StringEqual("http://example.com:1234"),
							),
						),
					),
				),
			},
		},
	}

	for name, test := range testCases {
		scenario := func(t *testing.T) {
			rnd := io.InfiniteReader([]byte("jwt-random-seed"))
			subject := fake.NewJWTHandler(start, rnd, logger)
			reqopt := []requestOption{
				WithRequestPath("jwt"),
				WithRequestPathValue("subject", test.HaveSubject),
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
