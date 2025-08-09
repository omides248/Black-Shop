package domain

import "time"

// Product Entity
type Product struct {
	ID              ProductID
	Name            string
	Description     string
	PrimaryImageURL *string
	AverageRating   float64
	CategoryID      CategoryID
	BrandID         *BrandID
	Tags            []TagID
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ProductVariant Entity
type ProductVariant struct {
	ID         productVariantID
	ProductID  ProductID
	SKU        string
	Price      float64
	Stock      int
	Images     []Image
	Attributes []Attribute
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Image struct {
	URL       string
	AltText   string
	IsPrimary bool
	Order     int
}

type Attribute struct {
	Name  string
	Value string
}
