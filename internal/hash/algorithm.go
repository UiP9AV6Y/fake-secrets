package hash

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pquerna/otp"
)

type Algorithm int

const (
	AlgorithmSHA1 Algorithm = 1 + iota
	AlgorithmSHA256
	AlgorithmSHA512
	AlgorithmMD5
)

var algorithmImpl = map[string]Algorithm{
	"SHA1":   AlgorithmSHA1,
	"SHA256": AlgorithmSHA256,
	"SHA512": AlgorithmSHA512,
	"MD5":    AlgorithmMD5,
}

var algorithmOTP = map[Algorithm]otp.Algorithm{
	AlgorithmSHA1:   otp.AlgorithmSHA1,
	AlgorithmSHA256: otp.AlgorithmSHA256,
	AlgorithmSHA512: otp.AlgorithmSHA512,
	AlgorithmMD5:    otp.AlgorithmMD5,
}

func ParseAlgorithm(a string) (algo Algorithm, err error) {
	if a == "" {
		algo = AlgorithmSHA1
		return
	}

	err = (&algo).UnmarshalText([]byte(a))

	return
}

func (a *Algorithm) UnmarshalText(text []byte) error {
	algo, ok := algorithmImpl[strings.ToUpper(string(text))]
	if !ok {
		return fmt.Errorf("invalid hash algorithm %q", text)
	}

	*a = algo

	return nil
}

func (a Algorithm) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a Algorithm) String() string {
	switch a {
	case AlgorithmSHA1:
		return "SHA1"
	case AlgorithmSHA256:
		return "SHA256"
	case AlgorithmSHA512:
		return "SHA512"
	case AlgorithmMD5:
		return "MD5"
	default:
		return "unknown crypto algorithm " + strconv.Itoa(int(a))
	}
}

func (a Algorithm) OTPAlgorithm() otp.Algorithm {
	return algorithmOTP[a]
}
