package domain

import "time"

// Image Value Object
type Image struct {
	URL       string
	AltText   string
	IsPrimary bool
	Order     int
}

type Tag struct {
	ID        TagID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Review struct {
	ID        ReviewID
	ProductID ProductID
	UserID    string
	Rating    int
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Brand struct {
	ID        BrandID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
