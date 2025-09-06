package services

import (
	"catalog/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
)

type Service struct {
	ProductService  ProductService
	CategoryService CategoryService
}

func NewService(productRepo domain.ProductRepository, categoryRepo domain.CategoryRepository, logger *zap.Logger) *Service {
	productSvc := NewProductService(productRepo, logger)
	categorySvc := NewCategoryService(categoryRepo, productRepo, logger)

	return &Service{
		ProductService:  productSvc,
		CategoryService: categorySvc,
	}
}

type ProductService interface {
	GetProduct(ctx context.Context, id domain.ProductID) (*domain.Product, error)
	FindAllProducts(ctx context.Context, filterQuery bson.M, sortOptions bson.D, page, limit int) ([]*domain.Product, int64, error)
	CreateProduct(ctx context.Context, name string) (*domain.Product, error)
}

type CategoryService interface {
	CreateCategory(ctx context.Context, name string, image, slug, parentID *string) (*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) error
	GetAllCategories(ctx context.Context) ([]*domain.Category, error)
	FindByID(ctx context.Context, id domain.CategoryID) (*domain.Category, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Category, error)
}
