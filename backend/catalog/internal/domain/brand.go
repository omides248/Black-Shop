package domain

import "time"

type Brand struct {
	ID        BrandID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
