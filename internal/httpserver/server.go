package httpserver

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
	"ssubench/internal/config"
	"ssubench/internal/handler"
	"ssubench/internal/repo"
	"ssubench/internal/service"
	"time"
)

func New(cfg config.Config, pool *pgxpool.Pool, logger *slog.Logger) *http.Server {
	mux := http.NewServeMux()

	// Layered wiring: repo -> service -> handler
	userRepo := repo.NewUserRepo(pool)
	authSvc := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authSvc)

	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	// Middleware (минимум): логирование + recover
	handler := RecoverMiddleware(logger, LoggingMiddleware(logger, mux))

	return &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
