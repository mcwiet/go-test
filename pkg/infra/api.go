package infra

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdkappsyncalpha/v2"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ApiStackProps struct {
	awscdk.StackProps
}

func NewApiStack(scope constructs.Construct, id string, props *awscdk.StackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, props)
	stackName := props.StackName

	// Lambda resolver
	lambdaName := *stackName + "-lambda"
	lambda := awscdklambdagoalpha.NewGoFunction(stack, &lambdaName, &awscdklambdagoalpha.GoFunctionProps{
		Entry:        jsii.String("./cmd/api"),
		FunctionName: &lambdaName,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(5)),
	})

	// Schema definition
	schema := awscdkappsyncalpha.Schema_FromAsset(jsii.String("./api/schema.graphql"))

	// AppSync API
	apiName := *stackName + "-appsync"
	api := awscdkappsyncalpha.NewGraphqlApi(stack, &apiName, &awscdkappsyncalpha.GraphqlApiProps{
		Name:   &apiName,
		Schema: schema,
	})

	// Data source(s)
	apiSourceName := "lambda_source" // Can't use '-' in data source name
	lambdaSource := api.AddLambdaDataSource(&apiSourceName, lambda, &awscdkappsyncalpha.DataSourceOptions{
		Name: &apiSourceName,
	})

	// Resolvers
	api.CreateResolver(&awscdkappsyncalpha.ExtendedResolverProps{
		TypeName:   jsii.String("Query"),
		FieldName:  jsii.String("person"),
		DataSource: lambdaSource,
	})
	api.CreateResolver(&awscdkappsyncalpha.ExtendedResolverProps{
		TypeName:   jsii.String("Query"),
		FieldName:  jsii.String("people"),
		DataSource: lambdaSource,
	})

	return stack
}
