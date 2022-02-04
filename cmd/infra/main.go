package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/mcwiet/go-test/pkg/infra"
)

func main() {
	app := awscdk.NewApp(nil)

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	apiStackName := "go-api-" + env
	infra.NewApiStack(app, apiStackName, &awscdk.StackProps{
		StackName: &apiStackName,
		Env:       newCdkEnvironment(),
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func newCdkEnvironment() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	// 	Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	// 	Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
