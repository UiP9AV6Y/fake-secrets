package assert

import (
	"testing"
)

func IntGreaterThan(want int, msg ...string) Assertion[int] {
	return func(t *testing.T, got int) {
		if got > want {
			return
		}

		t.Error(format(msg, "got %d, want >%d", got, want))
	}
}

func IntLessThan(want int, msg ...string) Assertion[int] {
	return func(t *testing.T, got int) {
		if got < want {
			return
		}

		t.Error(format(msg, "got %d, want <%d", got, want))
	}
}
