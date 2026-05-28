package service

import (
	"context"
	"errors"
	"ssubench/internal/model"
	"ssubench/internal/testingutils"
	"testing"
)

// Тест 7: Попытка заблокировать администратора
func TestBlockUser_CannotBlockAdmin(t *testing.T) {
	userRepo := testingutils.NewMockUserRepo()
	svc := NewUserService(userRepo)

	// Создаём администратора
	admin := &model.User{Email: "admin@test.com", Role: model.RoleAdmin}
	userRepo.Create(context.Background(), admin)

	// Попытка заблокировать администратора
	err := svc.BlockUser(context.Background(), admin.ID)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrCannotBlockAdmin) {
		t.Fatalf("expected ErrCannotBlockAdmin, got %v", err)
	}
}

// Тест 8: Попытка заблокировать уже заблокированного пользователя
func TestBlockUser_AlreadyBlocked(t *testing.T) {
	userRepo := testingutils.NewMockUserRepo()
	svc := NewUserService(userRepo)

	// Создаём пользователя
	user := &model.User{Email: "user@test.com", Role: model.RoleCustomer, IsBlocked: false}
	userRepo.Create(context.Background(), user)

	// Блокируем первый раз
	svc.BlockUser(context.Background(), user.ID)

	// Попытка заблокировать повторно
	err := svc.BlockUser(context.Background(), user.ID)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrUserAlreadyBlocked) {
		t.Fatalf("expected ErrUserAlreadyBlocked, got %v", err)
	}
}

// Тест 9: Попытка разблокировать незаблокированного пользователя
func TestUnblockUser_NotBlocked(t *testing.T) {
	userRepo := testingutils.NewMockUserRepo()
	svc := NewUserService(userRepo)

	// Создаём незаблокированного пользователя
	user := &model.User{Email: "user@test.com", Role: model.RoleCustomer, IsBlocked: false}
	userRepo.Create(context.Background(), user)

	// Попытка разблокировать незаблокированного
	err := svc.UnblockUser(context.Background(), user.ID)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrUserNotBlocked) {
		t.Fatalf("expected ErrUserNotBlocked, got %v", err)
	}
}
