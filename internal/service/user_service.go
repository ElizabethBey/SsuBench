package service

import (
	"context"
	"errors"
	"ssubench/internal/model"
)

var ErrCannotBlockAdmin = errors.New("cannot block an admin user")
var ErrUserAlreadyBlocked = errors.New("user is already blocked")
var ErrUserNotBlocked = errors.New("user is not blocked")

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) BlockUser(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Role == model.RoleAdmin {
		return ErrCannotBlockAdmin
	}
	if user.IsBlocked {
		return ErrUserAlreadyBlocked
	}

	return s.userRepo.UpdateBlockStatus(ctx, userID, true)
}

func (s *UserService) UnblockUser(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !user.IsBlocked {
		return ErrUserNotBlocked
	}

	return s.userRepo.UpdateBlockStatus(ctx, userID, false)
}

func (s *UserService) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.userRepo.List(ctx)
}
