package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/mcwiet/go-test/pkg/authorization"
	"github.com/mcwiet/go-test/pkg/controller"
	"github.com/mcwiet/go-test/pkg/data"
	"github.com/mcwiet/go-test/pkg/encoding"
	"github.com/mcwiet/go-test/pkg/service"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	petController  controller.PetController
	userController controller.UserController
)

func init() {
	session := session.Must(session.NewSession())
	ddbClient := dynamodb.New(session)
	cognitoClient := cognitoidentityprovider.New(session)
	cursorEncoder := encoding.NewCursorEncoder()

	// Authorization
	petAuth := authorization.NewPetAuthorizer()

	// Data
	primaryTableName := os.Getenv("DDB_PRIMARY_TABLE_NAME")
	petDao := data.NewPetDao(ddbClient, primaryTableName)
	userPoolId := os.Getenv("USER_POOL_ID")
	userDao := data.NewUserDao(cognitoClient, userPoolId)

	// Service
	petService := service.NewPetService(&petDao, &userDao, &petAuth, &cursorEncoder)
	userService := service.NewUserService(&userDao, &cursorEncoder)

	// Controller
	petController = controller.NewPetController(&petService)
	userController = controller.NewUserController(&userService)
}

func handle(ctx context.Context, req interface{}) (interface{}, error) {
	request := NewRequest(req)
	log.Println(request.ParentTypeName + " " + request.FieldName)

	var response controller.Response
	switch request.ParentTypeName {
	case "Query":
		switch request.FieldName {
		case "pet":
			response = petController.HandleGet(request)
		case "pets":
			response = petController.HandleList(request)
		case "user":
			response = userController.HandleGet(request)
		case "users":
			response = userController.HandleList(request)
		default:
			response = controller.Response{Error: errors.New("query not recognized")}
		}
	case "Mutation":
		switch request.FieldName {
		case "createPet":
			response = petController.HandleCreate(request)
		case "deletePet":
			response = petController.HandleDelete(request)
		case "updatePetOwner":
			response = petController.HandleUpdateOwner(request)
		default:
			response = controller.Response{Error: errors.New("mutation not recognized")}
		}
	default:
		response = controller.Response{Error: errors.New("request type not recognized")}
	}

	return response.Data, response.Error
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handle)
}
