package assert

import (
	"strings"
	"testing"
)

func StringNotEmpty(msg ...string) Assertion[string] {
	return func(t *testing.T, got string) {
		if got != "" {
			return
		}

		t.Error(format(msg, "got empty string"))
	}
}

func StringEqual(want string, msg ...string) Assertion[string] {
	return func(t *testing.T, got string) {
		if got == want {
			return
		}

		t.Error(format(msg, "got %q, want %q", got, want))
	}
}

func StringContains(want string, msg ...string) Assertion[string] {
	return func(t *testing.T, got string) {
		if strings.Contains(got, want) {
			return
		}

		t.Error(format(msg, "got %q, does not include %q", got, want))
	}
}

func StringLongerThan(want int, msg ...string) Assertion[string] {
	return func(t *testing.T, got string) {
		if len(got) > want {
			return
		}

		t.Error(format(msg, "got len(%q) == %d, want >%d", got, len(got), want))
	}
}

func StringShorterThan(want int, msg ...string) Assertion[string] {
	return func(t *testing.T, got string) {
		if len(got) < want {
			return
		}

		t.Error(format(msg, "got len(%q) == %d, want <%d", got, len(got), want))
	}
}
