package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Schema struct {
	Title     string             `json:"title" bson:"title"`
	User      primitive.ObjectID `json:"user" bson:"user,omitempty"`
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	CreatedAt string             `json:"createdAt" bson:"createdAt"`
}

type SchemaResponse struct {
	Schema *Schema
}
