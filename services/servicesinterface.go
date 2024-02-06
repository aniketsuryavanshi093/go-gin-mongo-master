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
	UpdateUser(*models.User) error
	DeleteUser(*string) error
}

type SchemaService interface {
	// product methods
	CreateSchema(*gin.Context, *models.Schema, string)
}
