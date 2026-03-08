package assert

import (
	"encoding/json"
	"net/http"
	"testing"
)

func HTTPResponseStatusCode(want int) Assertion[*http.Response] {
	return func(t *testing.T, got *http.Response) {
		if got.StatusCode != want {
			t.Errorf("invalid status code. got %d, want %d", got.StatusCode, want)
		}
	}
}

func HTTPResponseBodyJSON[B any](assertions ...Assertion[B]) Assertion[*http.Response] {
	return func(t *testing.T, got *http.Response) {
		var dto B
		if err := json.NewDecoder(got.Body).Decode(&dto); err != nil {
			t.Fatal(err)
		} else {
			Assert(t, dto, assertions)
		}
	}
}
