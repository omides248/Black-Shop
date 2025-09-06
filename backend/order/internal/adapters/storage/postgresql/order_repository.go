package postgresql

import (
	"context"
	"fmt"
	"order/internal/domain"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type orderRepo struct {
	conn   *pgx.Conn
	logger *zap.Logger
}

func NewOrderRepository(conn *pgx.Conn, logger *zap.Logger) (domain.OrderRepository, error) {

	_, err := conn.Exec(context.Background(), `CREATE EXTENSION IF NOT EXISTS "pgcrypto";`)
	if err != nil {
		return nil, fmt.Errorf("failed to enable pgcrypto extension: %w", err)
	}

	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS orders (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id TEXT NOT NULL,
			total_price NUMERIC(10, 2) NOT NULL,
			status TEXT NOT NULL,
			payment_method TEXT NOT NULL,
			payment_address TEXT,
			transaction_id TEXT,
			derivation_index BIGSERIAL UNIQUE,
			created_at TIMESTAMPTZ NOT NULL
		);

		CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
			product_id TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			price NUMERIC(10, 2) NOT NULL
		);
	`)

	if err != nil {
		logger.Fatal("failed to create order table", zap.Error(err))
		return nil, err
	}

	_, err = conn.Exec(context.Background(), `CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);`)
	if err != nil {
		logger.Fatal("failed to create index on orders status", zap.Error(err))
		return nil, err
	}

	return &orderRepo{
		conn:   conn,
		logger: logger.Named("postgres_order_repo"),
	}, nil
}

func (r *orderRepo) Save(ctx context.Context, o *domain.Order) error {
	r.logger.Info("saving a new order", zap.String("user_id", o.UserID))

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	orderQuery := `
					INSERT INTO orders (user_id, total_price, status, payment_method, payment_address, created_at) 
					VALUES  ($1, $2, $3, $4, $5, $6) 
					returning id, derivation_index
					`

	var orderID domain.OrderID
	var derivationIndex int64
	err = tx.QueryRow(ctx, orderQuery, o.UserID, o.TotalPrice, o.Status, o.PaymentMethod, o.PaymentAddress, o.CreatedAt).Scan(&orderID)
	if err != nil {
		r.logger.Error("failed to insert order", zap.Error(err))
		return err
	}
	o.ID = orderID
	o.DerivationIndex = &derivationIndex

	for _, item := range o.Items {
		itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) returning id`
		_, err := tx.Exec(ctx, itemQuery, orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			r.logger.Error("failed to insert order item", zap.Error(err))
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *orderRepo) Update(ctx context.Context, order *domain.Order) error {
	r.logger.Info("updating order", zap.String("order_id", string(order.ID)))

	query := `
        UPDATE orders 
        SET status = $2, transaction_id = $3 
        WHERE id = $1
    `
	_, err := r.conn.Exec(ctx, query, order.ID, order.Status, order.TransactionID)
	if err != nil {
		r.logger.Error("failed to update order", zap.String("order_id", string(order.ID)), zap.Error(err))
		return err
	}

	return nil
}

func (r *orderRepo) FindByID(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	return nil, nil
}

func (r *orderRepo) FindByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	return nil, nil
}

func (r *orderRepo) FindAwaitingPayment(ctx context.Context, limit, offset int) ([]*domain.Order, error) {
	r.logger.Info("finding orders awaiting payment")
	query := `
			SELECT id, user_id, total_price, status, payment_method, payment_address, derivation_index, created_at 
			FROM orders WHERE status = $1
			ORDER BY created_at ASC 
			LIMIT $2
			OFFSET $3
			`

	rows, err := r.conn.Query(ctx, query, domain.StatusAwaitingPayment, limit, offset)
	if err != nil {
		r.logger.Error("failed to query orders awaiting payment", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
		err := rows.Scan(&o.ID, &o.UserID, &o.TotalPrice, &o.Status, &o.PaymentMethod, &o.PaymentAddress, &o.DerivationIndex, &o.CreatedAt)
		if err != nil {
			r.logger.Error("failed to scan order row", zap.Error(err))
			continue
		}
		orders = append(orders, &o)
	}

	return orders, nil

}
