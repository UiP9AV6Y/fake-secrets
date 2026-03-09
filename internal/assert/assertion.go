package assert

import (
	"testing"
)

type Assertion[A any] func(*testing.T, A)

type Assertions[A any] []Assertion[A]

func Assert[A any](t *testing.T, got A, a Assertions[A]) {
	for _, f := range a {
		f(t, got)
	}
}
