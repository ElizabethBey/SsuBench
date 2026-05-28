package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"ssubench/internal/model"
	"ssubench/internal/service"
	"strconv"
	"strings"
)

type TaskHandler struct {
	svc *service.TaskService
}

type AcceptBidRequest struct {
	BidID int `json:"bid_id"`
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	customerID := r.Context().Value(model.UserIDKey).(int)

	var req model.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "bad_request", "invalid json")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "title is required")
		return
	}
	if req.Budget <= 0 {
		WriteError(w, http.StatusBadRequest, "validation_error", "budget must be greater than 0")
		return
	}

	task, err := h.svc.CreateTask(r.Context(), customerID, req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := model.TaskFilter{
		Limit: 10,
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			filter.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		status := model.TaskStatus(statusStr)
		filter.Status = &status
	}

	tasks, err := h.svc.ListTasks(r.Context(), filter)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) CreateBid(w http.ResponseWriter, r *http.Request) {
	executorID := r.Context().Value(model.UserIDKey).(int)

	var req model.CreateBidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "bad_request", "invalid json")
		return
	}

	if req.Amount <= 0 {
		WriteError(w, http.StatusBadRequest, "validation_error", "amount must be greater than 0")
		return
	}

	bid, err := h.svc.CreateBid(r.Context(), executorID, req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, bid)
}

func (h *TaskHandler) AcceptBid(w http.ResponseWriter, r *http.Request) {
	customerID := r.Context().Value(model.UserIDKey).(int)

	taskIDStr := r.PathValue("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "bad_request", "invalid task id")
		return
	}

	var req AcceptBidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "bad_request", "invalid json")
		return
	}

	task, err := h.svc.AcceptBid(r.Context(), customerID, taskID, req.BidID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			WriteError(w, http.StatusNotFound, "not_found", "task not found")
		case errors.Is(err, service.ErrBidNotFound):
			WriteError(w, http.StatusNotFound, "not_found", "bid not found")
		case errors.Is(err, service.ErrNotTaskOwner):
			WriteError(w, http.StatusForbidden, "forbidden", "you are not the owner of this task")
		case errors.Is(err, service.ErrInvalidTaskStatus):
			WriteError(w, http.StatusConflict, "invalid_status", err.Error())
		default:
			WriteError(w, http.StatusInternalServerError, "internal", err.Error())
		}
		return
	}

	WriteJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) MarkCompleted(w http.ResponseWriter, r *http.Request) {
	executorID := r.Context().Value(model.UserIDKey).(int)

	taskIDStr := r.PathValue("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "bad_request", "invalid task id")
		return
	}

	task, err := h.svc.MarkAsCompleted(r.Context(), executorID, taskID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			WriteError(w, http.StatusNotFound, "not_found", "task not found")
		case errors.Is(err, service.ErrNotAssignedExecutor):
			WriteError(w, http.StatusForbidden, "forbidden", "you are not the assigned executor")
		case errors.Is(err, service.ErrInvalidTaskStatus):
			WriteError(w, http.StatusConflict, "invalid_status", err.Error())
		default:
			WriteError(w, http.StatusInternalServerError, "internal", err.Error())
		}
		return
	}

	WriteJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) ConfirmCompletion(w http.ResponseWriter, r *http.Request) {
	customerID := r.Context().Value(model.UserIDKey).(int)

	taskIDStr := r.PathValue("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "bad_request", "invalid task id")
		return
	}

	payment, err := h.svc.ConfirmCompletion(r.Context(), customerID, taskID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			WriteError(w, http.StatusNotFound, "not_found", "task not found")
		case errors.Is(err, service.ErrNotTaskOwner):
			WriteError(w, http.StatusForbidden, "forbidden", "you are not the owner of this task")
		case errors.Is(err, service.ErrInvalidTaskStatus):
			WriteError(w, http.StatusConflict, "invalid_status", err.Error())
		case errors.Is(err, service.ErrInsufficientBalance):
			WriteError(w, http.StatusPaymentRequired, "insufficient_balance", "customer has insufficient balance")
		default:
			WriteError(w, http.StatusInternalServerError, "internal", err.Error())
		}
		return
	}

	WriteJSON(w, http.StatusOK, payment)
}
