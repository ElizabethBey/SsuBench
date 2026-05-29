package service

import (
	"context"
	"ssubench/internal/model"

	"github.com/jackc/pgx/v5"
)

type AuthUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context) ([]model.User, error)
	UpdateBlockStatus(ctx context.Context, userID int, blocked bool) error
	GetByIDInTx(ctx context.Context, tx pgx.Tx, id int) (*model.User, error)
	UpdateBalanceInTx(ctx context.Context, tx pgx.Tx, userID int, delta float64) error
}

type TaskRepository interface {
	Create(ctx context.Context, t *model.Task) error
	GetByID(ctx context.Context, id int) (*model.Task, error)
	List(ctx context.Context, filter model.TaskFilter) ([]model.Task, error)
	UpdateStatus(ctx context.Context, id int, status model.TaskStatus) error
	UpdateStatusInTx(ctx context.Context, tx pgx.Tx, id int, status model.TaskStatus) error
}

type BidRepository interface {
	Create(ctx context.Context, b *model.Bid) error
	GetByID(ctx context.Context, id int) (*model.Bid, error)
	UpdateStatus(ctx context.Context, id int, status model.BidStatus) error
	SetAllBidsStatus(ctx context.Context, taskID int, status model.BidStatus) error
	GetAcceptedBidByTaskID(ctx context.Context, taskID int) (*model.Bid, error)
}

type PaymentRepository interface {
	CreateInTx(ctx context.Context, tx pgx.Tx, p *model.Payment) error
}

type TxManager interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}
