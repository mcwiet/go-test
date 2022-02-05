package infra

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AuthStackProps struct {
	awscdk.StackProps
}

func NewAuthStack(scope constructs.Construct, id string, props *AuthStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)
	stackName := props.StackName

	// Cognito User Pool
	userPoolName := *stackName + "-user-pool"
	userPool := awscognito.NewUserPool(stack, &userPoolName, &awscognito.UserPoolProps{
		UserPoolName: &userPoolName,
		SignInAliases: &awscognito.SignInAliases{
			Email:    jsii.Bool(true),
			Username: jsii.Bool(false),
		},
		StandardAttributes: &awscognito.StandardAttributes{
			Email: &awscognito.StandardAttribute{
				Required: jsii.Bool(true),
			},
		},
	})
	NewInfraParameter(stack, userPoolName, "id", *userPool.UserPoolId())

	return stack
}
