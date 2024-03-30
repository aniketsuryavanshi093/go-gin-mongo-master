package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Name          string               `json:"name" bson:"name"`
	IsGoogleLogin bool                 `json:"isGoogleLogin" bson:"isGoogleLogin" default:"false"`
	ProfilePic    string               `json:"profilePic" bson:"profilePic"`
	Schemas       []primitive.ObjectID `json:"schemas" bson:"schemas,omitempty" default:"[]"`
	Folders       []Folder             `json:"folders" bson:"folders,omitempty" default:"[]"`
	Email         string               `json:"email" bson:"email,omitempty"`
	ID            primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	CreatedAt     primitive.DateTime   `json:"createdAt" bson:"createdAt"`
	Password      string               `json:"password" bson:"password,omitempty"`
}

type Folder struct {
	ID        primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Name      string               `json:"name" bson:"name"`
	SchemaIds []primitive.ObjectID `json:"schemaIds" bson:"schemaIds"`
	CreatedAt primitive.DateTime   `json:"createdAt" bson:"createdAt"`
}

type UserResponse struct {
	User  *User
	Token string
}
