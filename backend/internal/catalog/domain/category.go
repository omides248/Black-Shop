package domain

import (
	"errors"
	"time"
)

type Category struct {
	ID        CategoryID
	Name      string
	ImageURL  *string
	ParentID  *CategoryID
	Depth     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCategory(name string, imageURL *string, parentCategory *Category) (*Category, error) {

	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	depth := 0
	var parentID *CategoryID
	if parentCategory != nil {
		depth = parentCategory.Depth + 1
		parentID = &parentCategory.ID
	}

	return &Category{
		Name:      name,
		ImageURL:  imageURL,
		ParentID:  parentID,
		Depth:     depth,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
