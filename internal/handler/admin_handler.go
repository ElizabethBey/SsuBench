package handler

import (
	"errors"
	"net/http"
	"ssubench/internal/repo"
	"ssubench/internal/service"
	"strconv"
)

type AdminHandler struct {
	userSvc *service.UserService
}

func NewAdminHandler(userSvc *service.UserService) *AdminHandler {
	return &AdminHandler{userSvc: userSvc}
}

// BlockUser — блокировка пользователя администратором.
// POST /admin/users/{id}/block
func (h *AdminHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "bad_request", "invalid user id")
		return
	}

	if err := h.userSvc.BlockUser(r.Context(), userID); err != nil {
		switch {
		case errors.Is(err, service.ErrCannotBlockAdmin):
			WriteError(w, http.StatusForbidden, "forbidden", "cannot block an admin user")
		case errors.Is(err, service.ErrUserAlreadyBlocked):
			WriteError(w, http.StatusConflict, "already_blocked", "user is already blocked")
		case errors.Is(err, repo.ErrUserNotFound):
			WriteError(w, http.StatusNotFound, "not_found", "user not found")
		default:
			WriteError(w, http.StatusInternalServerError, "internal", err.Error())
		}
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "user blocked"})
}

// UnblockUser — разблокировка пользователя администратором.
// POST /admin/users/{id}/unblock
func (h *AdminHandler) UnblockUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "bad_request", "invalid user id")
		return
	}

	if err := h.userSvc.UnblockUser(r.Context(), userID); err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotBlocked):
			WriteError(w, http.StatusConflict, "not_blocked", "user is not blocked")
		case errors.Is(err, repo.ErrUserNotFound):
			WriteError(w, http.StatusNotFound, "not_found", "user not found")
		default:
			WriteError(w, http.StatusInternalServerError, "internal", err.Error())
		}
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "user unblocked"})
}

// ListUsers — список всех пользователей (только для администратора).
// GET /admin/users
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userSvc.ListUsers(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal", err.Error())
		return
	}

	// Скрываем хеши паролей (на всякий случай)
	for i := range users {
		users[i].PasswordHash = ""
	}

	WriteJSON(w, http.StatusOK, users)
}
