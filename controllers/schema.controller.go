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
	sc.SchemaService.CreateSchema(ctx, schema, userID)

}
func (sc *SchemaController) RegisterSchemaRoutes(rg *gin.RouterGroup) {
	_schemaroute := rg.Group("/schema")
	_schemaroute.POST("/create", helpers.AuthMiddleware(), sc.CreateSchema)
	// userroute.GET("/get/:name", uc.GetUser)
	// userroute.GET("/getall", uc.GetAll)
	// userroute.PATCH("/update", uc.UpdateUser)
	// userroute.DELETE("/delete/:name", uc.DeleteUser)
}
