package service

import (
	"context"
	"errors"
	"ssubench/internal/model"
	"ssubench/internal/repo"
)

var ErrTaskNotFound = errors.New("task not found")
var ErrInvalidTaskStatus = errors.New("task is not in a state that allows this action")
var ErrNotTaskOwner = errors.New("you are not the owner of this task")

type TaskService struct {
	taskRepo *repo.TaskRepo
	bidRepo  *repo.BidRepo
}

func NewTaskService(tr *repo.TaskRepo, br *repo.BidRepo) *TaskService {
	return &TaskService{taskRepo: tr, bidRepo: br}
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

func (s *TaskService) AcceptBid(ctx context.Context, customerID int, bidID int) error {
	// - Проверяем, что заказчик задачи == customerID
	// - Ставим отклику статус Accepted
	// - Всем остальным откликам на эту задачу ставим Rejected
	// - Ставим задаче статус InProgress

	return nil
}
