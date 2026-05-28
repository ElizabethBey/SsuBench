package repo

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"ssubench/internal/model"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")
var pgErr *pgconn.PgError

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (email, password_hash, role, balance) 
              VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	err := r.pool.QueryRow(ctx, query, user.Email, user.PasswordHash, user.Role, user.Balance).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, password_hash, role, balance, is_blocked, created_at FROM users WHERE email = $1`

	var u model.User
	err := r.pool.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.Balance, &u.IsBlocked, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `SELECT id, email, password_hash, role, balance, is_blocked, created_at FROM users WHERE id = $1`

	var u model.User
	err := r.pool.QueryRow(ctx, query, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.Balance, &u.IsBlocked, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByIDInTx(ctx context.Context, tx pgx.Tx, id int) (*model.User, error) {
	query := `SELECT id, email, password_hash, role, balance, is_blocked, created_at FROM users WHERE id = $1 FOR UPDATE`
	var u model.User
	err := tx.QueryRow(ctx, query, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.Balance, &u.IsBlocked, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) UpdateBalanceInTx(ctx context.Context, tx pgx.Tx, userID int, delta float64) error {
	query := `UPDATE users SET balance = balance + $1 WHERE id = $2`
	_, err := tx.Exec(ctx, query, delta, userID)
	return err
}

func (r *UserRepo) UpdateBlockStatus(ctx context.Context, userID int, blocked bool) error {
	query := `UPDATE users SET is_blocked = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, blocked, userID)
	return err
}

func (r *UserRepo) List(ctx context.Context) ([]model.User, error) {
	query := `SELECT id, email, password_hash, role, balance, is_blocked, created_at FROM users ORDER BY id`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.Balance, &u.IsBlocked, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
