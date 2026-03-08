package fake_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

const defaultHostPort = "example.test"

type requestOption func(*requestBuilder)

func WithRequestPath(p string) requestOption {
	return func(rb *requestBuilder) {
		rb.path = append(rb.path, p)
	}
}

func WithRequestPathValue(p, v string) requestOption {
	return func(rb *requestBuilder) {
		if v == "" {
			return
		}

		rb.path = append(rb.path, v)
		rb.pathValues[p] = v
	}
}

func WithRequestHeader(k, v string) requestOption {
	return func(rb *requestBuilder) {
		rb.header.Add(k, v)
	}
}

func WithRequestQuery(k, v string) requestOption {
	return func(rb *requestBuilder) {
		rb.query.Add(k, v)
	}
}

type requestBuilder struct {
	hostPost   string
	scheme     string
	method     string
	header     http.Header
	query      url.Values
	path       []string
	pathValues map[string]string
}

func newRequestBuilder() *requestBuilder {
	result := &requestBuilder{
		hostPost:   defaultHostPort,
		scheme:     "http",
		method:     http.MethodGet,
		header:     http.Header{},
		query:      url.Values{},
		path:       []string{},
		pathValues: map[string]string{},
	}

	return result
}

func (rb *requestBuilder) URL() string {
	path := strings.Join(rb.path, "/")
	query := rb.query.Encode()

	return rb.scheme + "://" + rb.hostPost + "/" + path + "?" + query
}

func (rb *requestBuilder) DecorateRequest(req *http.Request) *http.Request {
	req.Host = rb.hostPost

	for k, v := range rb.pathValues {
		req.SetPathValue(k, v)
	}

	for k, vs := range rb.header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	return req
}

func newRequest(ctx context.Context, opts ...requestOption) *http.Request {
	builder := newRequestBuilder()

	for _, o := range opts {
		o(builder)
	}

	return builder.DecorateRequest(
		httptest.NewRequestWithContext(ctx, builder.method, builder.URL(), nil),
	)
}
