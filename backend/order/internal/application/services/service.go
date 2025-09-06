package services

import (
	"context"
	"order/internal/domain"
)

type OrderService interface {
	CreateOrderFromCart(ctx context.Context, userID string, paymentMethod string) (*domain.Order, error)
	AddItemToCart(ctx context.Context, userID string, item domain.CartItem) (*domain.Cart, error)
	GetCart(ctx context.Context, userID string) (*domain.Cart, error)
}
