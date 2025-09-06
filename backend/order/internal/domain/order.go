package domain

import (
	"time"
)

type OrderID string
type OrderStatus string

const (
	StatusPending         OrderStatus = "PENDING"
	StatusAwaitingPayment OrderStatus = "AWAITING_PAYMENT"
	StatusPaid            OrderStatus = "PAID"
	StatusShipped         OrderStatus = "SHIPPED"
	StatusCancelled       OrderStatus = "CANCELLED"
)

type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}

type Order struct {
	ID              OrderID
	UserID          string
	Items           []OrderItem
	TotalPrice      float64
	PaymentMethod   string
	PaymentAddress  *string // Crypto address
	TransactionID   *string // Blockchain hash transaction
	DerivationIndex *int64
	Status          OrderStatus
	CreatedAt       time.Time
}

func NewOrder(userID, paymentMethod string, items []OrderItem) (*Order, error) {
	var totalPrice float64
	for _, item := range items {
		totalPrice += item.Price * float64(item.Quantity)
	}

	if userID == "" {
		return nil, ErrUserIDNotEmpty
	}

	return &Order{
		UserID:        userID,
		Items:         items,
		TotalPrice:    totalPrice,
		Status:        StatusPending,
		PaymentMethod: paymentMethod,
		CreatedAt:     time.Now(),
	}, nil
}
