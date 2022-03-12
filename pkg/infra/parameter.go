package infra

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsssm"
	"github.com/aws/constructs-go/constructs/v10"
)

// Defining all parameter names here - helps prevent typos and keep parameters organized
const (
	ParamAppSyncId                  = "appsync-id"
	ParamAppSyncUrl                 = "appsync-url"
	ParamAuthenticatedUserRoleArn   = "authenticated-user-role-arn"
	ParamIdentityPoolId             = "identity-pool-id"
	ParamUserPoolArn                = "user-pool-arn"
	ParamUserPoolId                 = "user-pool-id"
	ParamUserPoolApiClientId        = "user-pool-api-client-id"
	ParamUnauthenticatedUserRoleArn = "unauthenticated-user-role-arn"
)

// Returns the name for an SSM parameter
func NewInfraParameter(scope constructs.Construct, envName string, paramDescriptor string, value string) awsssm.StringParameter {
	paramName := "/go/" + envName + "/" + paramDescriptor
	return awsssm.NewStringParameter(scope, &paramName, &awsssm.StringParameterProps{
		ParameterName: &paramName,
		StringValue:   &value,
	})
}

// Gets the value of an existing parameter
func GetInfraParameter(scope constructs.Construct, envName string, paramDescriptor string) string {
	paramName := "/go/" + envName + "/" + paramDescriptor
	paramToken := awsssm.StringParameter_ValueForStringParameter(scope, &paramName, nil)
	return *paramToken
}
