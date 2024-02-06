package services

import (
	"context"
	"fmt"
	"gojinmongo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SchemaServiceImpl struct {
	schemacollection *mongo.Collection
	ctx              context.Context
	userCollection   *mongo.Collection
}

func NewSchemaService(schemacollection *mongo.Collection, userCollection *mongo.Collection, ctx context.Context) SchemaService {
	return &SchemaServiceImpl{
		schemacollection: schemacollection,
		ctx:              ctx,
		userCollection:   userCollection,
	}
}

func (s *SchemaServiceImpl) CreateSchema(ctx *gin.Context, schema *models.Schema, userID string) {
	userid, _ := primitive.ObjectIDFromHex(userID)
	schema.User = userid
	fmt.Println(userID)
	fmt.Println(userid)
	fmt.Println(schema)
	res, err := s.schemacollection.InsertOne(s.ctx, schema)
	if err != nil {
		appErr := &AppError{400, err.Error()}
		ctx.Error(appErr)
		return
	}
	schemaID := res.InsertedID.(primitive.ObjectID)
	// Update user's schema array
	filter := bson.M{"_id": userid}
	update := bson.M{"$push": bson.M{"schemas": schemaID}}
	upres, err := s.userCollection.UpdateOne(s.ctx, filter, update)
	if err != nil {
		appErr := &AppError{400, err.Error()}
		ctx.Error(appErr)
		return
	}
	fmt.Println(upres, filter)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "user schema updated",
	})
}
