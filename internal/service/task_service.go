package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"ssubench/internal/model"
	"ssubench/internal/repo"
)

var ErrTaskNotFound = errors.New("task not found")
var ErrInvalidTaskStatus = errors.New("task is not in a state that allows this action")
var ErrNotTaskOwner = errors.New("you are not the owner of this task")
var ErrBidNotFound = errors.New("bid not found")
var ErrNotAssignedExecutor = errors.New("you are not the assigned executor for this task")
var ErrInsufficientBalance = errors.New("customer has insufficient balance")

type TaskService struct {
	taskRepo    *repo.TaskRepo
	bidRepo     *repo.BidRepo
	userRepo    *repo.UserRepo
	paymentRepo *repo.PaymentRepo
	pool        *pgxpool.Pool
}

func NewTaskService(tr *repo.TaskRepo, br *repo.BidRepo, ur *repo.UserRepo, pr *repo.PaymentRepo, pool *pgxpool.Pool) *TaskService {
	return &TaskService{
		taskRepo:    tr,
		bidRepo:     br,
		userRepo:    ur,
		paymentRepo: pr,
		pool:        pool,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, customerID int, req model.CreateTaskRequest) (*model.Task, error) {
	task := &model.Task{
		CustomerID:  customerID,
		Title:       req.Title,
		Description: req.Description,
		Budget:      req.Budget,
		Status:      model.TaskStatusOpen,
	}
	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) ListTasks(ctx context.Context, filter model.TaskFilter) ([]model.Task, error) {
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	return s.taskRepo.List(ctx, filter)
}

func (s *TaskService) CreateBid(ctx context.Context, executorID int, req model.CreateBidRequest) (*model.Bid, error) {
	task, err := s.taskRepo.GetByID(ctx, req.TaskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	if task.Status != model.TaskStatusOpen {
		return nil, ErrInvalidTaskStatus
	}

	bid := &model.Bid{
		TaskID:     req.TaskID,
		ExecutorID: executorID,
		Amount:     req.Amount,
		Status:     model.BidStatusPending,
	}
	if err := s.bidRepo.Create(ctx, bid); err != nil {
		return nil, err
	}
	return bid, nil
}

func (s *TaskService) AcceptBid(ctx context.Context, customerID int, taskID int, bidID int) (*model.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	if task.CustomerID != customerID {
		return nil, ErrNotTaskOwner
	}

	if task.Status != model.TaskStatusOpen {
		return nil, ErrInvalidTaskStatus
	}

	bid, err := s.bidRepo.GetByID(ctx, bidID)
	if err != nil {
		return nil, ErrBidNotFound
	}

	if bid.TaskID != taskID {
		return nil, ErrBidNotFound
	}

	if bid.Status != model.BidStatusPending {
		return nil, ErrInvalidTaskStatus
	}

	if err := s.bidRepo.UpdateStatus(ctx, bidID, model.BidStatusAccepted); err != nil {
		return nil, err
	}

	if err := s.bidRepo.SetAllBidsStatus(ctx, taskID, model.BidStatusRejected); err != nil {
		return nil, err
	}

	if err := s.bidRepo.UpdateStatus(ctx, bidID, model.BidStatusAccepted); err != nil {
		return nil, err
	}

	if err := s.taskRepo.UpdateStatus(ctx, taskID, model.TaskStatusInProgress); err != nil {
		return nil, err
	}

	task.Status = model.TaskStatusInProgress
	return task, nil
}

func (s *TaskService) MarkAsCompleted(ctx context.Context, executorID int, taskID int) (*model.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	if task.Status != model.TaskStatusInProgress {
		return nil, ErrInvalidTaskStatus
	}

	acceptedBid, err := s.bidRepo.GetAcceptedBidByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if acceptedBid == nil || acceptedBid.ExecutorID != executorID {
		return nil, ErrNotAssignedExecutor
	}

	if err := s.taskRepo.UpdateStatus(ctx, taskID, model.TaskStatusCompleted); err != nil {
		return nil, err
	}

	task.Status = model.TaskStatusCompleted
	return task, nil
}

func (s *TaskService) ConfirmCompletion(ctx context.Context, customerID int, taskID int) (*model.Payment, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	if task.CustomerID != customerID {
		return nil, ErrNotTaskOwner
	}

	if task.Status != model.TaskStatusCompleted {
		return nil, ErrInvalidTaskStatus
	}

	acceptedBid, err := s.bidRepo.GetAcceptedBidByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if acceptedBid == nil {
		return nil, ErrBidNotFound
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	customer, err := s.userRepo.GetByIDInTx(ctx, tx, customerID)
	if err != nil {
		return nil, err
	}
	if customer.Balance < task.Budget {
		return nil, ErrInsufficientBalance
	}

	if err := s.userRepo.UpdateBalanceInTx(ctx, tx, customerID, -task.Budget); err != nil {
		return nil, err
	}

	if err := s.userRepo.UpdateBalanceInTx(ctx, tx, acceptedBid.ExecutorID, task.Budget); err != nil {
		return nil, err
	}

	payment := &model.Payment{
		TaskID:     taskID,
		FromUserID: customerID,
		ToUserID:   acceptedBid.ExecutorID,
		Amount:     task.Budget,
	}
	if err := s.paymentRepo.CreateInTx(ctx, tx, payment); err != nil {
		return nil, err
	}

	if err := s.taskRepo.UpdateStatusInTx(ctx, tx, taskID, model.TaskStatusConfirmed); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return payment, nil
}
