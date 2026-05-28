package handler

import (
	"encoding/json"
	"net/http"
	"regexp"
	"ssubench/internal/model"
	"ssubench/internal/service"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "bad_request", "invalid json")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || !emailRegex.MatchString(req.Email) {
		WriteError(w, http.StatusBadRequest, "validation_error", "invalid email format")
		return
	}
	if len(req.Password) < 6 {
		WriteError(w, http.StatusBadRequest, "validation_error", "password must be at least 6 characters")
		return
	}
	if req.Role != model.RoleCustomer && req.Role != model.RoleExecutor && req.Role != model.RoleAdmin {
		WriteError(w, http.StatusBadRequest, "validation_error", "invalid role, must be customer, executor or admin")
		return
	}

	user, err := h.svc.Register(r.Context(), req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "registration_failed", err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "bad_request", "invalid json")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "email is required")
		return
	}
	if req.Password == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "password is required")
		return
	}

	token, err := h.svc.Login(r.Context(), req)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid email or password")
		return
	}

	WriteJSON(w, http.StatusOK, model.LoginResponse{Token: token})
}
