package data_test

import (
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	SampleTableName  = "table"
	SampleUserPoolId = "user-pool-id"
)

type FakeDynamoDbClient struct {
	deleteItemOutput *dynamodb.DeleteItemOutput
	deleteItemErr    error
	getItemOutput    *dynamodb.GetItemOutput
	getItemErr       error
	putItemOutput    *dynamodb.PutItemOutput
	putItemErr       error
	queryOutput      *dynamodb.QueryOutput
	queryErr         error
}

func (f *FakeDynamoDbClient) DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return f.deleteItemOutput, f.deleteItemErr
}
func (f *FakeDynamoDbClient) GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return f.getItemOutput, f.getItemErr
}
func (f *FakeDynamoDbClient) PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return f.putItemOutput, f.putItemErr
}
func (f *FakeDynamoDbClient) Query(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return f.queryOutput, f.queryErr
}

type FakeUserPoolClient struct {
	adminGetUserOutput     *cognito.AdminGetUserOutput
	adminGetUserErr        error
	listUsersOutput        *cognito.ListUsersOutput
	listUsersErr           error
	describeUserPoolOutput *cognito.DescribeUserPoolOutput
	describeUserPoolErr    error
}

func (f *FakeUserPoolClient) AdminGetUser(*cognito.AdminGetUserInput) (*cognito.AdminGetUserOutput, error) {
	return f.adminGetUserOutput, f.adminGetUserErr
}
func (f *FakeUserPoolClient) ListUsers(*cognito.ListUsersInput) (*cognito.ListUsersOutput, error) {
	return f.listUsersOutput, f.listUsersErr
}
func (f *FakeUserPoolClient) DescribeUserPool(*cognito.DescribeUserPoolInput) (*cognito.DescribeUserPoolOutput, error) {
	return f.describeUserPoolOutput, f.describeUserPoolErr
}
