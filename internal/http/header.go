package http

import (
	nethttp "net/http"
)

const (
	HeaderXForwardedHost  = "X-Forwarded-Host"
	HeaderXForwardedProto = "X-Forwarded-Proto"
)

func ParseHeaderBaseURL(r *nethttp.Request) string {
	scheme := "http"
	host := r.Host

	if p := r.Header.Get(HeaderXForwardedProto); p != "" {
		scheme = p
	}

	if h := r.Header.Get(HeaderXForwardedHost); h != "" {
		host = h
	}

	return scheme + "://" + host
}
