package infrastructure

import (
	"context"
	"database/sql"

	"order/internal/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrDuplicateRequest = errors.New("duplicated request")
)

type order struct {
	ID     uuid.UUID `json:"id"`
	Price  float64   `json:"price"`
	UserID int       `json:"user_id"`
}

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(
	db *sql.DB,
) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(
	ctx context.Context,
	key uuid.UUID,
	id uuid.UUID,
	price float64,
	userID int,
) (*domain.Order, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, errors.Wrap(err, "begin tx")
	}

	var existingOrderID uuid.UUID
	err = tx.QueryRow(
		"SELECT order_id FROM idempotency_keys WHERE key = $1",
		key.String(),
	).
		Scan(&existingOrderID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if err == nil {
		var o order
		err = tx.QueryRow("SELECT id, user_id, price FROM orders WHERE id = $1", existingOrderID).
			Scan(&o.ID, &o.UserID, &o.Price)
		if err != nil {
			return nil, errors.Wrap(err, "select order")
		}

		_ = tx.Commit()

		return &domain.Order{
			ID:     o.ID,
			Price:  o.Price,
			UserID: o.UserID,
		}, ErrDuplicateRequest
	}

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders (id, user_id, price) VALUES ($1, $2, $3)", id, userID, price,
	)
	if err != nil {
		return nil, errors.Wrap(err, "insert order")
	}

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO idempotency_keys (key, order_id) VALUES ($1, $2)", key.String(), id,
	)
	if err != nil {
		return nil, errors.Wrap(err, "insert idempotency key")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit tx")
	}

	return &domain.Order{
		ID:     id,
		Price:  price,
		UserID: userID,
	}, nil
}

func (r *OrderRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var o order
	err := r.db.QueryRowContext(ctx, "SELECT id, user_id, price FROM orders WHERE id = $1", id).
		Scan(&o.ID, &o.UserID, &o.Price)
	if err != nil {
		return nil, errors.Wrap(err, "select order")
	}

	return &domain.Order{
		ID:     o.ID,
		Price:  o.Price,
		UserID: o.UserID,
	}, nil
}
