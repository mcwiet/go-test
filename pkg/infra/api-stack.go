package infra

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdkappsyncalpha/v2"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ApiStackProps struct {
	awscdk.StackProps
	EnvName string
}

func NewApiStack(scope constructs.Construct, id string, props *ApiStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)
	stackName := props.StackName

	// Lambda resolver
	lambdaName := *stackName + "-lambda"
	lambda := awscdklambdagoalpha.NewGoFunction(stack, &lambdaName, &awscdklambdagoalpha.GoFunctionProps{
		Entry:        jsii.String("./cmd/api"),
		FunctionName: &lambdaName,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(5)),
		Tracing:      awslambda.Tracing_ACTIVE,
	})

	// Schema definition
	schema := awscdkappsyncalpha.Schema_FromAsset(jsii.String("./api/schema.graphql"))

	// AppSync API
	apiName := *stackName + "-appsync"
	userPoolId := GetInfraParameter(stack, props.EnvName, ParamUserPoolId)
	userPool := awscognito.UserPool_FromUserPoolId(stack, &userPoolId, &userPoolId)
	api := awscdkappsyncalpha.NewGraphqlApi(stack, &apiName, &awscdkappsyncalpha.GraphqlApiProps{
		Name:   &apiName,
		Schema: schema,
		AuthorizationConfig: &awscdkappsyncalpha.AuthorizationConfig{
			DefaultAuthorization: &awscdkappsyncalpha.AuthorizationMode{
				AuthorizationType: awscdkappsyncalpha.AuthorizationType_USER_POOL,
				UserPoolConfig: &awscdkappsyncalpha.UserPoolConfig{
					UserPool: userPool,
				},
			},
		},
	})
	NewInfraParameter(stack, props.EnvName, ParamAppSyncUrl, *api.GraphqlUrl())

	// Data source(s)
	apiSourceName := "lambda_source" // Can't use '-' in data source name
	lambdaSource := api.AddLambdaDataSource(&apiSourceName, lambda, &awscdkappsyncalpha.DataSourceOptions{
		Name: &apiSourceName,
	})

	// Resolvers
	createResolver(api, "Query", "pet", lambdaSource)
	createResolver(api, "Query", "pets", lambdaSource)
	createResolver(api, "Query", "user", lambdaSource)
	createResolver(api, "Query", "users", lambdaSource)
	createResolver(api, "Mutation", "createPet", lambdaSource)
	createResolver(api, "Mutation", "deletePet", lambdaSource)
	createResolver(api, "Mutation", "updatePetOwner", lambdaSource)

	// Dynamo DB table
	tableName := *stackName + "-primary-table"
	partitionKey := awsdynamodb.Attribute{Name: jsii.String("Id"), Type: awsdynamodb.AttributeType_STRING}
	sortKey := awsdynamodb.Attribute{Name: jsii.String("Sort"), Type: awsdynamodb.AttributeType_STRING}
	table := awsdynamodb.NewTable(stack, &tableName, &awsdynamodb.TableProps{
		TableName:    &tableName,
		PartitionKey: &partitionKey,
		SortKey:      &sortKey,
		BillingMode:  awsdynamodb.BillingMode_PAY_PER_REQUEST,
	})
	table.AddGlobalSecondaryIndex(&awsdynamodb.GlobalSecondaryIndexProps{
		IndexName:      jsii.String("sort-key-gsi"),
		ProjectionType: awsdynamodb.ProjectionType_ALL,
		PartitionKey:   &sortKey,
	})

	// Permission for Lambda to access Dynamo DB table
	lambda.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings(
			"dynamodb:BatchGetItem",
			"dynamodb:DescribeTable",
			"dynamodb:GetItem",
			"dynamodb:Query",
			"dynamodb:Scan",
			"dynamodb:BatchWriteItem",
			"dynamodb:DeleteItem",
			"dynamodb:UpdateItem",
			"dynamodb:PutItem"),
		Resources: jsii.Strings(*table.TableArn(), *table.TableArn()+"/*"),
	}))

	// Permission for Lambda to access Cognito User Pool
	userPoolArn := GetInfraParameter(stack, props.EnvName, ParamUserPoolArn)
	lambda.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings(
			"cognito-idp:AdminGetUser",
			"cognito-idp:ListUsers",
			"cognito-idp:DescribeUserPool",
		),
		Resources: jsii.Strings(userPoolArn, userPoolArn+"/*"),
	}))

	// Add environment variables to Lambda to reference other infra
	lambda.AddEnvironment(jsii.String("DDB_TABLE_NAME"), &tableName, nil)
	lambda.AddEnvironment(jsii.String("USER_POOL_ID"), &userPoolId, nil)

	return stack
}

func createResolver(api awscdkappsyncalpha.GraphqlApi, typeName string, fieldName string, source awscdkappsyncalpha.BaseDataSource) {
	api.CreateResolver(&awscdkappsyncalpha.ExtendedResolverProps{
		TypeName:   &typeName,
		FieldName:  &fieldName,
		DataSource: source,
	})
}
