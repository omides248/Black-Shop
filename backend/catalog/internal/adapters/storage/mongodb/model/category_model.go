package model

import (
	"catalog/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoCategory struct {
	ID        bson.ObjectID      `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Slug      *string            `bson:"slug,omitempty"`
	Image     *string            `bson:"image,omitempty"`
	ParentID  *domain.CategoryID `bson:"parentId,omitempty"`
	Depth     int                `bson:"depth"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
