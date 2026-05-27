package repo

import (
	"context"
	"errors"
	"fmt"
	"ssubench/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepo(pool *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{pool: pool}
}

func (r *TaskRepo) Create(ctx context.Context, t *model.Task) error {
	query := `INSERT INTO tasks (customer_id, title, description, budget) 
              VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, t.CustomerID, t.Title, t.Description, t.Budget).Scan(&t.ID, &t.CreatedAt)
}

func (r *TaskRepo) GetByID(ctx context.Context, id int) (*model.Task, error) {
	query := `SELECT id, customer_id, title, description, budget, status, created_at FROM tasks WHERE id = $1`
	var t model.Task
	err := r.pool.QueryRow(ctx, query, id).Scan(&t.ID, &t.CustomerID, &t.Title, &t.Description, &t.Budget, &t.Status, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepo) List(ctx context.Context, filter model.TaskFilter) ([]model.Task, error) {
	query := `SELECT id, customer_id, title, description, budget, status, created_at FROM tasks`
	var args []interface{}

	if filter.Status != nil {
		query += ` WHERE status = $1`
		args = append(args, *filter.Status)
	}

	query += ` ORDER BY created_at DESC LIMIT $` +
		fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` +
		fmt.Sprintf("%d", len(args)+2)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.CustomerID, &t.Title, &t.Description, &t.Budget, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TaskRepo) UpdateStatus(ctx context.Context, id int, status model.TaskStatus) error {
	query := `UPDATE tasks SET status = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}
