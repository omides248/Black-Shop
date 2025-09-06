package domain

import "time"

type Tag struct {
	ID        TagID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
