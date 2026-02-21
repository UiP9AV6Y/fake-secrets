package http

import (
	nethttp "net/http"
	"strconv"

	"github.com/UiP9AV6Y/fake-secrets/internal/crypto"
)

func ParseFormString(r *nethttp.Request, field, fallback string) string {
	value := r.FormValue(field)
	if value == "" {
		return fallback
	}

	return value
}

func ParseFormBool(r *nethttp.Request, field string, fallback bool) (bool, error) {
	value := r.FormValue(field)
	if value == "" {
		return fallback, nil
	}

	result, err := strconv.ParseBool(value)
	if err != nil {
		return fallback, err
	}

	return result, nil
}

func ParseFormInt(r *nethttp.Request, field string, fallback int64) (int64, error) {
	value := r.FormValue(field)
	if value == "" {
		return fallback, nil
	}

	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil || result <= 0 {
		return fallback, err
	}

	return result, nil
}

func ParseFormCryptoAlgorithm(r *nethttp.Request, field string, fallback crypto.Algorithm) (crypto.Algorithm, error) {
	value := r.FormValue(field)
	if value == "" {
		return fallback, nil
	}

	result, err := crypto.ParseAlgorithm(value)
	if err != nil {
		return fallback, err
	}

	return result, nil
}

func ParseFormECDSACurve(r *nethttp.Request, field string, fallback crypto.ECDSACurve) (crypto.ECDSACurve, error) {
	value := r.FormValue(field)
	if value == "" {
		return fallback, nil
	}

	result, err := crypto.ParseECDSACurve(value)
	if err != nil {
		return fallback, err
	}

	return result, nil
}
