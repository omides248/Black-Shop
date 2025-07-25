package postgresql

import (
	"black-shop-service/internal/domain/order"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type orderRepo struct {
	conn   *pgx.Conn
	logger *zap.Logger
}

func NewOrderRepository(conn *pgx.Conn, logger *zap.Logger) (order.OrderRepository, error) {

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

	return &orderRepo{
		conn:   conn,
		logger: logger.Named("postgres_order_repo"),
	}, nil
}

func (r *orderRepo) Save(ctx context.Context, o *order.Order) error {
	r.logger.Info("saving a new order", zap.String("user_id", o.UserID))

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	orderQuery := `INSERT INTO orders (user_id, total_price, status, created_at) VALUES ($1, $2, $3, $4) returning id`

	var orderID order.OrderID
	err = tx.QueryRow(ctx, orderQuery, o.UserID, o.TotalPrice, o.Status, o.CreatedAt).Scan(&orderID)
	if err != nil {
		r.logger.Error("failed to insert order", zap.Error(err))
		return err
	}
	o.ID = orderID

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

func (r *orderRepo) FindByID(ctx context.Context, id order.OrderID) (*order.Order, error) {
	return nil, nil
}

func (r *orderRepo) FindByUserID(ctx context.Context, userID string) ([]*order.Order, error) {
	return nil, nil
}
