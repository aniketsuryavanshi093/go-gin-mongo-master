package controllers

import (
	"gojinmongo/helpers"
	"gojinmongo/models"
	"gojinmongo/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserService
}

func New(userservice services.UserService) UserController {
	return UserController{
		UserService: userservice,
	}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "isError": true})
		return
	}
	err := uc.UserService.CreateUser(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error(), "isError": true})
		return
	}
}

func (uc *UserController) LoginUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "isError": true})
		return
	}
	userResponse, err := uc.UserService.LoginUser(ctx, &user)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error(), "isError": true})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User Login successfully", "data": userResponse})
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	var userId = ctx.MustGet("user_id").(string)
	user, err := uc.UserService.GetUser(&userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) GetAll(ctx *gin.Context) {
	users, err := uc.UserService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := uc.UserService.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	var username string = ctx.Param("name")
	err := uc.UserService.DeleteUser(&username)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) getFolders(ctx *gin.Context) {
	var userId = ctx.MustGet("user_id").(string)

	folders, err := uc.UserService.GetFolders(ctx, &userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success", "data": folders})
}

func (uc *UserController) createFolder(ctx *gin.Context) {
	var userId = ctx.MustGet("user_id").(string)
	var folder models.Folder
	if err := ctx.ShouldBindJSON(&folder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := uc.UserService.CreateFolder(ctx, &folder, &userId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": " folder created successfully!", "data": folder})
}
func (uc *UserController) DeleteFolder(ctx *gin.Context) {
	var userId = ctx.MustGet("user_id").(string)
	folderid := ctx.Param("id")

	err := uc.UserService.DeleteFolder(ctx, folderid, userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": " folder Deleted successfully!"})
}

func (uc *UserController) getFolderdetails(ctx *gin.Context) {
	var userId = ctx.MustGet("user_id").(string)
	folderid := ctx.Param("id")

	folder, err := uc.UserService.GetFolderdetails(ctx, folderid, userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": " folder Fetched successfully!", "folder": folder})
}

func (uc *UserController) GetUserDaigrams(ctx *gin.Context) {
	var userId = ctx.MustGet("user_id").(string)
	type filter struct {
		sortby string
	}
	queryfilter := &filter{
		sortby: ctx.Query("sortby"),
	}
	daigrams, err := uc.UserService.GetUserDaigrams(ctx, userId, queryfilter.sortby)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": " daigrams Fetched successfully!", "data": daigrams})
}

func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	userroute := rg.Group("/user")
	userroute.POST("/create", uc.CreateUser)
	userroute.POST("/login", uc.LoginUser)
	userroute.GET("/get", helpers.AuthMiddleware(), uc.GetUser)
	userroute.GET("/getdaigrams", helpers.AuthMiddleware(), uc.GetUserDaigrams)
	userroute.GET("/folders", helpers.AuthMiddleware(), uc.getFolders)
	userroute.GET("/folder/:id", helpers.AuthMiddleware(), uc.getFolderdetails)
	userroute.POST("/folder", helpers.AuthMiddleware(), uc.createFolder)
	userroute.GET("/getall", uc.GetAll)
	userroute.PATCH("/update", uc.UpdateUser)
	userroute.DELETE("/delete/:name", uc.DeleteUser)
	userroute.DELETE("/deletefolder/:id", helpers.AuthMiddleware(), uc.DeleteFolder)
}
