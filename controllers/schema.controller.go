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
func (sc *SchemaController) AddSchematoFolder(ctx *gin.Context) {

	userID := ctx.MustGet("user_id").(string)

	var requestBody struct {
		SchemaId string `json:"schemaId"`
		FolderId string `json:"folderId"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	schemaId := requestBody.SchemaId
	folderId := requestBody.FolderId

	err := sc.SchemaService.AddSchematoFolder(ctx, folderId, schemaId, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Schema added to folder successfully!"})
}

func (sc *SchemaController) DeleteSchema(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	schemaid := ctx.Param("id")
	if schemaid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Schema ID is required"})
		return
	}
	err := sc.SchemaService.DeleteSchema(ctx, schemaid, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Schema deleted successfully!"})
}

func (sc *SchemaController) RegisterSchemaRoutes(rg *gin.RouterGroup) {
	_schemaroute := rg.Group("/schema")
	_schemaroute.POST("/create", helpers.AuthMiddleware(), sc.CreateSchema)
	_schemaroute.POST("/addfolder", helpers.AuthMiddleware(), sc.AddSchematoFolder)
	// userroute.GET("/get/:name", uc.GetUser)
	// userroute.GET("/getall", uc.GetAll)
	// userroute.PATCH("/update", uc.UpdateUser)
	_schemaroute.DELETE("/delete/:id", helpers.AuthMiddleware(), sc.DeleteSchema)
}
