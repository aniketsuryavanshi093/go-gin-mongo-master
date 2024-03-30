package services

import (
	"context"
	"fmt"
	"gojinmongo/models"
	"time"

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

func (s *SchemaServiceImpl) CreateSchema(ctx *gin.Context, schema *models.Schema, userID string) (*models.SchemaResponse, error) {
	userid, _ := primitive.ObjectIDFromHex(userID)
	schema.User = userid
	fmt.Print("userid", userid)
	schema.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	fmt.Println(userID)
	fmt.Println(userid)
	fmt.Println(schema)
	res, err := s.schemacollection.InsertOne(s.ctx, schema)
	if err != nil {
		appErr := &AppError{400, err.Error()}
		ctx.Error(appErr)
		return nil, appErr
	}
	schemaID := res.InsertedID.(primitive.ObjectID)
	// Update user's schema array
	filter := bson.M{"_id": userid}
	update := bson.M{"$push": bson.M{"schemas": schemaID}}
	_, err = s.userCollection.UpdateOne(s.ctx, filter, update)
	schema.ID = schemaID
	if err != nil {
		appErr := &AppError{400, err.Error()}
		ctx.Error(appErr)
		return nil, appErr
	}
	userResponse := &models.SchemaResponse{
		Schema: schema,
	}
	return userResponse, nil
}
