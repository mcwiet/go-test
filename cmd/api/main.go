package main

import (
	"context"
	"errors"

	"github.com/mcwiet/go-test/pkg/controller"

	"github.com/aws/aws-lambda-go/lambda"
)

func handle(ctx context.Context, rawRequest interface{}) (interface{}, error) {
	request := controller.NewRequest(rawRequest)
	var response controller.Response

	switch request.Info.FieldName {
	case "person":
		response = controller.GetPerson(request)
	case "people":
		response = controller.GetPeople(request)
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
