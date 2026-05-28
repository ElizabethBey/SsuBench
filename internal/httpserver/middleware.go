package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"ssubench/internal/handler"
	"ssubench/internal/model"
	"ssubench/internal/service"
	"strings"
	"time"
)

type statusCapturingResponseWriter struct {
	http.ResponseWriter
	status int
}

func AuthMiddleware(authSvc *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handler.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				handler.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid authorization format")
				return
			}

			tokenStr := parts[1]
			userID, role, err := authSvc.ValidateToken(tokenStr)
			if err != nil {
				handler.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), model.UserIDKey, userID)
			ctx = context.WithValue(ctx, model.UserRoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RoleMiddleware(allowedRoles ...model.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(model.UserRoleKey).(model.Role)
			if !ok {
				handler.WriteError(w, http.StatusForbidden, "forbidden", "user role not found in context")
				return
			}

			isAllowed := false
			for _, role := range allowedRoles {
				if role == userRole {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				handler.WriteError(w, http.StatusForbidden, "forbidden", "you don't have permission to access this resource")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RecoverMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered", "panic", err)
				handler.WriteError(w, http.StatusInternalServerError, "internal_error", "something went wrong")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusCapturingResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sw.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
