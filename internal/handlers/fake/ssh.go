package fake

import (
	"crypto/x509"
	"encoding/pem"
	"log/slog"
	"math/rand"
	nethttp "net/http"

	"golang.org/x/crypto/ssh"

	"github.com/UiP9AV6Y/fake-secrets/internal/cache"
	"github.com/UiP9AV6Y/fake-secrets/internal/http"
)

type SSHHandler struct {
	logger *slog.Logger
	rsa    *cache.RSACache
}

func NewSSHHandler(rnd *rand.Rand, logger *slog.Logger) *SSHHandler {
	rsa := cache.NewRSACache(rnd)
	result := &SSHHandler{
		logger: logger,
		rsa:    rsa,
	}

	return result
}

func (h *SSHHandler) ServeCertificate(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("hostname")
	meta, err := ParseSSHMeta(name, r)
	if err != nil {
		http.ServeError(w, nethttp.StatusBadRequest, err)
		return
	}

	h.logger.Debug("serving generated SSH certificate", "length", meta.Length, "hostname", meta.Hostname)

	key, err := h.rsa.Load(meta.Hostname, meta.Length)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	cert, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	data := ssh.MarshalAuthorizedKey(cert)

	http.ServeSecret(w, data, meta)
}

func (h *SSHHandler) ServePrivateKey(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("hostname")
	meta, err := ParseSSHMeta(name, r)
	if err != nil {
		http.ServeError(w, nethttp.StatusBadRequest, err)
		return
	}

	h.logger.Debug("serving generated SSH certificate", "length", meta.Length, "hostname", meta.Hostname)

	key, err := h.rsa.Load(meta.Hostname, meta.Length)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	der := x509.MarshalPKCS1PrivateKey(key)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: der,
	}

	data := pem.EncodeToMemory(block)

	http.ServeSecret(w, data, meta)
}
