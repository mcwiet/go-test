package infra

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsssm"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/constructs-go/constructs/v10"
)

var (
	ssmClient = new(ssm.SSM)
)

// Defining all parameter names here - helps prevent typos and keep parameters organized
const (
	ParamAppSyncUrl          = "appsync-url"
	ParamUserPoolId          = "user-pool-id"
	ParamUserPoolApiClientId = "user-pool-api-client-id"
)

func init() {
	session, _ := session.NewSession()
	ssmClient = ssm.New(session)
}

// Returns the name for an SSM parameter
func NewInfraParameter(scope constructs.Construct, envName string, paramDescriptor string, value string) awsssm.StringParameter {
	paramName := "/go/" + envName + "/" + paramDescriptor
	return awsssm.NewStringParameter(scope, &paramName, &awsssm.StringParameterProps{
		ParameterName: &paramName,
		StringValue:   &value,
	})
}

// Gets the value of an existing parameter
func GetInfraParameter(envName string, paramDescriptor string) string {
	paramName := "/go/" + envName + "/" + paramDescriptor
	output, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name: &paramName,
	})
	if err != nil {
		panic("Could not retrieve parameter: " + paramName)
	}
	return *output.Parameter.Value
}
