package services

import (
	"gojinmongo/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type UserService interface {
	CreateUser(*gin.Context, *models.User) error
	LoginUser(*gin.Context, *models.User) (*models.UserResponse, error)
	GetUser(*string) (*models.User, error)
	GetAll() ([]*models.User, error)
	GetUserDaigrams(*gin.Context, string, string) ([]bson.M, error)
	GetFolderdetails(*gin.Context, string, string) ([]bson.M, error)
	GetFolders(*gin.Context, *string) ([]models.Folder, error)
	CreateFolder(*gin.Context, *models.Folder, *string) error
	DeleteFolder(*gin.Context, string, string) error
	UpdateUser(*models.User) error
	DeleteUser(*string) error
}

type SchemaService interface {
	// product methods
	GetSchema(*gin.Context, string) (*models.SchemaResponse, error)
	CreateSchema(*gin.Context, *models.Schema, string) (*models.SchemaResponse, error)
	AddSchematoFolder(*gin.Context, string, string, string) error
	DeleteSchema(*gin.Context, string, string) error
	UpdateSchema(*gin.Context, string, models.Schema) error
}
