// main.go
package main

import (
	"context"
	"errors"

	"github.com/mcwiet/go-test/pkg/controller"

	"github.com/aws/aws-lambda-go/lambda"
)

func handle(ctx context.Context, rawRequest interface{}) (interface{}, error) {
	request := controller.NewRequest(rawRequest)
	var ret interface{} = nil
	var err error = nil

	switch request.Info.FieldName {
	case "person":
		ret = controller.GetPerson(request)
	case "people":
		ret = controller.GetPeople(request)
	default:
		err = errors.New("request not recognized")
	}

	return ret, err
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handle)
}
