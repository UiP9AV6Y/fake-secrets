package assert

import (
	"testing"
	"time"
)

func TimeEqual(want time.Time, msg ...string) Assertion[time.Time] {
	return func(t *testing.T, got time.Time) {
		if got.Equal(want) {
			return
		}

		t.Error(format(msg, "got %s, want %s", got, want))
	}
}

func TimeBefore(want time.Time, msg ...string) Assertion[time.Time] {
	return func(t *testing.T, got time.Time) {
		if got.Before(want) {
			return
		}

		t.Error(format(msg, "time %s is not before %s", got, want))
	}
}

func TimeAfter(want time.Time, msg ...string) Assertion[time.Time] {
	return func(t *testing.T, got time.Time) {
		if got.After(want) {
			return
		}

		t.Error(format(msg, "time %s is not after %s", got, want))
	}
}
