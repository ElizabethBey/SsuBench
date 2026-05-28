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
	paymentRepo := repo.NewPaymentRepo(pool)

	authSvc := service.NewAuthService(userRepo, cfg)
	taskSvc := service.NewTaskService(taskRepo, bidRepo, userRepo, paymentRepo, pool)
	userSvc := service.NewUserService(userRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	taskHandler := handler.NewTaskHandler(taskSvc)
	adminHandler := handler.NewAdminHandler(userSvc)

	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	authMiddleware := AuthMiddleware(authSvc)
	roleCustomer := RoleMiddleware(model.RoleCustomer)
	roleExecutor := RoleMiddleware(model.RoleExecutor)
	roleAdmin := RoleMiddleware(model.RoleAdmin)

	mux.Handle("POST /tasks", authMiddleware(roleCustomer(http.HandlerFunc(taskHandler.Create))))
	mux.Handle("GET /tasks", authMiddleware(http.HandlerFunc(taskHandler.List)))
	mux.Handle("POST /bids", authMiddleware(roleExecutor(http.HandlerFunc(taskHandler.CreateBid))))

	mux.Handle("POST /tasks/{id}/accept_bid", authMiddleware(roleCustomer(http.HandlerFunc(taskHandler.AcceptBid))))
	mux.Handle("POST /tasks/{id}/mark_completed", authMiddleware(roleExecutor(http.HandlerFunc(taskHandler.MarkCompleted))))
	mux.Handle("POST /tasks/{id}/confirm", authMiddleware(roleCustomer(http.HandlerFunc(taskHandler.ConfirmCompletion))))

	mux.Handle("POST /admin/users/{id}/block", authMiddleware(roleAdmin(http.HandlerFunc(adminHandler.BlockUser))))
	mux.Handle("POST /admin/users/{id}/unblock", authMiddleware(roleAdmin(http.HandlerFunc(adminHandler.UnblockUser))))
	mux.Handle("GET /admin/users", authMiddleware(roleAdmin(http.HandlerFunc(adminHandler.ListUsers))))

	// Middleware (минимум): логирование + recover
	handlr := RecoverMiddleware(logger, LoggingMiddleware(logger, mux))

	return &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handlr,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
