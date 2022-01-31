package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/mcwiet/go-test/pkg/controller"
	"github.com/mcwiet/go-test/pkg/data"
	"github.com/mcwiet/go-test/pkg/service"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	personController controller.PersonController
)

func init() {
	sess := session.Must(session.NewSession())
	ddbClient := dynamodb.New(sess)

	tableName := os.Getenv("DDB_TABLE_NAME")
	personDao := data.NewPersonDao(ddbClient, tableName)
	personService := service.NewPersonService(personDao)
	personController = controller.NewPerosnController(personService)
}

func handle(ctx context.Context, rawRequest interface{}) (interface{}, error) {
	request := controller.NewRequest(rawRequest)
	var response controller.Response

	log.Println(request.Info.ParentTypeName)
	log.Println(request.Info.FieldName)

	switch request.Info.ParentTypeName {
	case "Query":
		switch request.Info.FieldName {
		case "person":
			response = personController.HandleGetPerson(request)
		case "people":
			response = personController.HandleGetPeople(request)
		default:
			response = controller.Response{Error: errors.New("query not recognized")}
		}
	case "Mutation":
		switch request.Info.FieldName {
		case "createPerson":
			response = personController.HandleCreatePerson(request)
		case "deletePerson":
			response = personController.HandleDeletePerson(request)
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
