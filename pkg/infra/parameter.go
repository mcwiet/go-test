package infra

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsssm"
	"github.com/aws/constructs-go/constructs/v10"
)

// Returns the name for an SSM parameter
func NewInfraParameter(scope constructs.Construct, targetConstructName string, descriptor string, value string) awsssm.StringParameter {
	paramName := "/go/" + targetConstructName + "/" + descriptor
	return awsssm.NewStringParameter(scope, &paramName, &awsssm.StringParameterProps{
		ParameterName: &paramName,
		StringValue:   &value,
	})
}
