package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Schema struct {
	Title           string             `json:"title" bson:"title"`
	User            primitive.ObjectID `json:"user" bson:"user,omitempty"`
	ID              primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	CreatedAt       primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt       primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
	Tablesdata      string             `json:"tablesdata" bson:"tablesdata"`
	Tablesrelations string             `json:"tablesrelations" bson:"tablesrelations"`
}

type SchemaResponse struct {
	Schema *Schema
}
