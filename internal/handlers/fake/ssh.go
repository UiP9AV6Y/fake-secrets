package fake

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/pem"
	"io"
	"log/slog"
	nethttp "net/http"

	"golang.org/x/crypto/ssh"

	"github.com/UiP9AV6Y/fake-secrets/internal/cache"
	"github.com/UiP9AV6Y/fake-secrets/internal/http"
)

type SSHHandler struct {
	logger  *slog.Logger
	rand    io.Reader
	rsa     cache.Cacher[*rsa.PrivateKey]
	ecdsa   cache.Cacher[*ecdsa.PrivateKey]
	ed25519 cache.Cacher[ed25519.PrivateKey]
}

func NewSSHHandler(rnd io.Reader, logger *slog.Logger) *SSHHandler {
	rsa := cache.NewCacher[*rsa.PrivateKey]()
	ecdsa := cache.NewCacher[*ecdsa.PrivateKey]()
	ed25519 := cache.NewCacher[ed25519.PrivateKey]()
	result := &SSHHandler{
		logger:  logger,
		rand:    rnd,
		rsa:     rsa,
		ecdsa:   ecdsa,
		ed25519: ed25519,
	}

	return result
}

func (h *SSHHandler) RSACache() cache.Cacher[*rsa.PrivateKey] {
	return h.rsa
}

func (h *SSHHandler) ECDSACache() cache.Cacher[*ecdsa.PrivateKey] {
	return h.ecdsa
}

func (h *SSHHandler) ED25519Cache() cache.Cacher[ed25519.PrivateKey] {
	return h.ed25519
}

func (h *SSHHandler) ServeCertificate(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("hostname")
	meta, err := ParseSSHMeta(name, r)
	if err != nil {
		http.ServeError(w, nethttp.StatusBadRequest, err)
		return
	}

	h.logger.Debug("serving generated SSH certificate", "meta", meta)

	key, _, err := LoadHandlerKey(h, &meta.CryptoMeta, h.rand)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	cert, err := ssh.NewPublicKey(key)
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

	h.logger.Debug("serving generated SSH certificate", "meta", meta)

	_, key, err := LoadHandlerKey(h, &meta.CryptoMeta, h.rand)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	block, err := ssh.MarshalPrivateKey(key, "")
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	data := pem.EncodeToMemory(block)

	http.ServeSecret(w, data, meta)
}
