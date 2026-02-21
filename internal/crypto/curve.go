package crypto

import (
	"crypto/elliptic"
	"fmt"
	"strconv"
	"strings"
)

type ECDSACurve int

const (
	ECDSACurveP224 ECDSACurve = 1 + iota
	ECDSACurveP256
	ECDSACurveP384
	ECDSACurveP521
)

var ecdsaCurves = map[ECDSACurve]elliptic.Curve{
	ECDSACurveP224: elliptic.P224(),
	ECDSACurveP256: elliptic.P256(),
	ECDSACurveP384: elliptic.P384(),
	ECDSACurveP521: elliptic.P521(),
}

var ecdsaImpl = map[string]ECDSACurve{
	"P224":      ECDSACurveP224,
	"P-224":     ECDSACurveP224,
	"SECP224R1": ECDSACurveP224,
	"P256":      ECDSACurveP256,
	"P-256":     ECDSACurveP256,
	"SECP256R1": ECDSACurveP256,
	"P384":      ECDSACurveP384,
	"P-384":     ECDSACurveP384,
	"SECP384R1": ECDSACurveP384,
	"P521":      ECDSACurveP521,
	"P-521":     ECDSACurveP521,
	"SECP521R1": ECDSACurveP521,
}

func ParseECDSACurve(c string) (curve ECDSACurve, err error) {
	if c == "" {
		curve = ECDSACurveP256
		return
	}

	err = (&curve).UnmarshalText([]byte(c))

	return
}

func (c *ECDSACurve) UnmarshalText(text []byte) error {
	curve, ok := ecdsaImpl[strings.ToUpper(string(text))]
	if !ok {
		return fmt.Errorf("invalid ECDSA curve %q", text)
	}

	*c = curve

	return nil
}

func (c ECDSACurve) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

func (c ECDSACurve) String() string {
	if curve, ok := ecdsaCurves[c]; !ok {
		return "unknown ECDSA curve " + strconv.Itoa(int(c))
	} else {
		return curve.Params().Name
	}
}

func (c ECDSACurve) Curve() elliptic.Curve {
	return ecdsaCurves[c]
}
