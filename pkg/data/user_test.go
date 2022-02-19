package data_test

import (
	"testing"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/jsii-runtime-go"
	"github.com/mcwiet/go-test/pkg/data"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/stretchr/testify/assert"
)

type fakeUserPoolClient struct {
	adminGetUserOutput     *cognito.AdminGetUserOutput
	adminGetUserErr        error
	listUsersOutput        *cognito.ListUsersOutput
	listUsersErr           error
	describeUserPoolOutput *cognito.DescribeUserPoolOutput
	describeUserPoolErr    error
}

func (f *fakeUserPoolClient) AdminGetUser(*cognito.AdminGetUserInput) (*cognito.AdminGetUserOutput, error) {
	return f.adminGetUserOutput, f.adminGetUserErr
}

func (f *fakeUserPoolClient) ListUsers(*cognito.ListUsersInput) (*cognito.ListUsersOutput, error) {
	return f.listUsersOutput, f.listUsersErr
}

func (f *fakeUserPoolClient) DescribeUserPool(*cognito.DescribeUserPoolInput) (*cognito.DescribeUserPoolOutput, error) {
	return f.describeUserPoolOutput, f.describeUserPoolErr
}

var (
	sampleUserPoolId = "user-pool-id"
	sampleUser1      = model.User{
		Username: "test-user-1",
		Email:    "email1@email.com",
		Name:     "Test User 1",
	}
	sampleUser2 = model.User{
		Username: "test-user-2",
		Email:    "email2@email.com",
		Name:     "Test User 2",
	}
	sampleUserAttrs1 = []*cognito.AttributeType{
		{Name: jsii.String("email"), Value: &sampleUser1.Email},
		{Name: jsii.String("name"), Value: &sampleUser1.Name},
	}
	sampleUserAttrs2 = []*cognito.AttributeType{
		{Name: jsii.String("email"), Value: &sampleUser2.Email},
		{Name: jsii.String("name"), Value: &sampleUser2.Name},
	}
	samplePaginationToken = "test pagination token"
)

func TestGetByUsername(t *testing.T) {
	// Define test struct
	type Test struct {
		name           string
		userPoolClient fakeUserPoolClient
		username       string
		expectedUser   model.User
		expectErr      bool
	}

	// Define tests
	tests := []Test{
		{
			name: "valid get by username",
			userPoolClient: fakeUserPoolClient{
				adminGetUserOutput: &cognito.AdminGetUserOutput{
					UserAttributes: sampleUserAttrs1}},
			username:     sampleUser1.Username,
			expectedUser: sampleUser1,
			expectErr:    false,
		},
		{
			name: "DAO get user error",
			userPoolClient: fakeUserPoolClient{
				adminGetUserErr: assert.AnError},
			username:  sampleUser1.Username,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		userDao := data.NewUserDao(&test.userPoolClient, sampleUserPoolId)

		// Execute
		user, err := userDao.GetByUsername(test.username)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedUser, user, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

// func TestQuery
func TestList(t *testing.T) {
	// Define test struct
	type Test struct {
		name           string
		userPoolClient fakeUserPoolClient
		first          int
		after          string
		expectedUsers  []model.User
		expectedToken  string
		expectErr      bool
	}

	// Define tests
	tests := []Test{
		{
			name: "request more users than in DB",
			userPoolClient: fakeUserPoolClient{
				listUsersOutput: &cognito.ListUsersOutput{
					Users: []*cognito.UserType{
						{Username: &sampleUser1.Username, Attributes: sampleUserAttrs1},
						{Username: &sampleUser2.Username, Attributes: sampleUserAttrs2}}}},
			first: 10,
			after: "",
			expectedUsers: []model.User{
				sampleUser1,
				sampleUser2,
			},
			expectedToken: "",
			expectErr:     false,
		},
		{
			name: "request some users (beginning of list)",
			userPoolClient: fakeUserPoolClient{
				listUsersOutput: &cognito.ListUsersOutput{
					Users: []*cognito.UserType{
						{Username: &sampleUser1.Username, Attributes: sampleUserAttrs1}},
					PaginationToken: &samplePaginationToken}},
			first: 1,
			after: "",
			expectedUsers: []model.User{
				sampleUser1},
			expectedToken: samplePaginationToken,
			expectErr:     false,
		},
		{
			name: "request some users (end of list)",
			userPoolClient: fakeUserPoolClient{
				listUsersOutput: &cognito.ListUsersOutput{
					Users: []*cognito.UserType{
						{Username: &sampleUser2.Username, Attributes: sampleUserAttrs2}}}},
			first: 1,
			after: samplePaginationToken,
			expectedUsers: []model.User{
				sampleUser2},
			expectedToken: "",
			expectErr:     false,
		},
		{
			name: "DAO list error",
			userPoolClient: fakeUserPoolClient{
				listUsersErr: assert.AnError},
			first:     1,
			after:     "",
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		userDao := data.NewUserDao(&test.userPoolClient, sampleUserPoolId)

		// Execute
		users, token, err := userDao.List(test.first, test.after)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedUsers, users, test.name)
			assert.Equal(t, test.expectedToken, token, test.name)
			assert.Nil(t, err, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestGetTotalCountUsers(t *testing.T) {
	// Define test struct
	type Test struct {
		name           string
		userPoolClient fakeUserPoolClient
		expectedCount  int
		expectErr      bool
	}

	// Define tests
	tests := []Test{
		{
			name: "valid get total count",
			userPoolClient: fakeUserPoolClient{
				describeUserPoolOutput: &cognito.DescribeUserPoolOutput{
					UserPool: &cognito.UserPoolType{
						EstimatedNumberOfUsers: newInt64Temp(2)}}},
			expectedCount: 2,
			expectErr:     false,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		userDao := data.NewUserDao(&test.userPoolClient, sampleUserPoolId)

		// Execute
		count, err := userDao.GetTotalCount()

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedCount, count, test.name)
			assert.Nil(t, err, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func newInt64Temp(val int) *int64 {
	valInt64 := int64(val)
	return &valInt64
}
