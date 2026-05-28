package testingutils

import (
	"context"
	"ssubench/internal/model"
	"ssubench/internal/repo"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type MockTaskRepo struct {
	tasks map[int]*model.Task
}

func NewMockTaskRepo() *MockTaskRepo {
	return &MockTaskRepo{tasks: make(map[int]*model.Task)}
}

func (m *MockTaskRepo) Create(ctx context.Context, t *model.Task) error {
	t.ID = len(m.tasks) + 1
	t.CreatedAt = time.Now()
	t.Status = model.TaskStatusOpen
	m.tasks[t.ID] = t
	return nil
}

func (m *MockTaskRepo) GetByID(ctx context.Context, id int) (*model.Task, error) {
	t, ok := m.tasks[id]
	if !ok {
		return nil, repo.ErrTaskNotFound
	}
	return t, nil
}

func (m *MockTaskRepo) UpdateStatus(ctx context.Context, id int, status model.TaskStatus) error {
	t, ok := m.tasks[id]
	if !ok {
		return repo.ErrTaskNotFound
	}
	t.Status = status
	return nil
}

func (m *MockTaskRepo) UpdateStatusInTx(ctx context.Context, tx pgx.Tx, id int, status model.TaskStatus) error {
	return m.UpdateStatus(ctx, id, status)
}

func (m *MockTaskRepo) List(ctx context.Context, filter model.TaskFilter) ([]model.Task, error) {
	var result []model.Task
	for _, t := range m.tasks {
		result = append(result, *t)
	}
	return result, nil
}

type MockBidRepo struct {
	bids map[int]*model.Bid
}

func NewMockBidRepo() *MockBidRepo {
	return &MockBidRepo{bids: make(map[int]*model.Bid)}
}

func (m *MockBidRepo) Create(ctx context.Context, b *model.Bid) error {
	b.ID = len(m.bids) + 1
	b.CreatedAt = time.Now()
	b.Status = model.BidStatusPending
	m.bids[b.ID] = b
	return nil
}

func (m *MockBidRepo) GetByID(ctx context.Context, id int) (*model.Bid, error) {
	b, ok := m.bids[id]
	if !ok {
		return nil, repo.ErrBidNotFound
	}
	return b, nil
}

func (m *MockBidRepo) UpdateStatus(ctx context.Context, id int, status model.BidStatus) error {
	b, ok := m.bids[id]
	if !ok {
		return repo.ErrBidNotFound
	}
	b.Status = status
	return nil
}

func (m *MockBidRepo) SetAllBidsStatus(ctx context.Context, taskID int, status model.BidStatus) error {
	for _, b := range m.bids {
		if b.TaskID == taskID {
			b.Status = status
		}
	}
	return nil
}

func (m *MockBidRepo) GetAcceptedBidByTaskID(ctx context.Context, taskID int) (*model.Bid, error) {
	for _, b := range m.bids {
		if b.TaskID == taskID && b.Status == model.BidStatusAccepted {
			return b, nil
		}
	}
	return nil, nil
}

type MockUserRepo struct {
	users map[int]*model.User
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{users: make(map[int]*model.User)}
}

func (m *MockUserRepo) Create(ctx context.Context, u *model.User) error {
	u.ID = len(m.users) + 1
	u.CreatedAt = time.Now()
	m.users[u.ID] = u
	return nil
}

func (m *MockUserRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, repo.ErrUserNotFound
	}
	return u, nil
}

func (m *MockUserRepo) GetByIDInTx(ctx context.Context, tx pgx.Tx, id int) (*model.User, error) {
	return m.GetByID(ctx, id)
}

func (m *MockUserRepo) UpdateBalanceInTx(ctx context.Context, tx pgx.Tx, userID int, delta float64) error {
	u, ok := m.users[userID]
	if !ok {
		return repo.ErrUserNotFound
	}
	u.Balance += delta
	return nil
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, repo.ErrUserNotFound
}

func (m *MockUserRepo) UpdateBlockStatus(ctx context.Context, userID int, blocked bool) error {
	u, ok := m.users[userID]
	if !ok {
		return repo.ErrUserNotFound
	}
	u.IsBlocked = blocked
	return nil
}

func (m *MockUserRepo) List(ctx context.Context) ([]model.User, error) {
	var result []model.User
	for _, u := range m.users {
		result = append(result, *u)
	}
	return result, nil
}

type MockPaymentRepo struct {
	payments []*model.Payment
}

func NewMockPaymentRepo() *MockPaymentRepo {
	return &MockPaymentRepo{}
}

func (m *MockPaymentRepo) CreateInTx(ctx context.Context, tx pgx.Tx, p *model.Payment) error {
	p.ID = len(m.payments) + 1
	p.CreatedAt = time.Now()
	m.payments = append(m.payments, p)
	return nil
}

type MockTx struct {
	committed  bool
	rolledback bool
}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	return m, nil
}

func (m *MockTx) Commit(ctx context.Context) error {
	m.committed = true
	return nil
}

func (m *MockTx) Rollback(ctx context.Context) error {
	m.rolledback = true
	return nil
}

func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}

func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return nil
}

func (m *MockTx) LargeObjects() pgx.LargeObjects {
	return pgx.LargeObjects{}
}

func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}

func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (m *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return nil
}

func (m *MockTx) Conn() *pgx.Conn {
	return nil
}

type MockPool struct {
	tx *MockTx
}

func NewMockPool() *MockPool {
	return &MockPool{tx: &MockTx{}}
}

func (m *MockPool) Begin(ctx context.Context) (pgx.Tx, error) {
	return m.tx, nil
}
