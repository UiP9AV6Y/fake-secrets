package fake_test

import (
	"testing"

	"github.com/UiP9AV6Y/fake-secrets/internal/assert"
)

type DTO = map[string]any

func assertDTOString(field string, assertions ...assert.Assertion[string]) assert.Assertion[DTO] {
	return func(t *testing.T, got DTO) {
		dtoField, ok := got[field]
		if !ok {
			t.Fatalf("DTO %q field does not exist", field)

			return
		}

		if dto, ok := dtoField.(string); !ok {
			t.Fatalf("DTO %q field is not a string but a %T", field, dtoField)
		} else {
			assert.Assert(t, dto, assertions)
		}
	}
}
