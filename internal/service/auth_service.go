package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"ssubench/internal/config"
	"ssubench/internal/model"
	"ssubench/internal/repo"
	"time"
)

var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrInvalidToken = errors.New("invalid token")

type AuthService struct {
	repo      *repo.UserRepo
	jwtSecret string
}

func NewAuthService(r *repo.UserRepo, cfg config.Config) *AuthService {
	return &AuthService{
		repo:      r,
		jwtSecret: cfg.JWTSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         req.Role,
		Balance:      0,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (string, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenStr string) (int, model.Role, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return 0, "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", ErrInvalidToken
	}

	userID := int(claims["user_id"].(float64))
	role := model.Role(claims["role"].(string))

	return userID, role, nil
}
