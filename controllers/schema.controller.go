package controllers

import (
	"gojinmongo/helpers"
	"gojinmongo/models"
	"gojinmongo/services"
	"net/http"

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
func (sc *SchemaController) CreateSchema(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	var schema *models.Schema
	if err := ctx.ShouldBindJSON(&schema); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	schemaresponse, err := sc.SchemaService.CreateSchema(ctx, schema, userID)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User Schema created successfully!", "data": schemaresponse})
}
func (sc *SchemaController) RegisterSchemaRoutes(rg *gin.RouterGroup) {
	_schemaroute := rg.Group("/schema")
	_schemaroute.POST("/create", helpers.AuthMiddleware(), sc.CreateSchema)
	// userroute.GET("/get/:name", uc.GetUser)
	// userroute.GET("/getall", uc.GetAll)
	// userroute.PATCH("/update", uc.UpdateUser)
	// userroute.DELETE("/delete/:name", uc.DeleteUser)
}
