package httpserver

import (
	"log/slog"
	"net/http"
	"ssubench/internal/config"
	"ssubench/internal/repo"
	"time"
)

func New(cfg config.Config, logger *slog.Logger) *http.Server {
	mux := http.NewServeMux()

	// Layered wiring: repo -> service -> handler
	taskRepo, err := repo.NewPostgresRepo(cfg)

	// Middleware (минимум): логирование + recover
	handler := RecoverMiddleware(logger, LoggingMiddleware(logger, mux))

	return &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
