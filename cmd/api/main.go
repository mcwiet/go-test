package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/mcwiet/go-test/pkg/controller"
	"github.com/mcwiet/go-test/pkg/data"
	"github.com/mcwiet/go-test/pkg/encoding"
	"github.com/mcwiet/go-test/pkg/service"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	petController controller.PetController
)

func init() {
	session := session.Must(session.NewSession())
	ddbClient := dynamodb.New(session)
	cursorEncoder := encoding.NewCursorEncoder()

	tableName := os.Getenv("DDB_TABLE_NAME")
	petDao := data.NewPetDao(ddbClient, tableName)
	petService := service.NewPetService(&petDao, &cursorEncoder)
	petController = controller.NewPetController(&petService)
}

func handle(ctx context.Context, rawRequest interface{}) (interface{}, error) {
	request := controller.NewRequest(rawRequest)
	var response controller.Response

	log.Println(request.Info.ParentTypeName)
	log.Println(request.Info.FieldName)

	switch request.Info.ParentTypeName {
	case "Query":
		switch request.Info.FieldName {
		case "pet":
			response = petController.HandleGet(request)
		case "pets":
			response = petController.HandleList(request)
		default:
			response = controller.Response{Error: errors.New("query not recognized")}
		}
	case "Mutation":
		switch request.Info.FieldName {
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
