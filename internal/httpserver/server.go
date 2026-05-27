package httpserver

import (
	"log/slog"
	"net/http"
	"ssubench/internal/config"
	"ssubench/internal/handler"
	"ssubench/internal/model"
	"ssubench/internal/repo"
	"ssubench/internal/service"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg config.Config, pool *pgxpool.Pool, logger *slog.Logger) *http.Server {
	mux := http.NewServeMux()

	// Layered wiring: repo -> service -> handlr
	userRepo := repo.NewUserRepo(pool)
	taskRepo := repo.NewTaskRepo(pool)
	bidRepo := repo.NewBidRepo(pool)

	authSvc := service.NewAuthService(userRepo, cfg)
	taskSvc := service.NewTaskService(taskRepo, bidRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	taskHandler := handler.NewTaskHandler(taskSvc)

	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	authMiddleware := AuthMiddleware(authSvc)
	roleCustomer := RoleMiddleware(model.RoleCustomer)
	roleExecutor := RoleMiddleware(model.RoleExecutor)

	mux.Handle("POST /tasks", authMiddleware(roleCustomer(http.HandlerFunc(taskHandler.Create))))
	mux.Handle("GET /tasks", authMiddleware(http.HandlerFunc(taskHandler.List)))
	mux.Handle("POST /bids", authMiddleware(roleExecutor(http.HandlerFunc(taskHandler.CreateBid))))

	// Middleware (минимум): логирование + recover
	handlr := RecoverMiddleware(logger, LoggingMiddleware(logger, mux))

	return &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handlr,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
