package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CategoryRepository interface {
	Save(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	FindByID(ctx context.Context, id CategoryID) (*Category, error)
	FindAll(ctx context.Context) ([]*Category, error)
	HasChildren(ctx context.Context, id CategoryID) (bool, error)
	FindByNameAndParentID(ctx context.Context, name string, parentID *CategoryID) (*Category, error)
}

type ProductRepository interface {
	Save(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id ProductID) (*Product, error)
	FindAll(ctx context.Context, filterQuery bson.M, sortOptions bson.D, page, limit int) ([]*Product, int64, error)
	CategoryHasProducts(ctx context.Context, id CategoryID) (bool, error)
	// TODO: Add method for searching and filtering products
}

type ProductVariantRepository interface {
	Create(ctx context.Context, productVariant *ProductVariant) error
	Update(ctx context.Context, productVariant *ProductVariant) error
	FindByProductID(ctx context.Context, id ProductID) ([]*ProductVariant, error)
	FindByID(ctx context.Context, id productVariantID) (*ProductVariant, error)
}

type BrandRepository interface {
	Create(ctx context.Context, brand *Brand) error
	Update(ctx context.Context, brand *Brand) error
	FindByID(ctx context.Context, id BrandID) (*Brand, error)
	FindAll(ctx context.Context) ([]*Brand, error)
}

type TagRepository interface {
	Create(ctx context.Context, tag *Tag) error
	Update(ctx context.Context, tag *Tag) error
	FindByID(ctx context.Context, id TagID) (*Tag, error)
	FindAll(ctx context.Context) ([]*Tag, error)
}

type ReviewRepository interface {
	Create(ctx context.Context, review *Review) error
	Update(ctx context.Context, review *Review) error
	FindByProductID(ctx context.Context, id ProductID) (*Review, error)
}
