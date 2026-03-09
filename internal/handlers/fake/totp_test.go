package fake_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/UiP9AV6Y/fake-secrets/internal/assert"
	"github.com/UiP9AV6Y/fake-secrets/internal/handlers/fake"
	"github.com/UiP9AV6Y/fake-secrets/internal/io"
)

func TestTOTPHandlerServePrivateKey(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	testCases := map[string]struct {
		HaveAccount string
		HaveRequest []requestOption
		Want        assert.Assertions[*http.Response]
	}{
		"random_default": {
			HaveAccount: "default",
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.OTP(
							assert.OTPAccountName("default"),
							assert.OTPAlgorithm("SHA1"),
							assert.OTPDigits(6),
							assert.OTPIssuer("Vault"),
							assert.OTPType("totp"),
							assert.OTPSecretValue("ORXXI4BNOJQW4ZDPNUWXGZLFMR2G65DQ"),
						),
					),
				),
			},
		},
		"organization": {
			HaveAccount: "organization",
			HaveRequest: []requestOption{
				WithRequestQuery("organization", "Spec"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.OTP(
							assert.OTPIssuer("Spec"),
						),
					),
				),
			},
		},
		"length": {
			HaveAccount: "length",
			HaveRequest: []requestOption{
				WithRequestQuery("length", "28"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.OTP(
							assert.OTPSecretValue("ORXXI4BNOJQW4ZDPNUWXGZLFMR2G65DQFVZGC3TEN5WS2"),
						),
					),
				),
			},
		},
		"algorithm": {
			HaveAccount: "algorithm",
			HaveRequest: []requestOption{
				WithRequestQuery("algorithm", "md5"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.OTP(
							assert.OTPAlgorithm("MD5"),
						),
					),
				),
			},
		},
	}

	for name, test := range testCases {
		scenario := func(t *testing.T) {
			rnd := io.InfiniteReader([]byte("totp-random-seed"))
			subject := fake.NewTOTPHandler(rnd, logger)
			reqopt := []requestOption{
				WithRequestPath("totp"),
				WithRequestPathValue("account", test.HaveAccount),
				WithRequestPath("keys"),
			}

			if len(test.HaveRequest) > 0 {
				reqopt = append(reqopt, test.HaveRequest...)
			}

			req := newRequest(t.Context(), reqopt...)
			w := httptest.NewRecorder()

			subject.ServePrivateKey(w, req)

			assert.Assert(t, w.Result(), test.Want)
		}

		t.Run(name, scenario)
	}
}
