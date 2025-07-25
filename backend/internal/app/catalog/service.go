package catalog

import (
	"black-shop-service/internal/domain/catalog"
	"context"
	"go.uber.org/zap"
)

type Service struct {
	products catalog.ProductRepository
	logger   *zap.Logger
}

func NewService(repo catalog.ProductRepository, logger *zap.Logger) *Service {
	return &Service{
		products: repo,
		logger:   logger.Named("catalog_service"),
	}
}

func (s *Service) GetProduct(ctx context.Context, id catalog.ProductID) (*catalog.Product, error) {
	s.logger.Info("getting product by id", zap.String("product_id", string(id)))

	product, err := s.products.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get product by id from repository", zap.String("product_id", string(id)), zap.Error(err))
		return product, err
	}

	s.logger.Info("successfully found product", zap.String("product_id", string(id)))
	return s.products.FindByID(ctx, id)
}

func (s *Service) FindAllProducts(ctx context.Context) ([]*catalog.Product, error) {
	s.logger.Info("getting all products")

	products, err := s.products.FindAll(ctx)
	if err != nil {
		s.logger.Error("failed to get all products from repository", zap.Error(err))
		return nil, err
	}

	s.logger.Info("successfully found all products", zap.Int("count", len(products)))
	return products, nil
}

func (s *Service) CreateProduct(ctx context.Context, name string) (*catalog.Product, error) {
	s.logger.Info("creating a new product", zap.String("name", name))

	newProduct := &catalog.Product{
		Name: name,
	}

	if err := s.products.Save(ctx, newProduct); err != nil {
		s.logger.Error("failed to save product via repository", zap.Error(err))
		return nil, err
	}

	s.logger.Info("successfully created a new product", zap.String("product_id", string(newProduct.ID)))
	return newProduct, nil
}
