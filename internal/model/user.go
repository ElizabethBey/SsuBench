package model

import "time"

type ContextKey string

const (
	UserIDKey   ContextKey = "user_id"
	UserRoleKey ContextKey = "user_role"
)

type Role string

const (
	RoleCustomer Role = "customer"
	RoleExecutor Role = "executor"
	RoleAdmin    Role = "admin"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	Balance      float64   `json:"balance"`
	IsBlocked    bool      `json:"is_blocked"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
