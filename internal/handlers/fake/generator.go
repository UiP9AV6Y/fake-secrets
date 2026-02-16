package fake

import (
	"log/slog"
	"math/rand"
	nethttp "net/http"

	"github.com/UiP9AV6Y/fake-secrets/internal/http"
)

var (
	RandomPoolUpper   = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	RandomPoolLower   = []byte("abcdefghijklmnopqrstuvwxyz")
	RandomPoolNumeric = []byte("0123456789")
	RandomPoolSpecial = []byte("!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~")
)

type GeneratorHandler struct {
	logger *slog.Logger
	rnd    *rand.Rand
}

func NewGeneratorHandler(rnd *rand.Rand, logger *slog.Logger) *GeneratorHandler {
	result := &GeneratorHandler{
		logger: logger,
		rnd:    rnd,
	}

	return result
}

func (h *GeneratorHandler) ServeStatic(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("secret")
	meta := NewStaticMeta(r)

	h.logger.Debug("serving static secret")

	data := []byte(name)

	http.ServeSecret(w, data, meta)
}

func (h *GeneratorHandler) ServeToken(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("seed")
	meta := NewTokenMeta(name, r)

	var data []byte
	var err error
	if name == "" {
		h.logger.Debug("serving random token secret")
		data, err = generateRandomUUID(h.rnd)
	} else {
		h.logger.Debug("serving seeded token secret", "seed", name)
		data, err = generateSeededUUID(name)
	}

	if err != nil {
		http.ServeError(w, nethttp.StatusBadRequest, err)
		return
	}

	http.ServeSecret(w, data, meta)
}

func (h *GeneratorHandler) ServePassword(w nethttp.ResponseWriter, r *nethttp.Request) {
	meta, err := ParsePasswordMeta(r)
	if err != nil {
		http.ServeError(w, nethttp.StatusBadRequest, err)
		return
	}

	pool := generateRandomPool(meta.Upper, meta.Lower, meta.Numeric, meta.Special)

	h.logger.Debug("serving generated password", "length", meta.Length, "pool", len(pool))

	data := generatePassword(h.rnd, meta.Length, pool)

	http.ServeSecret(w, data, meta)
}
