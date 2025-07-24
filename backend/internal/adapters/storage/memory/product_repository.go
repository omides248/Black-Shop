package memory

import (
	"black-shop-service/internal/domain/catalog"
	"context"
)

type productRepo struct {
	Products map[catalog.ProductID]*catalog.Product
}

func NewProductRepository() catalog.ProductRepository {

	products := map[catalog.ProductID]*catalog.Product{
		"product-1": {ID: "1", Name: "Gaming laptop"},
		"product-2": {ID: "1", Name: "Wireless mouse"},
	}

	return &productRepo{Products: products}
}

func (r *productRepo) FindByID(ctx context.Context, id catalog.ProductID) (*catalog.Product, error) {
	if product, ok := r.Products[id]; ok {
		return product, nil
	}
	return nil, catalog.ErrProductNotFound
}

func (r *productRepo) Save(ctx context.Context, product *catalog.Product) error {
	return nil
}

func (r *productRepo) FindAll(ctx context.Context) ([]*catalog.Product, error) {
	return []*catalog.Product{}, nil
}
