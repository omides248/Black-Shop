package application

import (
	"catalog/internal/domain"
	"context"
	"errors"
	"go.uber.org/zap"
)

type categoryService struct {
	categoryRepo domain.CategoryRepository
	productRepo  domain.ProductRepository
	logger       *zap.Logger
}

func NewCategoryService(catRepo domain.CategoryRepository, prodRepo domain.ProductRepository, logger *zap.Logger) CategoryService {
	return &categoryService{
		categoryRepo: catRepo,
		productRepo:  prodRepo,
		logger:       logger.Named("category_service"),
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, name string, imageUrl, parentID *string) (*domain.Category, error) {
	s.logger.Info("creating a new category", zap.String("name", name))

	// Rule 1: Limit Depth
	// Rule 2: Do not add subcategories to a category that has products
	parentCategory, err := s.validateParentCategory(ctx, parentID)
	if err != nil {
		return nil, err
	}

	// Rule 3: Redundant duplicate category name for sub category
	if err := s.checkDuplicateCategoryName(ctx, name, parentCategory); err != nil {
		return nil, err
	}

	newCategory, err := domain.NewCategory(name, imageUrl, parentCategory)
	if err != nil {
		s.logger.Warn("failed to create new object category from domain factory", zap.String("name", name))
		return nil, err
	}

	if err := s.categoryRepo.Save(ctx, newCategory); err != nil {
		s.logger.Error("failed to create category via repository", zap.Error(err))
		return nil, err
	}

	s.logger.Info("category created successfully", zap.String("category_id", string(newCategory.ID)), zap.String("name", newCategory.Name))
	return newCategory, nil

}

func (s *categoryService) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	s.logger.Info("getting all categories")
	return s.categoryRepo.FindAll(ctx)
}

func (s *categoryService) validateParentCategory(ctx context.Context, parentID *string) (*domain.Category, error) {

	if parentID == nil {
		return nil, nil // Category is root
	}

	parentCatID := domain.CategoryID(*parentID)
	parentCategory, err := s.categoryRepo.FindByID(ctx, parentCatID)

	if err != nil {
		s.logger.Warn("parent category not found", zap.String("parent_id", *parentID))
		return nil, err
	}

	// Rule 1: Limit Depth
	if parentCategory.Depth >= 2 {
		s.logger.Warn("depth limit exceeded", zap.Error(domain.ErrCategoryDepthLimitExceeded), zap.Int("parent_depth", parentCategory.Depth))
		return nil, domain.ErrCategoryDepthLimitExceeded
	}

	// Rule 2: Do not add subcategories to a category that has products
	hasProducts, err := s.productRepo.CategoryHasProducts(ctx, parentCatID)
	if err != nil {
		s.logger.Error("failed to check for products in parent category", zap.Error(err))
		return nil, err
	}

	if hasProducts {
		s.logger.Warn("attempted to add sub-category to a category with products", zap.Error(domain.ErrCategoryHasProducts))
		return nil, domain.ErrCategoryHasProducts
	}

	return parentCategory, nil
}

func (s *categoryService) checkDuplicateCategoryName(ctx context.Context, name string, parentCategory *domain.Category) error {
	var parentCatID *domain.CategoryID
	if parentCategory != nil {
		parentCatID = &parentCategory.ID
	}

	_, err := s.categoryRepo.FindByNameAndParentID(ctx, name, parentCatID)
	if err == nil {
		return domain.ErrCategoryAlreadyExists
	}
	if !errors.Is(err, domain.ErrCategoryNotFound) {
		s.logger.Error("service: failed to check for duplicate category name", zap.Error(err))
		return err
	}

	return nil
}
