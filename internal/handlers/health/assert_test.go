package health_test

import (
	"testing"

	"github.com/UiP9AV6Y/fake-secrets/internal/assert"
	status "github.com/UiP9AV6Y/fake-secrets/internal/health"
)

func assertStatusStatus(assertions ...assert.Assertion[string]) assert.Assertion[status.Status] {
	return func(t *testing.T, got status.Status) {
		assert.Assert(t, got.Status, assertions)
	}
}

func assertStatusUptime(assertions ...assert.Assertion[int]) assert.Assertion[status.Status] {
	return func(t *testing.T, got status.Status) {
		assert.Assert(t, got.Uptime, assertions)
	}
}
