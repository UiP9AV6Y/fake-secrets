package http

import (
	"bytes"
	"encoding/json"
	nethttp "net/http"
)

const (
	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json"
)

func ServeError(w nethttp.ResponseWriter, code int, err error) {
	dto := map[string]any{
		"error": err.Error(),
	}

	w.WriteHeader(code)
	ServeJSON(w, dto)
}

func ServeSecretObject(w nethttp.ResponseWriter, data, meta any) {
	secret, err := json.Marshal(data)
	if err != nil {
		ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	ServeSecret(w, secret, meta)
}

func ServeSecret(w nethttp.ResponseWriter, data []byte, meta any) {
	secret := bytes.TrimSpace(data)
	dto := map[string]any{
		"secret": string(secret),
		"meta":   meta,
	}

	w.WriteHeader(nethttp.StatusOK)
	ServeJSON(w, dto)
}

func ServeJSON(w nethttp.ResponseWriter, dto any) {
	w.Header().Set(HeaderContentType, ContentTypeJSON+"; charset=utf-8")

	_ = json.NewEncoder(w).Encode(dto)
}
