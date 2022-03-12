package infra

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
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

	// Cognito User Pool Admin Group
	userPoolAdminGroupName := *stackName + "-user-pool-admin-group"
	awscognito.NewCfnUserPoolGroup(stack, &userPoolAdminGroupName, &awscognito.CfnUserPoolGroupProps{
		GroupName:   jsii.String("admin"),
		Description: jsii.String("Application administrators"),
		UserPoolId:  userPool.UserPoolId(),
	})

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

	// Cognito Identity Pool
	identityPoolName := *stackName + "-identity-pool"
	identityPool := awscognito.NewCfnIdentityPool(stack, &identityPoolName, &awscognito.CfnIdentityPoolProps{
		AllowUnauthenticatedIdentities: true,
		CognitoIdentityProviders: []map[string]*string{{
			"clientId":     appClient.UserPoolClientId(),
			"providerName": userPool.UserPoolProviderName(),
		}},
		IdentityPoolName: &identityPoolName,
	})
	NewInfraParameter(stack, props.EnvName, ParamIdentityPoolId, *identityPool.Ref())

	// Unauthenticated User Role
	unauthRoleName := *stackName + "-unauthenticated-user-role"
	unauthRole := awsiam.NewRole(stack, &unauthRoleName, &awsiam.RoleProps{
		RoleName: &unauthRoleName,
		AssumedBy: awsiam.NewFederatedPrincipal(jsii.String("cognito-identity.amazonaws.com"), &map[string]interface{}{
			"StringEquals":           map[string]interface{}{"cognito-identity.amazonaws.com:aud": identityPool.Ref()},
			"ForAnyValue:StringLike": map[string]interface{}{"cognito-identity.amazonaws.com:amr": "unauthenticated"},
		}, jsii.String("sts:AssumeRoleWithWebIdentity")),
	})
	NewInfraParameter(stack, props.EnvName, ParamUnauthenticatedUserRoleArn, *unauthRole.RoleArn())

	// Cognito Identity Pool Role Attachment
	identityPoolRoleAttachmentName := *stackName + "identity-pool-role-attachment"
	awscognito.NewCfnIdentityPoolRoleAttachment(stack, &identityPoolRoleAttachmentName, &awscognito.CfnIdentityPoolRoleAttachmentProps{
		IdentityPoolId: identityPool.Ref(),
		Roles: map[string]*string{
			"unauthenticated": unauthRole.RoleArn(),
		},
	})

	return stack
}
