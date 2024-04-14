package services

import (
	"context"
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
	schema.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	res, err := s.schemacollection.InsertOne(s.ctx, schema)
	if err != nil {
		appErr := &AppError{400, err.Error(), true}
		ctx.Error(appErr)
		return nil, appErr
	}
	schemaID := res.InsertedID.(primitive.ObjectID)
	// Update user's schema array
	filter := bson.M{"_id": userid}
	update := bson.M{"$addToSet": bson.M{"schemas": schemaID}}
	_, err = s.userCollection.UpdateOne(s.ctx, filter, update)
	schema.ID = schemaID
	if err != nil {
		appErr := &AppError{400, err.Error(), true}
		ctx.Error(appErr)
		return nil, appErr
	}
	userResponse := &models.SchemaResponse{
		Schema: schema,
	}
	return userResponse, nil
}

func (s *SchemaServiceImpl) AddSchematoFolder(context *gin.Context, folderId string, schemaId string, userid string) error {
	userId, _ := primitive.ObjectIDFromHex(userid)
	folderid, _ := primitive.ObjectIDFromHex(folderId)
	schemaid, _ := primitive.ObjectIDFromHex(schemaId)
	filter := bson.M{"_id": userId, "folders._id": folderid}
	update := bson.M{"$addToSet": bson.M{"folders.$.schemaIds": schemaid}}
	_, err := s.userCollection.UpdateOne(s.ctx, filter, update)
	if err != nil {
		appErr := &AppError{400, err.Error(), true}
		context.Error(appErr)
		return appErr
	}
	return nil
}

func (s *SchemaServiceImpl) DeleteSchema(context *gin.Context, schemaId string, userId string) error {
	schemaid, _ := primitive.ObjectIDFromHex(schemaId)
	userid, _ := primitive.ObjectIDFromHex(userId)
	_, err := s.schemacollection.DeleteOne(s.ctx, bson.M{"_id": schemaid})
	_, updateerr := s.userCollection.UpdateByID(s.ctx, userid, bson.M{"$pull": bson.M{"schemas": schemaid}})
	if updateerr != nil {
		appErr := &AppError{400, updateerr.Error(), true}
		context.Error(appErr)
		return appErr
	}
	if err != nil {
		appErr := &AppError{400, err.Error(), true}
		context.Error(appErr)
		return appErr
	}
	return nil
}

func (s *SchemaServiceImpl) GetSchema(context *gin.Context, schemaId string) (*models.SchemaResponse, error) {
	schemaid, _ := primitive.ObjectIDFromHex(schemaId)
	var tempschema *models.Schema

	err := s.schemacollection.FindOne(s.ctx, bson.M{"_id": schemaid}).Decode(&tempschema)
	if err != nil {
		appErr := &AppError{400, err.Error(), true}
		context.Error(appErr)
		return nil, appErr
	}
	tempschemares := &models.SchemaResponse{
		Schema: tempschema,
	}
	return tempschemares, nil
}
