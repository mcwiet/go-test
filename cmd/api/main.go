package main

import (
	"context"
	"errors"

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

	personDao := data.NewPersonDao(ddbClient, "go-api-primary-table")
	personService := service.NewPersonService(personDao)
	personController = controller.NewPerosnController(personService)
}

func handle(ctx context.Context, rawRequest interface{}) (interface{}, error) {
	request := controller.NewRequest(rawRequest)
	var response controller.Response

	switch request.Info.FieldName {
	case "person":
		response = personController.GetPerson(request)
	case "people":
		response = personController.GetPeople(request)
	default:
		response = controller.Response{
			Error: errors.New("request not recognized"),
		}
	}

	return response.Data, response.Error
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handle)
}
