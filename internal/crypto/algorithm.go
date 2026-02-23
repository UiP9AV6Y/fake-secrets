package crypto

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"fmt"
	"strconv"
	"strings"
)

type Algorithm int

const (
	AlgorithmRSA Algorithm = 1 + iota
	AlgorithmECDSA
	AlgorithmED25519
)

var algorithmImpl = map[string]Algorithm{
	"RSA":     AlgorithmRSA,
	"ECDSA":   AlgorithmECDSA,
	"EDDSA":   AlgorithmED25519,
	"ED25519": AlgorithmED25519,
}

func ParseAlgorithm(a string) (algo Algorithm, err error) {
	if a == "" {
		algo = AlgorithmRSA
		return
	}

	err = (&algo).UnmarshalText([]byte(a))

	return
}

func (a *Algorithm) UnmarshalText(text []byte) error {
	algo, ok := algorithmImpl[strings.ToUpper(string(text))]
	if !ok {
		return fmt.Errorf("invalid crypto algorithm %q", text)
	}

	*a = algo

	return nil
}

func (a Algorithm) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a Algorithm) String() string {
	switch a {
	case AlgorithmRSA:
		return "RSA"
	case AlgorithmECDSA:
		return "ECDSA"
	case AlgorithmED25519:
		return "ED25519"
	default:
		return "unknown crypto algorithm " + strconv.Itoa(int(a))
	}
}

func PublicKey(key interface{}) interface{} {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}
