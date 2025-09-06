package domain

import "time"

type Review struct {
	ID        ReviewID
	ProductID ProductID
	UserID    string
	Rating    int
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
