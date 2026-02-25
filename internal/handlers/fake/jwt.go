package fake

import (
	"log/slog"
	"math/rand"
	nethttp "net/http"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"

	"github.com/UiP9AV6Y/fake-secrets/internal/cache"
	"github.com/UiP9AV6Y/fake-secrets/internal/http"
)

type JWTHandler struct {
	logger  *slog.Logger
	rsa     *cache.RSACache
	ecdsa   *cache.ECDSACache
	ed25519 *cache.ED25519Cache
	start   time.Time
}

func NewJWTHandler(start time.Time, rnd *rand.Rand, logger *slog.Logger) *JWTHandler {
	rsa := cache.NewRSACache(rnd)
	ecdsa := cache.NewECDSACache(rnd)
	ed25519 := cache.NewED25519Cache(rnd)
	result := &JWTHandler{
		logger:  logger,
		start:   start,
		rsa:     rsa,
		ecdsa:   ecdsa,
		ed25519: ed25519,
	}

	return result
}

func (h *JWTHandler) RSACache() *cache.RSACache {
	return h.rsa
}

func (h *JWTHandler) ECDSACache() *cache.ECDSACache {
	return h.ecdsa
}

func (h *JWTHandler) ED25519Cache() *cache.ED25519Cache {
	return h.ed25519
}

func (h *JWTHandler) ServeToken(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("issuer")
	meta, err := ParseJWTMeta(name, h.start, r)
	if err != nil {
		http.ServeError(w, nethttp.StatusBadRequest, err)
		return
	}

	h.logger.Debug("serving generated JWT token", "meta", meta)

	_, key, err := LoadHandlerKey(h, &meta.CryptoMeta)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	k, err := jwk.Import(key)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	if err := jwk.AssignKeyID(k); err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	b := jwt.NewBuilder().
		Expiration(meta.ExpirationClaim()).
		NotBefore(meta.NotBeforeClaim()).
		IssuedAt(meta.IssuedAtClaim()).
		Issuer(meta.IssuerClaim())

	if meta.Audience != "" {
		b.Audience([]string{meta.Audience})
	}

	if meta.Subject != "" {
		b.Subject(meta.Subject)
	}

	token, err := b.Build()
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	alg := meta.SignatureAlgorithm()
	data, err := jwt.Sign(token, jwt.WithKey(alg, k))
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	http.ServeSecret(w, data, meta)
}

func (h *JWTHandler) ServeCertificate(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("issuer")
	meta, err := ParseJWTMeta(name, h.start, r)
	if err != nil {
		http.ServeError(w, nethttp.StatusBadRequest, err)
		return
	}

	h.logger.Debug("serving generated JWK public keyset", "meta", meta)

	key, _, err := LoadHandlerKey(h, &meta.CryptoMeta)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	k, err := jwk.Import(key)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	serveJWTKeySet(w, k, meta)
}

func (h *JWTHandler) ServePrivateKey(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("issuer")
	meta, err := ParseJWTMeta(name, h.start, r)
	if err != nil {
		http.ServeError(w, nethttp.StatusBadRequest, err)
		return
	}

	h.logger.Debug("serving generated JWK private keyset", "meta", meta)

	_, key, err := LoadHandlerKey(h, &meta.CryptoMeta)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	k, err := jwk.Import(key)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	serveJWTKeySet(w, k, meta)
}

func serveJWTKeySet(w nethttp.ResponseWriter, key jwk.Key, meta interface{}) {
	if err := jwk.AssignKeyID(key); err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	data := jwk.NewSet()
	if err := data.AddKey(key); err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	http.ServeSecretObject(w, data, meta)
}
