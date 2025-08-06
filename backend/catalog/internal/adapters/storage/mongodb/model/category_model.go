package model

import (
	"catalog/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type MongoCategory struct {
	ID        bson.ObjectID      `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	ImageURL  *string            `bson:"imageUrl,omitempty"`
	ParentID  *domain.CategoryID `bson:"parentId,omitempty"`
	Depth     int                `bson:"depth"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
