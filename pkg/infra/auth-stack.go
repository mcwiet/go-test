package infra

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AuthStackProps struct {
	awscdk.StackProps
	EnvName string
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
	NewInfraParameter(stack, props.EnvName, ParamUserPoolArn, *userPool.UserPoolArn())
	NewInfraParameter(stack, props.EnvName, ParamUserPoolId, *userPool.UserPoolId())

	// API App Client
	appClientName := userPoolName + "-api-client"
	appClient := userPool.AddClient(&appClientName, &awscognito.UserPoolClientOptions{
		UserPoolClientName: &appClientName,
		GenerateSecret:     jsii.Bool(false),
		OAuth: &awscognito.OAuthSettings{
			Flows: &awscognito.OAuthFlows{
				ImplicitCodeGrant: jsii.Bool(true),
			},
			Scopes: &[]awscognito.OAuthScope{
				awscognito.OAuthScope_PROFILE(),
				awscognito.OAuthScope_EMAIL(),
			},
		},
		AuthFlows: &awscognito.AuthFlow{
			UserPassword: jsii.Bool(true),
			UserSrp:      jsii.Bool(true),
		},
	})
	NewInfraParameter(stack, props.EnvName, ParamUserPoolApiClientId, *appClient.UserPoolClientId())

	return stack
}
