package services

import (
	"gojinmongo/models"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(*gin.Context, *models.User) error
	GetUser(*string) (*models.User, error)
	GetAll() ([]*models.User, error)
	UpdateUser(*models.User) error
	DeleteUser(*string) error
}

type SchemaService interface {
	// product methods
}
