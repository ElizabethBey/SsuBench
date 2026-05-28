package repo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"ssubench/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrBidNotFound = errors.New("bid not found")

type BidRepo struct {
	pool *pgxpool.Pool
}

func NewBidRepo(pool *pgxpool.Pool) *BidRepo {
	return &BidRepo{pool: pool}
}

func (r *BidRepo) Create(ctx context.Context, b *model.Bid) error {
	query := `INSERT INTO bids (task_id, executor_id, amount) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, b.TaskID, b.ExecutorID, b.Amount).Scan(&b.ID, &b.CreatedAt)
}

func (r *BidRepo) UpdateStatus(ctx context.Context, id int, status model.BidStatus) error {
	query := `UPDATE bids SET status = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}

func (r *BidRepo) SetAllBidsStatus(ctx context.Context, taskID int, status model.BidStatus) error {
	query := `UPDATE bids SET status = $1 WHERE task_id = $2`
	_, err := r.pool.Exec(ctx, query, status, taskID)
	return err
}

func (r *BidRepo) GetByID(ctx context.Context, id int) (*model.Bid, error) {
	query := `SELECT id, task_id, executor_id, amount, status, created_at FROM bids WHERE id = $1`
	var b model.Bid
	err := r.pool.QueryRow(ctx, query, id).Scan(&b.ID, &b.TaskID, &b.ExecutorID, &b.Amount, &b.Status, &b.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBidNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (r *BidRepo) GetAcceptedBidByTaskID(ctx context.Context, taskID int) (*model.Bid, error) {
	query := `SELECT id, task_id, executor_id, amount, status, created_at FROM bids 
	          WHERE task_id = $1 AND status = $2 LIMIT 1`
	var b model.Bid
	err := r.pool.QueryRow(ctx, query, taskID, model.BidStatusAccepted).Scan(&b.ID, &b.TaskID, &b.ExecutorID, &b.Amount, &b.Status, &b.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}
