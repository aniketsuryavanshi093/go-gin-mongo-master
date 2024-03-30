package services

import (
	"gojinmongo/models"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(*gin.Context, *models.User) error
	LoginUser(*gin.Context, *models.User) (*models.UserResponse, error)
	GetUser(*string) (*models.User, error)
	GetAll() ([]*models.User, error)
	GetFolders(*gin.Context, *string) ([]models.Folder, error)
	CreateFolder(*gin.Context, *models.Folder, *string) error
	UpdateUser(*models.User) error
	DeleteUser(*string) error
}

type SchemaService interface {
	// product methods
	CreateSchema(*gin.Context, *models.Schema, string) (*models.SchemaResponse, error)
}
