package order

import "context"

type OrderRepository interface {
	Save(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id OrderID) (*Order, error)
	FindByUserID(ctx context.Context, userID string) ([]*Order, error)
}

type CartRepository interface {
	Save(ctx context.Context, cart *Cart) error
	GetByUserID(ctx context.Context, userID string) (*Cart, error)
	Delete(ctx context.Context, userID string) error
}
