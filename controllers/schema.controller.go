package controllers

import (
	"fmt"
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
	if folderId == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid Folder id", "isError": true})
		return
	}
	if schemaId == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid Schema id", "isError": true})
		return
	}
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
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Schema ID is required", "isError": true})
		return
	}
	err := sc.SchemaService.DeleteSchema(ctx, schemaid, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Schema deleted successfully!"})
}

func (sc *SchemaController) GetSchema(ctx *gin.Context) {
	schemaid := ctx.Param("id")
	if schemaid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Schema ID is required", "isError": true})
		return
	}
	schema, err := sc.SchemaService.GetSchema(ctx, schemaid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Schema deleted successfully!", "data": schema})
}

func (sc *SchemaController) UpdateSchema(ctx *gin.Context) {
	// Parse the JSON request body into the defined struct
	fmt.Print("update schema")
	type RequestBody struct {
		Token string `json:"token"`
		// Add other fields if present in the request body
		Schema models.Schema `json:"schema"`
	}
	var requestBody RequestBody
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Access the token field from the parsed request body
	token := requestBody.Token
	// Continue with the rest of the code...
	helpers.CheckAuth(ctx, token)
	schemaid := ctx.Param("id")
	if schemaid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Schema ID is required", "isError": true})
		return
	}
	err := sc.SchemaService.UpdateSchema(ctx, schemaid, requestBody.Schema)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Schema updated successfully!"})
}

func (sc *SchemaController) RegisterSchemaRoutes(rg *gin.RouterGroup) {
	_schemaroute := rg.Group("/schema")
	_schemaroute.POST("/create", helpers.AuthMiddleware(), sc.CreateSchema)
	_schemaroute.POST("/addfolder", helpers.AuthMiddleware(), sc.AddSchematoFolder)
	_schemaroute.GET("/get/:id", helpers.AuthMiddleware(), sc.GetSchema)
	// userroute.GET("/getall", uc.GetAll)
	_schemaroute.POST("/update/:id", sc.UpdateSchema)
	_schemaroute.DELETE("/delete/:id", helpers.AuthMiddleware(), sc.DeleteSchema)
}
