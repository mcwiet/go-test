// main.go
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func hello(ctx context.Context, request interface{}) (string, error) {
	return fmt.Sprintf("%+v\n", request), nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(hello)
}
