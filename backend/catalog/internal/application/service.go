package application

import (
	"catalog/internal/domain"
	"context"
)

type ProductService interface {
	GetProduct(ctx context.Context, id domain.ProductID) (*domain.Product, error)
	FindAllProducts(ctx context.Context) ([]*domain.Product, error)
	CreateProduct(ctx context.Context, name string) (*domain.Product, error)
}

type CategoryService interface {
	CreateCategory(ctx context.Context, name string, imageUrl, parentID *string) (*domain.Category, error)
	GetAllCategories(ctx context.Context) ([]*domain.Category, error)
}
