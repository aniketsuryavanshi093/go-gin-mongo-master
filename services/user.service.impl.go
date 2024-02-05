package services

import (
	"context"
	"errors"
	"gojinmongo/helpers"
	"gojinmongo/models"
	"net/http"

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
		appErr := &AppError{400, "User already exists"}
		ctx.Error(appErr)
		return appErr
	}
	hashedPassword, err := helpers.HashPassword(user.Password)
	user.Password = hashedPassword
	res, err := u.usercollection.InsertOne(u.ctx, user)
	tokenString := helpers.GenerateToken(user)
	user.Password = ""
	user.ID = res.InsertedID.(primitive.ObjectID)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User created",
		"token":   tokenString,
		"data":    user,
	})
	return err
}
func (u *UserServiceImpl) LoginUser(ctx *gin.Context, user *models.User) (*models.UserResponse, error) {
	var tempuser *models.User
	err := u.usercollection.FindOne(u.ctx, bson.M{"email": user.Email}).Decode(&tempuser)
	if tempuser == nil {
		appErr := &AppError{400, "User does not exist"}
		ctx.Error(appErr)
		return nil, appErr
	}
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(tempuser.Password), []byte(user.Password))
	if err != nil {
		appErr := &AppError{400, "Incorrect password"}
		ctx.Error(appErr)
		return nil, appErr
	}
	tokenString := helpers.GenerateToken(user)
	user.Password = ""
	userResponse := &models.UserResponse{
		User:  user,
		Token: tokenString,
	}
	return userResponse, nil
}

func (u *UserServiceImpl) GetUser(name *string) (*models.User, error) {
	var user *models.User
	query := bson.D{bson.E{Key: "name", Value: name}}
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
