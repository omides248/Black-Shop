package domain

import "errors"

var ErrCartNotFound = errors.New("cart not found")

type CartItem struct {
	ProductID string
	Quantity  int
}

type Cart struct {
	UserID string
	Items  []CartItem
}
