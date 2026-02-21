package fake

import (
	"encoding/pem"
	"log/slog"
	"math/rand"
	nethttp "net/http"

	"golang.org/x/crypto/ssh"

	"github.com/UiP9AV6Y/fake-secrets/internal/cache"
	"github.com/UiP9AV6Y/fake-secrets/internal/http"
)

type SSHHandler struct {
	logger  *slog.Logger
	rsa     *cache.RSACache
	ecdsa   *cache.ECDSACache
	ed25519 *cache.ED25519Cache
}

func NewSSHHandler(rnd *rand.Rand, logger *slog.Logger) *SSHHandler {
	rsa := cache.NewRSACache(rnd)
	ecdsa := cache.NewECDSACache(rnd)
	ed25519 := cache.NewED25519Cache(rnd)
	result := &SSHHandler{
		logger:  logger,
		rsa:     rsa,
		ecdsa:   ecdsa,
		ed25519: ed25519,
	}

	return result
}

func (h *SSHHandler) RSACache() *cache.RSACache {
	return h.rsa
}

func (h *SSHHandler) ECDSACache() *cache.ECDSACache {
	return h.ecdsa
}

func (h *SSHHandler) ED25519Cache() *cache.ED25519Cache {
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

	key, _, err := LoadHandlerKey(h, &meta.CryptoMeta)
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

	_, key, err := LoadHandlerKey(h, &meta.CryptoMeta)
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
