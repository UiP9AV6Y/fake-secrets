package fake

import (
	"crypto/x509"
	"encoding/pem"
	"log/slog"
	"math/rand"
	"net"
	nethttp "net/http"
	"time"

	"github.com/UiP9AV6Y/fake-secrets/internal/cache"
	"github.com/UiP9AV6Y/fake-secrets/internal/crypto"
	"github.com/UiP9AV6Y/fake-secrets/internal/http"
)

type TLSHandler struct {
	logger  *slog.Logger
	rsa     *cache.RSACache
	ecdsa   *cache.ECDSACache
	ed25519 *cache.ED25519Cache
	cert    *cache.CertCache
	start   time.Time
}

func NewTLSHandler(start time.Time, rnd *rand.Rand, logger *slog.Logger) *TLSHandler {
	rsa := cache.NewRSACache(rnd)
	ecdsa := cache.NewECDSACache(rnd)
	ed25519 := cache.NewED25519Cache(rnd)
	cert := cache.NewCertCache(rnd, nil)
	result := &TLSHandler{
		logger:  logger,
		start:   start,
		cert:    cert,
		rsa:     rsa,
		ecdsa:   ecdsa,
		ed25519: ed25519,
	}

	return result
}

func (h *TLSHandler) RSACache() *cache.RSACache {
	return h.rsa
}

func (h *TLSHandler) ECDSACache() *cache.ECDSACache {
	return h.ecdsa
}

func (h *TLSHandler) ED25519Cache() *cache.ED25519Cache {
	return h.ed25519
}

func (h *TLSHandler) ServeCertificate(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("hostname")
	meta, err := ParseTLSMeta(name, h.start, r)
	if err != nil {
		http.ServeError(w, nethttp.StatusBadRequest, err)
		return
	}

	h.logger.Debug("serving generated TLS certificate", "meta", meta)

	_, key, err := LoadHandlerKey(h, &meta.CryptoMeta)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	template := &x509.Certificate{
		Subject:               meta.Subject(),
		NotBefore:             meta.NotBefore(),
		NotAfter:              meta.NotAfter(),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if meta.Algorithm == crypto.AlgorithmRSA {
		template.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
	} else {
		template.KeyUsage = x509.KeyUsageDigitalSignature
	}

	for _, h := range meta.SubjectAltNames() {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	der, err := h.cert.Load(template, key)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	block := &pem.Block{
		Type:    "CERTIFICATE",
		Headers: nil,
		Bytes:   der,
	}

	data := pem.EncodeToMemory(block)

	http.ServeSecret(w, data, meta)
}

func (h *TLSHandler) ServePrivateKey(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("hostname")
	meta, err := ParseTLSMeta(name, h.start, r)
	if err != nil {
		http.ServeError(w, nethttp.StatusBadRequest, err)
		return
	}

	h.logger.Debug("serving generated TLS certificate", "meta", meta)

	_, key, err := LoadHandlerKey(h, &meta.CryptoMeta)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	block := &pem.Block{
		Type:    "PRIVATE KEY",
		Headers: nil,
		Bytes:   der,
	}

	data := pem.EncodeToMemory(block)

	http.ServeSecret(w, data, meta)
}
