package domain

import "errors"

var (
	ErrOrderNotFound  = errors.New("order not found")
	ErrCartNotFound   = errors.New("cart not found")
	ErrUserIDNotEmpty = errors.New("user ID cannot be empty")
)
