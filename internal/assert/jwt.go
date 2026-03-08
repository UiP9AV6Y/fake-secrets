package assert

import (
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwt"
)

func JWT(assertions ...Assertion[jwt.Token]) Assertion[string] {
	return func(t *testing.T, got string) {
		if token, err := jwt.ParseInsecure([]byte(got)); err != nil {
			t.Fatalf("malformed JWT input")
		} else {
			Assert(t, token, assertions)
		}
	}
}

func JWTNoAudience(msg ...string) Assertion[jwt.Token] {
	return func(t *testing.T, got jwt.Token) {
		if _, ok := got.Audience(); ok {
			t.Error(format(msg, "JWT audience: unwanted claim present"))
		}
	}
}

func JWTAudience(assertions ...Assertion[string]) Assertion[jwt.Token] {
	return func(t *testing.T, got jwt.Token) {
		if value, ok := got.Audience(); !ok || len(value) == 0 {
			t.Errorf("JWT audience: no such claim")
		} else {
			Assert(t, value[0], assertions)
		}
	}
}

func JWTNoExpiration(msg ...string) Assertion[jwt.Token] {
	return func(t *testing.T, got jwt.Token) {
		if _, ok := got.Expiration(); ok {
			t.Error(format(msg, "JWT expiration: unwanted claim present"))
		}
	}
}

func JWTExpiration(assertions ...Assertion[time.Time]) Assertion[jwt.Token] {
	return func(t *testing.T, got jwt.Token) {
		if value, ok := got.Expiration(); !ok {
			t.Errorf("JWT expiration: no such claim")
		} else {
			Assert(t, value, assertions)
		}
	}
}

func JWTNoIssuedAt(msg ...string) Assertion[jwt.Token] {
	return func(t *testing.T, got jwt.Token) {
		if _, ok := got.IssuedAt(); ok {
			t.Error(format(msg, "JWT issued_at: unwanted claim present"))
		}
	}
}

func JWTIssuedAt(assertions ...Assertion[time.Time]) Assertion[jwt.Token] {
	return func(t *testing.T, got jwt.Token) {
		if value, ok := got.IssuedAt(); !ok {
			t.Errorf("JWT issued_at: no such claim")
		} else {
			Assert(t, value, assertions)
		}
	}
}

func JWTNoIssuer(msg ...string) Assertion[jwt.Token] {
	return func(t *testing.T, got jwt.Token) {
		if _, ok := got.Issuer(); ok {
			t.Error(format(msg, "JWT issuer: unwanted claim present"))
		}
	}
}

func JWTIssuer(assertions ...Assertion[string]) Assertion[jwt.Token] {
	return func(t *testing.T, got jwt.Token) {
		if value, ok := got.Issuer(); !ok {
			t.Errorf("JWT issuer: no such claim")
		} else {
			Assert(t, value, assertions)
		}
	}
}

func JWTNoJwtID(msg ...string) Assertion[jwt.Token] {
	return func(t *testing.T, got jwt.Token) {
		if _, ok := got.JwtID(); ok {
			t.Error(format(msg, "JWT jwt_id: unwanted claim present"))
		}
	}
}

func JWTJwtID(assertions ...Assertion[string]) Assertion[jwt.Token] {
	return func(t *testing.T, got jwt.Token) {
		if value, ok := got.JwtID(); !ok {
			t.Errorf("JWT jwt_id: no such claim")
		} else {
			Assert(t, value, assertions)
		}
	}
}

func JWTNoNotBefore(msg ...string) Assertion[jwt.Token] {
	return func(t *testing.T, got jwt.Token) {
		if _, ok := got.NotBefore(); ok {
			t.Error(format(msg, "JWT not_before: unwanted claim present"))
		}
	}
}

func JWTNotBefore(assertions ...Assertion[time.Time]) Assertion[jwt.Token] {
	return func(t *testing.T, got jwt.Token) {
		if value, ok := got.NotBefore(); !ok {
			t.Errorf("JWT not_before: no such claim")
		} else {
			Assert(t, value, assertions)
		}
	}
}

func JWTNoSubject(msg ...string) Assertion[jwt.Token] {
	return func(t *testing.T, got jwt.Token) {
		if _, ok := got.Subject(); ok {
			t.Error(format(msg, "JWT subject: unwanted claim present"))
		}
	}
}

func JWTSubject(assertions ...Assertion[string]) Assertion[jwt.Token] {
	return func(t *testing.T, got jwt.Token) {
		if value, ok := got.Subject(); !ok {
			t.Errorf("JWT subject: no such claim")
		} else {
			Assert(t, value, assertions)
		}
	}
}
