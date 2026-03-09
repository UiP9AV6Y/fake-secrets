package fake_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/UiP9AV6Y/fake-secrets/internal/assert"
	"github.com/UiP9AV6Y/fake-secrets/internal/handlers/fake"
)

func TestFileHandlerServeHTTP(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	root := os.DirFS("testdata/secrets")
	testCases := map[string]struct {
		HaveRequest []requestOption
		Want        assert.Assertions[*http.Response]
	}{
		"regular_file": {
			HaveRequest: []requestOption{
				WithRequestPathValue("filename", "test.txt"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.StringEqual("test"),
					),
				),
			},
		},
		"empty_file": {
			HaveRequest: []requestOption{
				WithRequestPathValue("filename", "empty.txt"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusOK),
				assert.HTTPResponseBodyJSON(
					assertDTOString("secret",
						assert.StringEqual(""),
					),
				),
			},
		},
		"nonexistent_file": {
			HaveRequest: []requestOption{
				WithRequestPathValue("filename", "404.txt"),
			},
			Want: assert.Assertions[*http.Response]{
				assert.HTTPResponseStatusCode(http.StatusNotFound),
				assert.HTTPResponseBodyJSON(
					assertDTOString("error",
						assert.StringContains("no such file"),
					),
				),
			},
		},
	}

	for name, test := range testCases {
		scenario := func(t *testing.T) {
			subject := fake.NewFileHandler(root, logger)
			reqopt := []requestOption{
				WithRequestPath("files"),
			}

			if len(test.HaveRequest) > 0 {
				reqopt = append(reqopt, test.HaveRequest...)
			}

			req := newRequest(t.Context(), reqopt...)
			w := httptest.NewRecorder()

			subject.ServeHTTP(w, req)

			assert.Assert(t, w.Result(), test.Want)
		}

		t.Run(name, scenario)
	}
}
