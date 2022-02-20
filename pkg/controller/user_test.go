package controller_test

import (
	"testing"

	"github.com/mcwiet/go-test/pkg/controller"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/stretchr/testify/assert"
)

var (
	SampleUser = model.User{
		Username: "user-1",
		Email:    "user1@email.com",
		Name:     "User 1",
	}
	SampleUserConnection = model.UserConnection{
		TotalCount: 1,
		Edges: []model.UserEdge{
			{Node: SampleUser},
		},
		PageInfo: model.PageInfo{
			EndCursor:   "cursor",
			HasNextPage: false,
		},
	}
)

type UserTest struct {
	name             string
	userService      FakeUserService
	request          controller.Request
	expectedResponse controller.Response
	expectErr        bool
}

func TestUserHandleGet(t *testing.T) {
	tests := []UserTest{
		{
			name: "valid get",
			userService: FakeUserService{
				getByUsernameUser: SampleUser,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"username": SampleUser.Username,
				}},
			},
			expectedResponse: controller.Response{
				Data: SampleUser,
			},
		},
		{
			name: "service get by username error",
			userService: FakeUserService{
				getByUsernameErr: assert.AnError,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"username": SampleUser.Username,
				}},
			},
			expectErr: true,
		},
	}

	for _, test := range tests {
		// Setup
		controller := controller.NewUserController(&test.userService)

		// Execute
		response := controller.HandleGet(test.request)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedResponse, response)
		} else {
			assert.NotNil(t, response.Error, test.name)
		}
	}
}

func TestUserHandleList(t *testing.T) {
	tests := []UserTest{
		{
			name: "valid list",
			userService: FakeUserService{
				listConnection: SampleUserConnection,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"first": 10,
					"after": "token",
				}},
			},
			expectedResponse: controller.Response{
				Data: SampleUserConnection,
			},
		},
		{
			name: "service list error",
			userService: FakeUserService{
				listErr: assert.AnError,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"username": SampleUser.Username,
				}},
			},
			expectErr: true,
		},
	}

	for _, test := range tests {
		// Setup
		controller := controller.NewUserController(&test.userService)

		// Execute
		response := controller.HandleList(test.request)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedResponse, response)
		} else {
			assert.NotNil(t, response.Error, test.name)
		}
	}
}
