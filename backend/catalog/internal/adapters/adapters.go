package adapters

import (
	"catalog/internal/adapters/storage/mongodb"
	"catalog/internal/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type Adapter struct {
	ProductRepo  domain.ProductRepository
	CategoryRepo domain.CategoryRepository
}

func NewAdapter(db *mongo.Database, logger *zap.Logger) (*Adapter, error) {

	productRepo, err := mongodb.NewProductRepository(db, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create product repository: %w", err)
	}

	categoryRepo, err := mongodb.NewCategoryRepository(db, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create category repository: %w", err)
	}

	return &Adapter{
		ProductRepo:  productRepo,
		CategoryRepo: categoryRepo,
	}, nil

}
