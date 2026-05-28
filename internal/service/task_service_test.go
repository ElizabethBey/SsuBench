package service

import (
	"context"
	"errors"
	"ssubench/internal/model"
	"ssubench/internal/testingutils"
	"testing"
)

// Тест 1: Попытка принять заявку не владельцем задачи
func TestAcceptBid_NotTaskOwner(t *testing.T) {
	taskRepo := testingutils.NewMockTaskRepo()
	bidRepo := testingutils.NewMockBidRepo()
	userRepo := testingutils.NewMockUserRepo()
	paymentRepo := testingutils.NewMockPaymentRepo()
	pool := testingutils.NewMockPool()

	svc := NewTaskService(taskRepo, bidRepo, userRepo, paymentRepo, pool)

	customer := &model.User{Email: "customer@test.com", Role: model.RoleCustomer, Balance: 1000}
	executor := &model.User{Email: "executor@test.com", Role: model.RoleExecutor}
	userRepo.Create(context.Background(), customer)
	userRepo.Create(context.Background(), executor)

	task := &model.Task{CustomerID: customer.ID, Title: "Test Task", Budget: 100}
	taskRepo.Create(context.Background(), task)

	bid := &model.Bid{TaskID: task.ID, ExecutorID: executor.ID, Amount: 100}
	bidRepo.Create(context.Background(), bid)

	// Попытка принять заявку другим пользователем
	_, err := svc.AcceptBid(context.Background(), 999, task.ID, bid.ID)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrNotTaskOwner) {
		t.Fatalf("expected ErrNotTaskOwner, got %v", err)
	}
}

// Тест 2: Попытка принять заявку когда задача не в статусе Open
func TestAcceptBid_InvalidStatus(t *testing.T) {
	taskRepo := testingutils.NewMockTaskRepo()
	bidRepo := testingutils.NewMockBidRepo()
	userRepo := testingutils.NewMockUserRepo()
	paymentRepo := testingutils.NewMockPaymentRepo()
	pool := testingutils.NewMockPool()

	svc := NewTaskService(taskRepo, bidRepo, userRepo, paymentRepo, pool)

	customer := &model.User{Email: "customer@test.com", Role: model.RoleCustomer}
	executor := &model.User{Email: "executor@test.com", Role: model.RoleExecutor}
	userRepo.Create(context.Background(), customer)
	userRepo.Create(context.Background(), executor)

	task := &model.Task{CustomerID: customer.ID, Title: "Test Task", Budget: 100}
	taskRepo.Create(context.Background(), task)

	taskRepo.UpdateStatus(context.Background(), task.ID, model.TaskStatusInProgress)

	bid := &model.Bid{TaskID: task.ID, ExecutorID: executor.ID, Amount: 100}
	bidRepo.Create(context.Background(), bid)

	// Попытка принять заявку когда задача уже не Open
	_, err := svc.AcceptBid(context.Background(), customer.ID, task.ID, bid.ID)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidTaskStatus) {
		t.Fatalf("expected ErrInvalidTaskStatus, got %v", err)
	}
}

// Тест 3: Попытка пометить задачу выполненной не назначенным исполнителем
func TestMarkAsCompleted_NotAssignedExecutor(t *testing.T) {
	taskRepo := testingutils.NewMockTaskRepo()
	bidRepo := testingutils.NewMockBidRepo()
	userRepo := testingutils.NewMockUserRepo()
	paymentRepo := testingutils.NewMockPaymentRepo()
	pool := testingutils.NewMockPool()

	svc := NewTaskService(taskRepo, bidRepo, userRepo, paymentRepo, pool)

	customer := &model.User{Email: "customer@test.com", Role: model.RoleCustomer}
	executor1 := &model.User{Email: "executor1@test.com", Role: model.RoleExecutor}
	executor2 := &model.User{Email: "executor2@test.com", Role: model.RoleExecutor}
	userRepo.Create(context.Background(), customer)
	userRepo.Create(context.Background(), executor1)
	userRepo.Create(context.Background(), executor2)

	task := &model.Task{CustomerID: customer.ID, Title: "Test Task", Budget: 100}
	taskRepo.Create(context.Background(), task)

	bid := &model.Bid{TaskID: task.ID, ExecutorID: executor1.ID, Amount: 100}
	bidRepo.Create(context.Background(), bid)

	svc.AcceptBid(context.Background(), customer.ID, task.ID, bid.ID)
	_, err := svc.MarkAsCompleted(context.Background(), executor2.ID, task.ID)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrNotAssignedExecutor) {
		t.Fatalf("expected ErrNotAssignedExecutor, got %v", err)
	}
}

