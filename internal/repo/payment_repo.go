package repo

import (
	"context"
	"ssubench/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepo struct {
	pool *pgxpool.Pool
}

func NewPaymentRepo(pool *pgxpool.Pool) *PaymentRepo {
	return &PaymentRepo{pool: pool}
}

func (r *PaymentRepo) CreateInTx(ctx context.Context, tx pgx.Tx, p *model.Payment) error {
	query := `INSERT INTO payments (task_id, from_user_id, to_user_id, amount)
	          VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return tx.QueryRow(ctx, query, p.TaskID, p.FromUserID, p.ToUserID, p.Amount).Scan(&p.ID, &p.CreatedAt)
}
