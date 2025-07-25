package order

import (
	"errors"
	"time"
)

var ErrOrderNotFound = errors.New("order not found")

type OrderID string
type OrderStatus string

const (
	StatusPending    OrderStatus = "PENDING"
	StatusProcessing OrderStatus = "PAID"
	StatusShipped    OrderStatus = "SHIPPED"
	StatusCancelled  OrderStatus = "CANCELLED"
)

type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}

type Order struct {
	ID         OrderID
	UserID     string
	Items      []OrderItem
	TotalPrice float64
	Status     OrderStatus
	CreatedAt  time.Time
}

func NewOrder(userID string, items []OrderItem) (*Order, error) {
	var totalPrice float64
	for _, item := range items {
		totalPrice += item.Price * float64(item.Quantity)
	}

	return &Order{
		UserID:     userID,
		Items:      items,
		TotalPrice: totalPrice,
		Status:     StatusPending,
		CreatedAt:  time.Now(),
	}, nil
}
