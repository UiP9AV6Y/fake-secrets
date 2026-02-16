package fake

import (
	"errors"
	"io"
	"io/fs"
	"log/slog"
	nethttp "net/http"

	"github.com/UiP9AV6Y/fake-secrets/internal/http"
)

var ErrNotAFile = errors.New("requested secret is not backed by a file")

type FileHandler struct {
	logger *slog.Logger
	root   fs.FS
}

func NewFileHandler(root fs.FS, logger *slog.Logger) *FileHandler {
	result := &FileHandler{
		logger: logger,
		root:   root,
	}

	return result
}

func (h *FileHandler) ServeHTTP(w nethttp.ResponseWriter, r *nethttp.Request) {
	name := r.PathValue("filename")

	h.logger.Debug("serving secret file", "filename", name)

	f, err := h.root.Open(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.ServeError(w, nethttp.StatusNotFound, err)
		} else {
			http.ServeError(w, nethttp.StatusInternalServerError, err)
		}

		return
	}

	defer func() {
		_ = f.Close()
	}()

	info, err := f.Stat()
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	if !info.Mode().IsRegular() {
		http.ServeError(w, nethttp.StatusNotFound, ErrNotAFile)
		return
	}

	meta := NewFileMeta(info, r)
	data, err := io.ReadAll(f)
	if err != nil {
		http.ServeError(w, nethttp.StatusInternalServerError, err)
		return
	}

	http.ServeSecret(w, data, meta)
}
