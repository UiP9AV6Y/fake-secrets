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
	hotp := fake.NewHOTPHandler(random, logger)
	totp := fake.NewTOTPHandler(random, logger)

	router.HandleFunc("/", index.ServeHTTP)
	router.HandleFunc(generator.RouteStatic(cfg))
	router.HandleFunc(generator.RoutePassword(cfg))
	router.HandleFunc(generator.RouteRandomToken(cfg))
	router.HandleFunc(generator.RouteSeededToken(cfg))
	router.HandleFunc(generator.RouteRandomAPIKey(cfg))
	router.HandleFunc(generator.RouteSeededAPIKey(cfg))
	router.HandleFunc(ssh.RouteCertificate(cfg))
	router.HandleFunc(ssh.RoutePrivateKey(cfg))
	router.HandleFunc(tls.RouteCertificate(cfg))
	router.HandleFunc(tls.RoutePrivateKey(cfg))
	router.HandleFunc(jwt.RouteCertificate(cfg))
	router.HandleFunc(jwt.RoutePrivateKey(cfg))
	router.HandleFunc(jwt.RouteToken(cfg))
	router.HandleFunc(hotp.RoutePrivateKey(cfg))
	router.HandleFunc(hotp.RouteCode(cfg))
	router.HandleFunc(totp.RoutePrivateKey(cfg))
	router.HandleFunc(totp.RouteCode(cfg))
	router.Handle(status.Route(cfg), status)

	if cfg.StorageDir != "" {
		root, err := os.OpenRoot(cfg.StorageDir)
		if err != nil {
			return nil, err
		}

		file := fake.NewFileHandler(root.FS(), logger)

		router.Handle(file.Route(cfg), file)
	}

	return router, nil
}
