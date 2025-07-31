package application

import (
	"black-shop/internal/catalog/domain"
	"context"
	"go.uber.org/zap"
)

func (s *Service) GetProduct(ctx context.Context, id domain.ProductID) (*domain.Product, error) {
	s.logger.Info("getting product by id", zap.String("product_id", string(id)))

	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get product by id from repository", zap.String("product_id", string(id)), zap.Error(err))
		return product, err
	}

	s.logger.Info("successfully found product", zap.String("product_id", string(id)))
	return s.productRepo.FindByID(ctx, id)
}

func (s *Service) FindAllProducts(ctx context.Context) ([]*domain.Product, error) {
	s.logger.Info("getting all products")

	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		s.logger.Error("failed to get all products from repository", zap.Error(err))
		return nil, err
	}

	s.logger.Info("successfully found all products", zap.Int("count", len(products)))
	return products, nil
}

func (s *Service) CreateProduct(ctx context.Context, name string) (*domain.Product, error) {
	s.logger.Info("creating a new product", zap.String("name", name))

	newProduct := &domain.Product{
		Name: name,
	}

	if err := s.productRepo.Create(ctx, newProduct); err != nil {
		s.logger.Error("failed to save product via repository", zap.Error(err))
		return nil, err
	}

	s.logger.Info("successfully created a new product", zap.String("product_id", string(newProduct.ID)))
	return newProduct, nil
}
