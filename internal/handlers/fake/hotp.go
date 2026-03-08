package fake

import (
	"io"
	"log/slog"
	nethttp "net/http"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"

	"github.com/UiP9AV6Y/fake-secrets/internal/cache"
	"github.com/UiP9AV6Y/fake-secrets/internal/config"
	"github.com/UiP9AV6Y/fake-secrets/internal/http"
)

type HOTPHandler struct {
	logger *slog.Logger
	rand   io.Reader
	keys   cache.Cacher[*otp.Key]
}

func NewHOTPHandler(rnd io.Reader, logger *slog.Logger) *HOTPHandler {
	keys := cache.NewCacher[*otp.Key]()
	result := &HOTPHandler{
		logger: logger,
		rand:   rnd,
		keys:   keys,
	}

	return result
}

func (h *HOTPHandler) RoutePrivateKey(cfg *config.Config) (string, nethttp.HandlerFunc) {
	return cfg.HandlerPattern("hotp", "{account}", "keys"), h.ServePrivateKey
}

func (h *HOTPHandler) ServePrivateKey(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("account")
	meta, err := ParseHOTPMeta(name, r)
	if err != nil {
		http.ServeError(w, nethttp.StatusBadRequest, err)
		return
	}

	h.logger.Debug("serving generated HOTP key", "meta", meta)

	req := &cache.HOTPLoader{
		Issuer:      meta.Organization,
		AccountName: meta.Subject,
		SecretSize:  uint(meta.Length),
		Algorithm:   meta.Algorithm,
		Random:      h.rand,
	}
	key, err := h.keys.Load(req)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	data := []byte(key.String())

	http.ServeSecret(w, data, meta)
}

func (h *HOTPHandler) RouteCode(cfg *config.Config) (string, nethttp.HandlerFunc) {
	return cfg.HandlerPattern("hotp", "{account}", "codes"), h.ServeCode
}

func (h *HOTPHandler) ServeCode(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("account")
	meta, err := ParseHOTPMeta(name, r)
	if err != nil {
		http.ServeError(w, nethttp.StatusBadRequest, err)
		return
	}

	h.logger.Debug("serving generated HOTP code", "meta", meta)

	req := &cache.HOTPLoader{
		Issuer:      meta.Organization,
		AccountName: meta.Subject,
		SecretSize:  uint(meta.Length),
		Algorithm:   meta.Algorithm,
		Random:      h.rand,
	}
	key, err := h.keys.Load(req)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	opts := hotp.ValidateOpts{
		Algorithm: meta.Algorithm.OTPAlgorithm(),
	}
	code, err := hotp.GenerateCodeCustom(key.Secret(), uint64(meta.Counter), opts)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	data := []byte(code)

	http.ServeSecret(w, data, meta)
}
