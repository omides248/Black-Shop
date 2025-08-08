package application

import (
	"catalog/internal/domain"
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductService interface {
	GetProduct(ctx context.Context, id domain.ProductID) (*domain.Product, error)
	FindAllProducts(ctx context.Context, filterQuery bson.M, sortOptions bson.D, page, limit int) ([]*domain.Product, int64, error)
	CreateProduct(ctx context.Context, name string) (*domain.Product, error)
}

type CategoryService interface {
	CreateCategory(ctx context.Context, name string, imageUrl, parentID *string) (*domain.Category, error)
	GetAllCategories(ctx context.Context) ([]*domain.Category, error)
}
