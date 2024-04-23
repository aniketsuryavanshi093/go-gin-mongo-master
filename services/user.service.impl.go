package services

import (
	"context"
	"errors"
	"gojinmongo/helpers"
	"gojinmongo/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	usercollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(usercollection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{
		usercollection: usercollection,
		ctx:            ctx,
	}
}
func (u *UserServiceImpl) CreateUser(ctx *gin.Context, user *models.User) error {
	var tempuser *models.User
	u.usercollection.FindOne(u.ctx, bson.M{"email": user.Email}).Decode(&tempuser)
	if tempuser != nil {
		appErr := &AppError{400, "User already exists", true}
		ctx.Error(appErr)
		return appErr
	}
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		// Handle error
		return err
	}
	user.Password = hashedPassword

	if user.Schemas == nil {
		user.Schemas = []primitive.ObjectID{}
	}

	if user.Folders == nil {
		user.Folders = []models.Folder{}
	}

	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	res, err := u.usercollection.InsertOne(u.ctx, user)
	if err != nil {
		// Handle error
		return err
	}

	tokenString := helpers.GenerateToken(user)
	user.Password = ""
	user.ID = res.InsertedID.(primitive.ObjectID)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User created",
		"token":   tokenString,
		"data":    user,
	})
	return nil
}

func (u *UserServiceImpl) LoginUser(ctx *gin.Context, user *models.User) (*models.UserResponse, error) {
	var tempuser *models.User
	query := bson.D{bson.E{Key: "email", Value: user.Email}}
	err := u.usercollection.FindOne(u.ctx, query).Decode(&tempuser)
	if tempuser == nil {
		appErr := &AppError{400, "User does not exist", true}
		ctx.Error(appErr)
		return nil, appErr
	}
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(tempuser.Password), []byte(user.Password))
	if err != nil {
		appErr := &AppError{400, "Incorrect password", true}
		ctx.Error(appErr)
		return nil, appErr
	}

	tempuser.Password = ""
	tokenString := helpers.GenerateToken(tempuser)
	userResponse := &models.UserResponse{
		User:  tempuser,
		Token: tokenString,
	}
	return userResponse, nil
}

func (u *UserServiceImpl) GetUser(name *string) (*models.User, error) {
	var user *models.User
	userObjectID, _ := primitive.ObjectIDFromHex(*name)
	query := bson.D{bson.E{Key: "_id", Value: userObjectID}}
	err := u.usercollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}

