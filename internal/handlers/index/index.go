package index

import (
	"fmt"
	nethttp "net/http"

	"github.com/UiP9AV6Y/fake-secrets/internal/http"
)

func ServeHTTP(w nethttp.ResponseWriter, r *nethttp.Request) {
	err := fmt.Errorf("invalid service request %q", r.URL.Path)

	http.ServeError(w, nethttp.StatusNotFound, err)
}
