package service_test

import (
	"testing"

	"github.com/mcwiet/go-test/pkg/model"
	"github.com/mcwiet/go-test/pkg/service"
	"github.com/stretchr/testify/assert"
)

var (
	SampleUser1 = model.User{
		Username: "user-1",
		Email:    "user1@email.com",
		Name:     "User 1",
	}
	SampleUser2 = model.User{
		Username: "user-2",
		Email:    "user2@email.com",
		Name:     "User 2",
	}
	SampleUser1Edge = model.UserEdge{
		Node: SampleUser1,
	}
	SampleUser2Edge = model.UserEdge{
		Node: SampleUser2,
	}
)

func TestUserGetByUsername(t *testing.T) {
	type Test struct {
		name         string
		userDao      FakeUserDao
		username     string
		expectedUser model.User
		expectErr    bool
	}

	tests := []Test{
		{
			name: "valid get by username",
			userDao: FakeUserDao{
				getByUsernameUser: SampleUser1,
			},
			username:     SampleUser1.Username,
			expectedUser: SampleUser1,
			expectErr:    false,
		},
		{
			name: "DAO get by username error",
			userDao: FakeUserDao{
				getByUsernameErr: assert.AnError,
			},
			username:  SampleUser1.Username,
			expectErr: true,
		},
	}

	for _, test := range tests {
		service := service.NewUserService(&test.userDao, &SampleEncoder)

		user, err := service.GetByUsername(test.username)

		if !test.expectErr {
			assert.Nil(t, err, test.name)
			assert.Equal(t, test.expectedUser, user, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestUserList(t *testing.T) {
	type Test struct {
		name               string
		userDao            FakeUserDao
		encoder            FakeEncoder
		first              int
		after              string
		expectedConnection model.UserConnection
		expectErr          bool
	}

	tests := []Test{
		{
			name: "list all users",
			userDao: FakeUserDao{
				listUsers: []model.User{
					SampleUser1,
					SampleUser2,
				},
				listToken:          "",
				getTotalCountValue: 2,
			},
			encoder: SampleEncoder,
			first:   10,
			after:   "",
			expectedConnection: model.UserConnection{
				TotalCount: 2,
				Edges: []model.UserEdge{
					SampleUser1Edge,
					SampleUser2Edge,
				},
				PageInfo: model.PageInfo{
					EndCursor:   "",
					HasNextPage: false,
				},
			},
		},
		{
			name: "list first of two users",
			userDao: FakeUserDao{
				listUsers: []model.User{
					SampleUser1,
				},
				listToken:          "token",
				getTotalCountValue: 2,
			},
			encoder: SampleEncoder,
			first:   1,
			after:   "",
			expectedConnection: model.UserConnection{
				TotalCount: 2,
				Edges: []model.UserEdge{
					SampleUser1Edge,
				},
				PageInfo: model.PageInfo{
					EndCursor:   SampleEncoder.Encode("token"),
					HasNextPage: true,
				},
			},
		},
		{
			name: "list second of two users",
			userDao: FakeUserDao{
				listUsers: []model.User{
					SampleUser2,
				},
				listToken:          "",
				getTotalCountValue: 2,
			},
			encoder: SampleEncoder,
			first:   1,
			after:   SampleEncoder.Encode("token"),
			expectedConnection: model.UserConnection{
				TotalCount: 2,
				Edges: []model.UserEdge{
					SampleUser2Edge,
				},
				PageInfo: model.PageInfo{
					EndCursor:   "",
					HasNextPage: false,
				},
			},
		},
		{
			name: "decode error",
			userDao: FakeUserDao{
				getTotalCountErr: assert.AnError,
			},
			encoder:   SampleEncoder,
			first:     1,
			after:     "",
			expectErr: true,
		},
		{
			name: "DAO get total count error",
			userDao: FakeUserDao{
				getTotalCountErr: assert.AnError,
			},
			encoder:   SampleEncoder,
			first:     1,
			after:     "",
			expectErr: true,
		},
		{
			name: "DAO list error",
			userDao: FakeUserDao{
				listErr: assert.AnError,
			},
			encoder:   SampleEncoder,
			first:     1,
			after:     "",
			expectErr: true,
		},
	}

	for _, test := range tests {
		// Setup
		service := service.NewUserService(&test.userDao, &test.encoder)

		// Execute
		connection, err := service.List(test.first, test.after)

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
			assert.Equal(t, test.expectedConnection, connection, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}
