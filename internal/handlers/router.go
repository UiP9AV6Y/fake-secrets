package handlers

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/UiP9AV6Y/fake-secrets/internal/config"
	"github.com/UiP9AV6Y/fake-secrets/internal/handlers/fake"
	"github.com/UiP9AV6Y/fake-secrets/internal/handlers/health"
	"github.com/UiP9AV6Y/fake-secrets/internal/handlers/index"
)

func NewRouter(cfg *config.Config, logger *slog.Logger) (http.Handler, error) {
	now := cfg.RandomSeedTime()
	router := http.NewServeMux()
	random := cfg.RandomGenerator()
	status := health.NewHandler(now, logger)
	generator := fake.NewGeneratorHandler(random, logger)
	ssh := fake.NewSSHHandler(random, logger)
	tls := fake.NewTLSHandler(now, random, logger)
	jwt := fake.NewJWTHandler(now, random, logger)

	router.HandleFunc("/", index.ServeHTTP)
	router.HandleFunc(cfg.HandlerPattern("passwords"), generator.ServePassword)
	router.HandleFunc(cfg.HandlerPattern("passwords", "{secret}"), generator.ServeStatic)
	router.HandleFunc(cfg.HandlerPattern("tokens"), generator.ServeToken)
	router.HandleFunc(cfg.HandlerPattern("tokens", "{seed}"), generator.ServeToken)
	router.HandleFunc(cfg.HandlerPattern("apikeys"), generator.ServeAPIKey)
	router.HandleFunc(cfg.HandlerPattern("apikeys", "{seed}"), generator.ServeAPIKey)
	router.HandleFunc(cfg.HandlerPattern("ssh", "{hostname}", "certificates"), ssh.ServeCertificate)
	router.HandleFunc(cfg.HandlerPattern("ssh", "{hostname}", "keys"), ssh.ServePrivateKey)
	router.HandleFunc(cfg.HandlerPattern("tls", "{hostname}", "certificates"), tls.ServeCertificate)
	router.HandleFunc(cfg.HandlerPattern("tls", "{hostname}", "keys"), tls.ServePrivateKey)
	router.HandleFunc(cfg.HandlerPattern("jwt", "{issuer}", "certificates"), jwt.ServeCertificate)
	router.HandleFunc(cfg.HandlerPattern("jwt", "{issuer}", "keys"), jwt.ServePrivateKey)
	router.HandleFunc(cfg.HandlerPattern("jwt", "{issuer}", "tokens"), jwt.ServeToken)
	router.Handle(cfg.HandlerPattern(health.URLPath), status)

	if cfg.StorageDir != "" {
		root, err := os.OpenRoot(cfg.StorageDir)
		if err != nil {
			return nil, err
		}

		file := fake.NewFileHandler(root.FS(), logger)

		router.Handle(cfg.HandlerPattern("files", "{filename}"), file)
	}

	return router, nil
}
