package data_test

import (
	"testing"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/jsii-runtime-go"
	"github.com/mcwiet/go-test/pkg/data"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/assert"
)

var (
	SampleUser1 = model.User{
		Username: "test-user-1",
		Email:    "email1@email.com",
		Name:     "Test User 1",
	}
	SampleUser2 = model.User{
		Username: "test-user-2",
		Email:    "email2@email.com",
		Name:     "Test User 2",
	}
	SampleUser1Attrs = []*cognito.AttributeType{
		{Name: jsii.String("email"), Value: &SampleUser1.Email},
		{Name: jsii.String("name"), Value: &SampleUser1.Name},
	}
	SampleUser2Attrs = []*cognito.AttributeType{
		{Name: jsii.String("email"), Value: &SampleUser2.Email},
		{Name: jsii.String("name"), Value: &SampleUser2.Name},
	}
	SamplePaginationToken = "test pagination token"
)

func TestUserGetByUsername(t *testing.T) {
	// Define test struct
	type Test struct {
		name           string
		userPoolClient FakeUserPoolClient
		username       string
		expectedUser   model.User
		expectErr      bool
	}

	// Define tests
	tests := []Test{
		{
			name: "valid get by username",
			userPoolClient: FakeUserPoolClient{
				adminGetUserOutput: &cognito.AdminGetUserOutput{
					UserAttributes: SampleUser1Attrs,
				},
			},
			username:     SampleUser1.Username,
			expectedUser: SampleUser1,
			expectErr:    false,
		},
		{
			name: "DAO get user error",
			userPoolClient: FakeUserPoolClient{
				adminGetUserErr: assert.AnError,
			},
			username:  SampleUser1.Username,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		userDao := data.NewUserDao(&test.userPoolClient, SampleUserPoolId)

		// Execute
		user, err := userDao.GetByUsername(test.username)

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
			assert.Equal(t, test.expectedUser, user, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestUserList(t *testing.T) {
	// Define test struct
	type Test struct {
		name           string
		userPoolClient FakeUserPoolClient
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
			userPoolClient: FakeUserPoolClient{
				listUsersOutput: &cognito.ListUsersOutput{
					Users: []*cognito.UserType{
						{Username: &SampleUser1.Username, Attributes: SampleUser1Attrs},
						{Username: &SampleUser2.Username, Attributes: SampleUser2Attrs},
					},
				},
			},
			first: 10,
			after: "",
			expectedUsers: []model.User{
				SampleUser1,
				SampleUser2,
			},
			expectedToken: "",
			expectErr:     false,
		},
		{
			name: "request some users (beginning of list)",
			userPoolClient: FakeUserPoolClient{
				listUsersOutput: &cognito.ListUsersOutput{
					Users: []*cognito.UserType{
						{Username: &SampleUser1.Username, Attributes: SampleUser1Attrs},
					},
					PaginationToken: &SamplePaginationToken,
				},
			},
			first: 1,
			after: "",
			expectedUsers: []model.User{
				SampleUser1,
			},
			expectedToken: SamplePaginationToken,
			expectErr:     false,
		},
		{
			name: "request some users (end of list)",
			userPoolClient: FakeUserPoolClient{
				listUsersOutput: &cognito.ListUsersOutput{
					Users: []*cognito.UserType{
						{Username: &SampleUser2.Username, Attributes: SampleUser2Attrs},
					},
				},
			},
			first: 1,
			after: SamplePaginationToken,
			expectedUsers: []model.User{
				SampleUser2,
			},
			expectedToken: "",
			expectErr:     false,
		},
		{
			name: "DAO list error",
			userPoolClient: FakeUserPoolClient{
				listUsersErr: assert.AnError,
			},
			first:     1,
			after:     "",
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		userDao := data.NewUserDao(&test.userPoolClient, SampleUserPoolId)

		// Execute
		users, token, err := userDao.List(test.first, test.after)

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
			assert.Equal(t, test.expectedUsers, users, test.name)
			assert.Equal(t, test.expectedToken, token, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestUserGetTotalCountUsers(t *testing.T) {
	// Define test struct
	type Test struct {
		name           string
		userPoolClient FakeUserPoolClient
		expectedCount  int
		expectErr      bool
	}

	// Define tests
	tests := []Test{
		{
			name: "valid get total count",
			userPoolClient: FakeUserPoolClient{
				describeUserPoolOutput: &cognito.DescribeUserPoolOutput{
					UserPool: &cognito.UserPoolType{
						EstimatedNumberOfUsers: pointy.Int64(2),
					},
				},
			},
			expectedCount: 2,
			expectErr:     false,
		},
		{
			name: "DAO describe user pool error",
			userPoolClient: FakeUserPoolClient{
				describeUserPoolErr: assert.AnError,
			},
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		userDao := data.NewUserDao(&test.userPoolClient, SampleUserPoolId)

		// Execute
		count, err := userDao.GetTotalCount()

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
			assert.Equal(t, test.expectedCount, count, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}
