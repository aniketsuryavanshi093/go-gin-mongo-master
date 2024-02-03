package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Schema struct {
	Title     string             `json:"title" bson:"title"`
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	CreatedAt string             `json:"createdAt" bson:"createdAt"`
}
