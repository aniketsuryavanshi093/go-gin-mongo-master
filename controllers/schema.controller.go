package controllers

import (
	"gojinmongo/services"

	"github.com/gin-gonic/gin"
)

type SchemaController struct {
	SchemaService services.SchemaService
}

func CreateSchemaController(schemaService services.SchemaService) SchemaController {
	return SchemaController{
		SchemaService: schemaService,
	}
}

func (sc *SchemaController) RegisterSchemaRoutes(rg *gin.RouterGroup) {
	// _schemaroute := rg.Group("/schema")
	// userroute.POST("/create", uc.CreateUser)
	// userroute.GET("/get/:name", uc.GetUser)
	// userroute.GET("/getall", uc.GetAll)
	// userroute.PATCH("/update", uc.UpdateUser)
	// userroute.DELETE("/delete/:name", uc.DeleteUser)
}
