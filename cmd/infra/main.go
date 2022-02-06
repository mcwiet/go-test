package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
	"github.com/mcwiet/go-test/pkg/infra"
)

func main() {
	app := awscdk.NewApp(nil)

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	stackNamePrefix := "go-" + env

	// Auth
	authStackName := stackNamePrefix + "-auth"
	authStack := infra.NewAuthStack(app, authStackName, &infra.AuthStackProps{
		StackProps: awscdk.StackProps{
			StackName: &authStackName,
			Env:       newCdkEnvironment(),
		},
		EnvName: env,
	})

	// API
	apiStackName := stackNamePrefix + "-api"
	apiStack := infra.NewApiStack(app, apiStackName, &infra.ApiStackProps{
		StackProps: awscdk.StackProps{
			StackName: &apiStackName,
			Env:       newCdkEnvironment(),
		},
		EnvName: env,
	})

	// Define dependencies (from parameters)
	apiStack.AddDependency(authStack, jsii.String("AppSync API needs reference to Cognito User Pool"))

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func newCdkEnvironment() *awscdk.Environment {
	awsAccount := os.Getenv("AWS_ACCOUNT")
	awsRegion := os.Getenv("AWS_REGION")

	if awsAccount != "" && awsRegion != "" {
		return &awscdk.Environment{
			Account: &awsAccount,
			Region:  &awsRegion,
		}
	} else {
		return nil
	}
}