func (u *UserServiceImpl) GetAll() ([]*models.User, error) {
	var users []*models.User
	cursor, err := u.usercollection.Find(u.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(u.ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(u.ctx)

	if len(users) == 0 {
		return nil, errors.New("documents not found")
	}
	return users, nil
}

func (u *UserServiceImpl) UpdateUser(user *models.User) error {
	filter := bson.D{primitive.E{Key: "name", Value: user.Name}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "name", Value: user.Name}, primitive.E{Key: "age"}}}}
	result, _ := u.usercollection.UpdateOne(u.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (u *UserServiceImpl) DeleteUser(name *string) error {
	filter := bson.D{primitive.E{Key: "name", Value: name}}
	result, _ := u.usercollection.DeleteOne(u.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}

func (u *UserServiceImpl) GetFolders(ctx *gin.Context, userid *string) ([]models.Folder, error) {
	var folders []models.Folder
	var user *models.User
	userObjectID, err := primitive.ObjectIDFromHex(*userid)
	if err != nil {
		appErr := &AppError{400, "Invalid user ID", true}
		ctx.Error(appErr)
		return nil, appErr
	}
	u.usercollection.FindOne(u.ctx, bson.M{"_id": userObjectID}).Decode(&user)
	if user == nil {
		appErr := &AppError{400, "User does not exist", true}
		ctx.Error(appErr)
		return nil, appErr
	}
	for _, folder := range user.Folders {
		folders = append(folders, folder)
	}

	return folders, nil
}

func (u *UserServiceImpl) CreateFolder(ctx *gin.Context, folder *models.Folder, userid *string) error {
	folder.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	folder.ID = primitive.NewObjectID()
	folder.SchemaIds = []primitive.ObjectID{}
	userObjectID, err := primitive.ObjectIDFromHex(*userid)
	if err != nil {
		appErr := &AppError{400, "Invalid user ID", true}
		ctx.Error(appErr)
		return appErr
	}

	update := bson.M{"$push": bson.M{"folders": folder}}
	_, err = u.usercollection.UpdateByID(u.ctx, userObjectID, update)
	if err != nil {
		appErr := &AppError{400, err.Error(), true}
		ctx.Error(appErr)
		return appErr
	}

	return nil
}

func (u *UserServiceImpl) DeleteFolder(ctx *gin.Context, folderId string, userId string) error {
	userObjectID, _ := primitive.ObjectIDFromHex(userId)
	folderid, _ := primitive.ObjectIDFromHex(folderId)
	filter := bson.M{"_id": userObjectID, "folders._id": folderid}
	update := bson.M{"$pull": bson.M{"folders": bson.M{"_id": folderid}}}
	_, err := u.usercollection.UpdateOne(u.ctx, filter, update)
	if err != nil {
		appErr := &AppError{400, err.Error(), true}
		ctx.Error(appErr)
		return appErr
	}
	return nil
}

func (u *UserServiceImpl) GetFolderdetails(ctx *gin.Context, folderId string, userId string) ([]bson.M, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	folderObjectID, err := primitive.ObjectIDFromHex(folderId)
	if err != nil {
		return nil, err
	}

	// Match user and folder
	// Match user and folder
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "_id", Value: userObjectID},
			{Key: "folders._id", Value: folderObjectID},
		}},
	}

	// Unwind folders array
	unwindStage := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$folders"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}},
	}

	// Match the folder
	folderMatchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "folders._id", Value: folderObjectID},
		}},
	}

	// Lookup schemas
	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "schemas"},
			{Key: "localField", Value: "folders.schemaIds"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "schemas"},
		}},
	}

	// Reshape result
	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "folder", Value: "$folders"},
			{Key: "schemas", Value: 1},
		}},
	}

	pipeline := mongo.Pipeline{
		matchStage,
		unwindStage,
		folderMatchStage,
		lookupStage,
		projectStage,
	}

	cursor, err := u.usercollection.Aggregate(ctx, pipeline)

	var results []bson.M

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
func (u *UserServiceImpl) GetUserDaigrams(ctx *gin.Context, userId string, filter string) ([]bson.M, error) {
	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	// Match user
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: userID}}}}

	// Lookup schemas
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "schemas"},
		{Key: "localField", Value: "schemas"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "schemas"},
	}}}

	var unwindStage, sortStage, groupStage bson.D
	switch filter {
	case "Alphabetical":
		unwindStage = bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$schemas"}}}}
		sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: "schemas.title", Value: 1}}}}
		groupStage = bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$_id"}, {Key: "schemas", Value: bson.D{{Key: "$push", Value: "$schemas"}}}}}}
	case "LastUpdatedAt":
		unwindStage = bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$schemas"}}}}
		sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: "schemas.updatedAt", Value: -1}}}}
		groupStage = bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$_id"}, {Key: "schemas", Value: bson.D{{Key: "$push", Value: "$schemas"}}}}}}
	case "LastCreatedAt":
		unwindStage = bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$schemas"}}}}
		sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: "schemas.createdAt", Value: -1}}}}
		groupStage = bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$_id"}, {Key: "schemas", Value: bson.D{{Key: "$push", Value: "$schemas"}}}}}}
	default:
		unwindStage = bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$schemas"}}}}
		sortStage = bson.D{{Key: "$sort", Value: bson.D{{Key: "schemas.title", Value: 1}}}}
		groupStage = bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$_id"}, {Key: "schemas", Value: bson.D{{Key: "$push", Value: "$schemas"}}}}}}
	}

	// Reshape result
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "schemas", Value: "$schemas"},
	}}}

	cursor, err := u.usercollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		sortStage,
		groupStage,
		projectStage,
	})
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
