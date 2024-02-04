package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Name       string               `json:"name" bson:"name"`
	ProfilePic string               `json:"profilePic" bson:"profilePic"`
	Schemas    []primitive.ObjectID `json:"schemas" bson:"schemas"`
	Email      string               `json:"email" bson:"email,omitempty"`
	ID         primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	CreatedAt  string               `json:"createdAt" bson:"createdAt"`
	Password   string               `json:"password" bson:"password,omitempty"`
}
