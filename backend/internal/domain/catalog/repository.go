package catalog

import (
	"context"
)

type ProductRepository interface {
	FindByID(ctx context.Context, id ProductID) (*Product, error)
	Save(ctx context.Context, product *Product) error
	FindAll(ctx context.Context) ([]*Product, error)
}
