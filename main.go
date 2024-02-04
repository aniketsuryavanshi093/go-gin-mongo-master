package main

import (
	"context"
	"fmt"
	"gojinmongo/controllers"
	"gojinmongo/services"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server           *gin.Engine
	userservice      services.UserService
	schemaservice    services.SchemaService
	usercontrollers  controllers.UserController
	schemacontroller controllers.SchemaController
	ctx              context.Context
	schemac          *mongo.Collection
	userc            *mongo.Collection
	mongoclient      *mongo.Client
	err              error
)

func init() {
	ctx = context.TODO()

	mongoconn := options.Client().ApplyURI("mongodb+srv://aniketsuryavanshi093:kMsaFYSHPe1MU1Bl@golangpractise.fy4qwfr.mongodb.net/")
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal("error while connecting with mongo", err)
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo", err)
	}

	fmt.Println("mongo connection established")

	// user initiation
	userc = mongoclient.Database("userdb").Collection("users")
	userservice = services.NewUserService(userc, ctx)
	usercontrollers = controllers.New(userservice)

	// user schemas initiation
	schemac = mongoclient.Database("userdb").Collection("schemas")
	schemaservice = services.NewSchemaService(schemac, ctx)
	schemacontroller = controllers.CreateSchemaController(schemaservice)

	// assigning servers to global variable
	server = gin.Default()
}

func main() {
	defer mongoclient.Disconnect(ctx)
	basepath := server.Group("/v1")
	// routes for user
	usercontrollers.RegisterUserRoutes(basepath)
	// routes for user schemas
	schemacontroller.RegisterSchemaRoutes(basepath)
	log.Fatal(server.Run(":9090"))
}