// Тест 4: Попытка подтвердить выполнение не владельцем задачи
func TestConfirmCompletion_NotTaskOwner(t *testing.T) {
	taskRepo := testingutils.NewMockTaskRepo()
	bidRepo := testingutils.NewMockBidRepo()
	userRepo := testingutils.NewMockUserRepo()
	paymentRepo := testingutils.NewMockPaymentRepo()
	pool := testingutils.NewMockPool()

	svc := NewTaskService(taskRepo, bidRepo, userRepo, paymentRepo, pool)

	customer := &model.User{Email: "customer@test.com", Role: model.RoleCustomer, Balance: 1000}
	executor := &model.User{Email: "executor@test.com", Role: model.RoleExecutor}
	userRepo.Create(context.Background(), customer)
	userRepo.Create(context.Background(), executor)

	task := &model.Task{CustomerID: customer.ID, Title: "Test Task", Budget: 100}
	taskRepo.Create(context.Background(), task)

	bid := &model.Bid{TaskID: task.ID, ExecutorID: executor.ID, Amount: 100}
	bidRepo.Create(context.Background(), bid)

	svc.AcceptBid(context.Background(), customer.ID, task.ID, bid.ID)
	svc.MarkAsCompleted(context.Background(), executor.ID, task.ID)

	// Попытка подтвердить от другого пользователя
	_, err := svc.ConfirmCompletion(context.Background(), 999, task.ID)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrNotTaskOwner) {
		t.Fatalf("expected ErrNotTaskOwner, got %v", err)
	}
}

// Тест 5: Попытка подтвердить выполнение при недостаточном балансе
func TestConfirmCompletion_InsufficientBalance(t *testing.T) {
	taskRepo := testingutils.NewMockTaskRepo()
	bidRepo := testingutils.NewMockBidRepo()
	userRepo := testingutils.NewMockUserRepo()
	paymentRepo := testingutils.NewMockPaymentRepo()
	pool := testingutils.NewMockPool()

	svc := NewTaskService(taskRepo, bidRepo, userRepo, paymentRepo, pool)

	customer := &model.User{Email: "customer@test.com", Role: model.RoleCustomer, Balance: 50}
	executor := &model.User{Email: "executor@test.com", Role: model.RoleExecutor}
	userRepo.Create(context.Background(), customer)
	userRepo.Create(context.Background(), executor)

	task := &model.Task{CustomerID: customer.ID, Title: "Test Task", Budget: 100}
	taskRepo.Create(context.Background(), task)

	bid := &model.Bid{TaskID: task.ID, ExecutorID: executor.ID, Amount: 100}
	bidRepo.Create(context.Background(), bid)

	svc.AcceptBid(context.Background(), customer.ID, task.ID, bid.ID)
	svc.MarkAsCompleted(context.Background(), executor.ID, task.ID)

	// Попытка подтвердить при недостаточном балансе
	_, err := svc.ConfirmCompletion(context.Background(), customer.ID, task.ID)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrInsufficientBalance) {
		t.Fatalf("expected ErrInsufficientBalance, got %v", err)
	}
}

// Тест 6: Успешное подтверждение выполнения с переводом баланса
func TestConfirmCompletion_Success(t *testing.T) {
	taskRepo := testingutils.NewMockTaskRepo()
	bidRepo := testingutils.NewMockBidRepo()
	userRepo := testingutils.NewMockUserRepo()
	paymentRepo := testingutils.NewMockPaymentRepo()
	pool := testingutils.NewMockPool()

	svc := NewTaskService(taskRepo, bidRepo, userRepo, paymentRepo, pool)

	customer := &model.User{Email: "customer@test.com", Role: model.RoleCustomer, Balance: 1000}
	executor := &model.User{Email: "executor@test.com", Role: model.RoleExecutor, Balance: 0}
	userRepo.Create(context.Background(), customer)
	userRepo.Create(context.Background(), executor)

	task := &model.Task{CustomerID: customer.ID, Title: "Test Task", Budget: 100}
	taskRepo.Create(context.Background(), task)

	bid := &model.Bid{TaskID: task.ID, ExecutorID: executor.ID, Amount: 100}
	bidRepo.Create(context.Background(), bid)

	svc.AcceptBid(context.Background(), customer.ID, task.ID, bid.ID)
	svc.MarkAsCompleted(context.Background(), executor.ID, task.ID)

	payment, err := svc.ConfirmCompletion(context.Background(), customer.ID, task.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if payment == nil {
		t.Fatalf("expected payment, got nil")
	}
	if payment.Amount != 100 {
		t.Fatalf("payment amount=%f want=%f", payment.Amount, 100.0)
	}

	customerUpdated, _ := userRepo.GetByID(context.Background(), customer.ID)
	executorUpdated, _ := userRepo.GetByID(context.Background(), executor.ID)

	if customerUpdated.Balance != 900 {
		t.Fatalf("customer balance=%f want=%f", customerUpdated.Balance, 900.0)
	}
	if executorUpdated.Balance != 100 {
		t.Fatalf("executor balance=%f want=%f", executorUpdated.Balance, 100.0)
	}

	taskUpdated, _ := taskRepo.GetByID(context.Background(), task.ID)
	if taskUpdated.Status != model.TaskStatusConfirmed {
		t.Fatalf("task status=%s want=%s", taskUpdated.Status, model.TaskStatusConfirmed)
	}
}
