package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Name    string               `json:"name" bson:"name"`
	Schemas []primitive.ObjectID `json:"schemas" bson:"schemas"`
}
