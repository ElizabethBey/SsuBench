package repo

import (
	"context"
	"ssubench/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
