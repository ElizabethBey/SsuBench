package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"ssubench/internal/model"
	"ssubench/internal/service"
	"ssubench/internal/testingutils"
	"testing"
)

// Тест 11: HTTP 402 при недостаточном балансе для подтверждения
func TestTaskHandler_ConfirmCompletion_PaymentRequired(t *testing.T) {
	taskRepo := testingutils.NewMockTaskRepo()
	bidRepo := testingutils.NewMockBidRepo()
	userRepo := testingutils.NewMockUserRepo()
	paymentRepo := testingutils.NewMockPaymentRepo()
	pool := testingutils.NewMockPool()

	taskSvc := service.NewTaskService(taskRepo, bidRepo, userRepo, paymentRepo, pool)
	handler := NewTaskHandler(taskSvc)

	// Создаём заказчика с недостаточным балансом
	customer := &model.User{Email: "customer@test.com", Role: model.RoleCustomer, Balance: 50}
	executor := &model.User{Email: "executor@test.com", Role: model.RoleExecutor}
	userRepo.Create(context.Background(), customer)
	userRepo.Create(context.Background(), executor)

	task := &model.Task{CustomerID: customer.ID, Title: "Test Task", Budget: 100}
	taskRepo.Create(context.Background(), task)

	bid := &model.Bid{TaskID: task.ID, ExecutorID: executor.ID, Amount: 100}
	bidRepo.Create(context.Background(), bid)

	taskSvc.AcceptBid(context.Background(), customer.ID, task.ID, bid.ID)
	taskSvc.MarkAsCompleted(context.Background(), executor.ID, task.ID)

	// HTTP запрос на подтверждение
	req := httptest.NewRequest(http.MethodPost, "/tasks/1/confirm", nil)
	req.SetPathValue("id", "1")

	// Добавляем customerID в контекст
	ctx := context.WithValue(req.Context(), model.UserIDKey, customer.ID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.ConfirmCompletion(rr, req)

	if rr.Code != http.StatusPaymentRequired {
		t.Fatalf("status=%d want=%d", rr.Code, http.StatusPaymentRequired)
	}
}

// Тест 12: HTTP 403 при попытке принять заявку не владельцем
func TestTaskHandler_AcceptBid_Forbidden(t *testing.T) {
	taskRepo := testingutils.NewMockTaskRepo()
	bidRepo := testingutils.NewMockBidRepo()
	userRepo := testingutils.NewMockUserRepo()
	paymentRepo := testingutils.NewMockPaymentRepo()
	pool := testingutils.NewMockPool()

	taskSvc := service.NewTaskService(taskRepo, bidRepo, userRepo, paymentRepo, pool)
	handler := NewTaskHandler(taskSvc)

	customer := &model.User{Email: "customer@test.com", Role: model.RoleCustomer}
	executor := &model.User{Email: "executor@test.com", Role: model.RoleExecutor}
	hacker := &model.User{Email: "hacker@test.com", Role: model.RoleCustomer}
	userRepo.Create(context.Background(), customer)
	userRepo.Create(context.Background(), executor)
	userRepo.Create(context.Background(), hacker)

	task := &model.Task{CustomerID: customer.ID, Title: "Test Task", Budget: 100}
	taskRepo.Create(context.Background(), task)

	bid := &model.Bid{TaskID: task.ID, ExecutorID: executor.ID, Amount: 100}
	bidRepo.Create(context.Background(), bid)

	// HTTP запрос на принятие заявки от постороннего пользователя
	body, _ := json.Marshal(AcceptBidRequest{BidID: bid.ID})
	req := httptest.NewRequest(http.MethodPost, "/tasks/1/accept", bytes.NewReader(body))
	req.SetPathValue("id", "1")

	// Добавляем hacker ID в контекст (не владелец задачи)
	ctx := context.WithValue(req.Context(), model.UserIDKey, hacker.ID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.AcceptBid(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("status=%d want=%d", rr.Code, http.StatusForbidden)
	}
}
