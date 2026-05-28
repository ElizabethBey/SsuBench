package service

import (
	"context"
	"errors"
	"ssubench/internal/config"
	"ssubench/internal/model"
	"ssubench/internal/testingutils"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// Тест 10: Попытка входа с неверными учетными данными
func TestLogin_InvalidCredentials(t *testing.T) {
	userRepo := testingutils.NewMockUserRepo()
	cfg := config.Config{JWTSecret: "test-secret"}
	svc := NewAuthService(userRepo, cfg)

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &model.User{
		Email:        "test@example.com",
		PasswordHash: string(hash),
		Role:         model.RoleCustomer,
	}
	userRepo.Create(context.Background(), user)

	req := model.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, err := svc.Login(context.Background(), req)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
