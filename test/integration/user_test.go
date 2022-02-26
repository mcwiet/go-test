package integration_test

import (
	"context"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

// Sequentially run functions involved for testing user API operations
func TestUserApi(t *testing.T) {
	// List users
	listUsers(t)

	// Get a user
	getUser(t, UserToken.Username)
}

func getUser(t *testing.T, username string) {
	// Setup
	request := graphql.NewRequest(`
		query ($username: String!) {
			user (input: { username: $username }) {
				username
				email
			}
		}
	`)
	request.Var("username", username)
	request.Header.Set("Authorization", UserToken.AccessTokenString)

	// Execute
	var response map[string]interface{}
	err := GraphQlClient.Run(context.Background(), request, &response)
	var user model.User
	mapstructure.Decode(response["user"], &user)

	// Verify
	stepName := "getUser"
	assert.Nil(t, err, stepName+": should not error")
	assert.Equal(t, username, user.Username, stepName+": should find the correct user (username)")
	assert.Equal(t, TestUserEmail, user.Email, stepName+": should find the correct user (email)")
}

func listUsers(t *testing.T) {
	// Setup
	request := graphql.NewRequest(`
		query {
			users (input: { first: 1, after: "" }) {
				totalCount
				edges {
					node {
						username
					}
				}
				pageInfo {
					endCursor
					hasNextPage
				}
			}
		}
	`)
	request.Header.Set("Authorization", UserToken.AccessTokenString)

	// Execute
	var response map[string]interface{}
	err := GraphQlClient.Run(context.Background(), request, &response)
	var connection model.UserConnection
	mapstructure.Decode(response["users"], &connection)

	// Verify
	stepName := "listUsers"
	assert.Nil(t, err, stepName+": should not error")
	assert.Equal(t, 1, len(connection.Edges), stepName+": should return 1 user")
	assert.GreaterOrEqual(t, connection.TotalCount, 1, stepName+": should have total count of at least 1")
}
