package application

import (
	"catalog/internal/domain"
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
)

type productService struct {
	productRepo domain.ProductRepository
	logger      *zap.Logger
}

func NewProductService(productRepo domain.ProductRepository, logger *zap.Logger) ProductService {
	return &productService{productRepo: productRepo, logger: logger.Named("product_service")}
}

func (s *productService) GetProduct(ctx context.Context, id domain.ProductID) (*domain.Product, error) {
	s.logger.Info("getting product by id", zap.String("product_id", string(id)))

	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get product by id from repository", zap.String("product_id", string(id)), zap.Error(err))
		return product, err
	}

	s.logger.Info("successfully found product", zap.String("product_id", string(id)))
	return s.productRepo.FindByID(ctx, id)
}

func (s *productService) FindAllProducts(ctx context.Context, filterQuery bson.M, sortOptions bson.D, page, limit int) ([]*domain.Product, int64, error) {
	s.logger.Info("getting all products")

	products, total, err := s.productRepo.FindAll(ctx, filterQuery, sortOptions, page, limit)
	if err != nil {
		s.logger.Error("failed to get all products from repository", zap.Error(err))
		return nil, 0, err
	}

	s.logger.Info("successfully found all products", zap.Int("count", len(products)))
	return products, total, nil
}

func (s *productService) CreateProduct(ctx context.Context, name string) (*domain.Product, error) {
	s.logger.Info("creating a new product", zap.String("name", name))

	newProduct := &domain.Product{
		Name: name,
	}

	if err := s.productRepo.Save(ctx, newProduct); err != nil {
		s.logger.Error("failed to save product via repository", zap.Error(err))
		return nil, err
	}

	s.logger.Info("successfully created a new product", zap.String("product_id", string(newProduct.ID)))
	return newProduct, nil
}
