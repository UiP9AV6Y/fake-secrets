package fake

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log/slog"
	"net"
	nethttp "net/http"
	"time"

	"github.com/UiP9AV6Y/fake-secrets/internal/cache"
	"github.com/UiP9AV6Y/fake-secrets/internal/crypto"
	"github.com/UiP9AV6Y/fake-secrets/internal/http"
)

type TLSHandler struct {
	logger  *slog.Logger
	start   time.Time
	rand    io.Reader
	rsa     cache.Cacher[*rsa.PrivateKey]
	ecdsa   cache.Cacher[*ecdsa.PrivateKey]
	ed25519 cache.Cacher[ed25519.PrivateKey]
	cert    cache.Cacher[crypto.Certificate]
}

func NewTLSHandler(start time.Time, rnd io.Reader, logger *slog.Logger) *TLSHandler {
	rsa := cache.NewCacher[*rsa.PrivateKey]()
	ecdsa := cache.NewCacher[*ecdsa.PrivateKey]()
	ed25519 := cache.NewCacher[ed25519.PrivateKey]()
	cert := cache.NewCacher[crypto.Certificate]()
	result := &TLSHandler{
		logger:  logger,
		start:   start,
		rand:    rnd,
		rsa:     rsa,
		ecdsa:   ecdsa,
		ed25519: ed25519,
		cert:    cert,
	}

	return result
}

func (h *TLSHandler) RSACache() cache.Cacher[*rsa.PrivateKey] {
	return h.rsa
}

func (h *TLSHandler) ECDSACache() cache.Cacher[*ecdsa.PrivateKey] {
	return h.ecdsa
}

func (h *TLSHandler) ED25519Cache() cache.Cacher[ed25519.PrivateKey] {
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

	_, key, err := LoadHandlerKey(h, &meta.CryptoMeta, h.rand)
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

	req := &cache.CertLoader{
		Template: template,
		Key:      key,
		Random:   h.rand,
	}
	der, err := h.cert.Load(req)
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

	_, key, err := LoadHandlerKey(h, &meta.CryptoMeta, h.rand)
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
