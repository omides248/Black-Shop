package catalog

import (
	"errors"
)

var ErrProductNotFound = errors.New("product not found")

type ProductID string
type Product struct {
	ID   ProductID
	Name string
}
