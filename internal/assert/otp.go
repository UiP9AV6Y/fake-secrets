package assert

import (
	"testing"

	"github.com/pquerna/otp"
)

func OTP(assertions ...Assertion[*otp.Key]) Assertion[string] {
	return func(t *testing.T, got string) {
		if key, err := otp.NewKeyFromURL(got); err != nil {
			t.Fatalf("malformed OTP input")
		} else {
			Assert(t, key, assertions)
		}
	}
}

func OTPAccountName(want string, msg ...string) Assertion[*otp.Key] {
	return func(t *testing.T, got *otp.Key) {
		if value := got.AccountName(); value != want {
			t.Error(format(msg, "OTP accountname: got %q, want %q", value, want))
		}
	}
}

func OTPAlgorithm(want string, msg ...string) Assertion[*otp.Key] {
	return func(t *testing.T, got *otp.Key) {
		if value := got.Algorithm().String(); value != want {
			t.Error(format(msg, "OTP algorithm: got %q, want %q", value, want))
		}
	}
}

func OTPDigits(want int, msg ...string) Assertion[*otp.Key] {
	return func(t *testing.T, got *otp.Key) {
		if value := int(got.Digits()); value != want {
			t.Error(format(msg, "OTP digits: got %d, want %d", value, want))
		}
	}
}

func OTPEncoder(want string, msg ...string) Assertion[*otp.Key] {
	return func(t *testing.T, got *otp.Key) {
		if value := string(got.Encoder()); value != want {
			t.Error(format(msg, "OTP encoder: got %q, want %q", value, want))
		}
	}
}

func OTPIssuer(want string, msg ...string) Assertion[*otp.Key] {
	return func(t *testing.T, got *otp.Key) {
		if value := got.Issuer(); value != want {
			t.Error(format(msg, "OTP issuer: got %q, want %q", value, want))
		}
	}
}

func OTPPeriod(want uint64, msg ...string) Assertion[*otp.Key] {
	return func(t *testing.T, got *otp.Key) {
		if value := got.Period(); value != want {
			t.Error(format(msg, "OTP period: got %d, want %d", value, want))
		}
	}
}

func OTPType(want string, msg ...string) Assertion[*otp.Key] {
	return func(t *testing.T, got *otp.Key) {
		if value := got.Type(); value != want {
			t.Error(format(msg, "OTP type: got %q, want %q", value, want))
		}
	}
}

func OTPSecretValue(want string, msg ...string) Assertion[*otp.Key] {
	return func(t *testing.T, got *otp.Key) {
		if value := got.Secret(); value != want {
			t.Error(format(msg, "OTP secret value: got %q, want %q", value, want))
		}
	}
}

func OTPSecretLength(want int, msg ...string) Assertion[*otp.Key] {
	return func(t *testing.T, got *otp.Key) {
		if value := got.Secret(); len(value) != want {
			t.Error(format(msg, "OTP secret length: got %d, want %d", len(value), want))
		}
	}
}
