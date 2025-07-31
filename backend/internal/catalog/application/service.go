package application

import (
	"black-shop/internal/catalog/domain"
	"go.uber.org/zap"
)

type ServiceDependencies struct {
	CategoryRepo domain.CategoryRepository
	ProductRepo  domain.ProductRepository
	ReviewRepo   domain.ReviewRepository
	BrandRepo    domain.BrandRepository
	TagRepo      domain.TagRepository

	Logger *zap.Logger
}

type Service struct {
	categoryRepo domain.CategoryRepository
	productRepo  domain.ProductRepository
	reviewRepo   domain.ReviewRepository
	brandRepo    domain.BrandRepository
	tagRepo      domain.TagRepository

	logger *zap.Logger
}

func NewService(deps ServiceDependencies) *Service {
	return &Service{
		categoryRepo: deps.CategoryRepo,
		productRepo:  deps.ProductRepo,
		reviewRepo:   deps.ReviewRepo,
		brandRepo:    deps.BrandRepo,
		tagRepo:      deps.TagRepo,
		logger:       deps.Logger.Named("catalog_service"),
	}
}
