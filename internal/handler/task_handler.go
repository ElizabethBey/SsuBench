package handler

import (
	"encoding/json"
	"net/http"
	"ssubench/internal/model"
	"ssubench/internal/service"
)

type TaskHandler struct {
	svc *service.TaskService
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
	// Тут можно добавить парсинг query параметров для фильтрации и пагинации
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

	bid, err := h.svc.CreateBid(r.Context(), executorID, req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, bid)
}
