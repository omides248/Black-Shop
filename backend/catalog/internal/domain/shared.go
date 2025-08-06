package domain

import "errors"

// TODO move to errors pkg
var (
	ErrProductNotFound            = errors.New("product not found")
	ErrCategoryNotFound           = errors.New("category not found")
	ErrBrandNotFound              = errors.New("brand not found")
	ErrReviewNotFound             = errors.New("review not found")
	ErrVariantNotFound            = errors.New("variant not found")
	ErrTagNotFound                = errors.New("tag not found")
	ErrCategoryAlreadyExists      = errors.New("a category with this name already exists at this level")
	ErrCategoryDepthLimitExceeded = errors.New("category depth limit exceeded")
	ErrCategoryHasProducts        = errors.New("cannot add sub-category to a category that already contains products")
)

type ProductID string
type CategoryID string
type BrandID string
type ReviewID string
type productVariantID string
type TagID string
