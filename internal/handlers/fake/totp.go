package fake

import (
	"io"
	"log/slog"
	nethttp "net/http"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/UiP9AV6Y/fake-secrets/internal/cache"
	"github.com/UiP9AV6Y/fake-secrets/internal/config"
	"github.com/UiP9AV6Y/fake-secrets/internal/http"
)

type TOTPHandler struct {
	logger *slog.Logger
	rand   io.Reader
	keys   cache.Cacher[*otp.Key]
}

func NewTOTPHandler(rnd io.Reader, logger *slog.Logger) *TOTPHandler {
	keys := cache.NewCacher[*otp.Key]()
	result := &TOTPHandler{
		logger: logger,
		rand:   rnd,
		keys:   keys,
	}

	return result
}

func (h *TOTPHandler) RoutePrivateKey(cfg *config.Config) (string, nethttp.HandlerFunc) {
	return cfg.HandlerPattern("totp", "{account}", "keys"), h.ServePrivateKey
}

func (h *TOTPHandler) ServePrivateKey(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("account")
	meta, err := ParseTOTPMeta(name, r)
	if err != nil {
		http.ServeError(w, nethttp.StatusBadRequest, err)
		return
	}

	h.logger.Debug("serving generated TOTP key", "meta", meta)

	req := &cache.TOTPLoader{
		Issuer:      meta.Organization,
		AccountName: meta.Subject,
		SecretSize:  uint(meta.Length),
		Algorithm:   meta.Algorithm,
		Period:      uint(meta.ValidFor),
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

func (h *TOTPHandler) RouteCode(cfg *config.Config) (string, nethttp.HandlerFunc) {
	return cfg.HandlerPattern("totp", "{account}", "codes"), h.ServeCode
}

func (h *TOTPHandler) ServeCode(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("account")
	meta, err := ParseTOTPMeta(name, r)
	if err != nil {
		http.ServeError(w, nethttp.StatusBadRequest, err)
		return
	}

	h.logger.Debug("serving generated TOTP code", "meta", meta)

	req := &cache.TOTPLoader{
		Issuer:      meta.Organization,
		AccountName: meta.Subject,
		SecretSize:  uint(meta.Length),
		Algorithm:   meta.Algorithm,
		Period:      uint(meta.ValidFor),
		Random:      h.rand,
	}
	key, err := h.keys.Load(req)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	opts := totp.ValidateOpts{
		Algorithm: meta.Algorithm.OTPAlgorithm(),
		Period:    uint(meta.ValidFor),
	}
	now := time.Unix(meta.ValidAt, 0)
	code, err := totp.GenerateCodeCustom(key.Secret(), now.UTC(), opts)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	data := []byte(code)

	http.ServeSecret(w, data, meta)
}
