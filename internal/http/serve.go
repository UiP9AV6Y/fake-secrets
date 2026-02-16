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
	dto := map[string]interface{}{
		"error": err.Error(),
	}

	w.WriteHeader(code)
	ServeJSON(w, dto)
}

func ServeSecret(w nethttp.ResponseWriter, data []byte, meta interface{}) {
	secret := bytes.TrimSpace(data)
	dto := map[string]interface{}{
		"secret": string(secret),
		"meta":   meta,
	}

	w.WriteHeader(nethttp.StatusOK)
	ServeJSON(w, dto)
}

func ServeJSON(w nethttp.ResponseWriter, dto interface{}) {
	w.Header().Set(HeaderContentType, ContentTypeJSON+"; charset=utf-8")

	_ = json.NewEncoder(w).Encode(dto)
}
