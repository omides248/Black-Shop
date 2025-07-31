package memory

import (
	"black-shop/internal/catalog/domain"
	"context"
)

type productRepo struct {
	Products map[domain.ProductID]*domain.Product
}

func NewProductRepository() domain.ProductRepository {

	products := map[domain.ProductID]*domain.Product{
		"product-1": {ID: "1", Name: "Gaming laptop"},
		"product-2": {ID: "1", Name: "Wireless mouse"},
	}

	return &productRepo{Products: products}
}

func (r *productRepo) FindByID(ctx context.Context, id domain.ProductID) (*domain.Product, error) {
	if product, ok := r.Products[id]; ok {
		return product, nil
	}
	return nil, domain.ErrProductNotFound
}

func (r *productRepo) Save(ctx context.Context, product *domain.Product) error {
	return nil
}

func (r *productRepo) FindAll(ctx context.Context) ([]*domain.Product, error) {
	return []*domain.Product{}, nil
}

func (r *productRepo) Create(ctx context.Context, product *domain.Product) error {
	//TODO implement me
	panic("implement me")
}

func (r *productRepo) Update(ctx context.Context, product *domain.Product) error {
	//TODO implement me
	panic("implement me")
}

func (r *productRepo) CategoryHasProducts(ctx context.Context, id domain.CategoryID) (bool, error) {
	//TODO implement me
	panic("implement me")
}
