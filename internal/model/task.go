package model

import "time"

type TaskStatus string

const (
	TaskStatusOpen       TaskStatus = "open"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
	TaskStatusConfirmed  TaskStatus = "confirmed"
)

type Task struct {
	ID          int        `json:"id"`
	CustomerID  int        `json:"customer_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Budget      float64    `json:"budget"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}

type CreateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Budget      float64 `json:"budget"`
}

type TaskFilter struct {
	Status *TaskStatus `json:"status"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}
