package model

import "go.mongodb.org/mongo-driver/v2/bson"

type MongoProduct struct {
	ID   bson.ObjectID `bson:"_id,omitempty"`
	Name string        `bson:"name"`
}
