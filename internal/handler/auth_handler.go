package handler

import (
	"encoding/json"
	"net/http"
	"ssubench/internal/model"
	"ssubench/internal/service"
)

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

	token, err := h.svc.Login(r.Context(), req)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid email or password")
		return
	}

	WriteJSON(w, http.StatusOK, model.LoginResponse{Token: token})
}
