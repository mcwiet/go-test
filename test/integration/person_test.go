package integration_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/machinebox/graphql"
	"github.com/mcwiet/go-test/pkg/authentication"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/mcwiet/go-test/test/integration"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

var (
	apiUrl string
	token  string
	client *graphql.Client
)

func init() {
	// Setup the GraphQL client
	apiUrl = integration.GetRequiredEnv("API_URL")
	client = graphql.NewClient(apiUrl)

	// Get an access token for making API calls
	clientId := integration.GetRequiredEnv("USER_POOL_APP_CLIENT_ID")
	session, _ := session.NewSession()
	cognitoClient := cognito.New(session)
	auth := authentication.NewCognitoAuthenticator(cognitoClient, clientId)
	email := integration.GetRequiredEnv("TEST_USER_EMAIL")
	password := integration.GetRequiredEnv("TEST_USER_PASSWORD")
	token, _ = auth.Login(email, password)
}

// Sequentially run functions involved for testing person API operations
func TestPersonApi(t *testing.T) {

	// Create some people
	person1 := createPerson(t)
	person2 := createPerson(t)

	// List the people
	listPeople(t)

	// Get a person
	getPerson(t, person1.Id, &person1)

	// Delete the people
	deletePerson(t, &person1)
	deletePerson(t, &person2)

	// Attempt to get a person again
	getPerson(t, person1.Id, nil)
}

func createPerson(t *testing.T) model.Person {
	// Setup
	personName := "Integration Test"
	personAge := 10
	request := graphql.NewRequest(`
		mutation ($name: String!, $age: Int!) {
			createPerson (name: $name, age: $age) {
				id
				name
				age
			}
		}
	`)
	request.Var("name", personName)
	request.Var("age", personAge)
	request.Header.Set("Authorization", token)

	// Execute
	var response map[string]interface{}
	err := client.Run(context.Background(), request, &response)
	var person model.Person
	mapstructure.Decode(response["createPerson"], &person)

	// Verify
	stepName := "createPerson"
	assert.Nil(t, err, stepName+": should not error")
	if err != nil {
		return model.Person{}
	}
	assert.NotNil(t, person.Id, stepName+"id should exist")
	assert.Equal(t, personName, person.Name, stepName+"name should match")
	assert.Equal(t, personAge, person.Age, stepName+"age should match")

	return person
}

func deletePerson(t *testing.T, person *model.Person) {
	// Setup
	request := graphql.NewRequest(`
		mutation ($id: ID!) {
			deletePerson (id: $id)
		}
	`)
	request.Var("id", person.Id)
	request.Header.Set("Authorization", token)

	// Execute
	var response map[string]interface{}
	err := client.Run(context.Background(), request, &response)

	// Verify
	stepName := "deletePerson"
	assert.Nil(t, err, stepName+": should not error")
}

func getPerson(t *testing.T, id string, expectedPerson *model.Person) {
	// Setup
	request := graphql.NewRequest(`
		query ($id: ID!) {
			person (id: $id) {
				id
				name
				age
			}
		}
	`)
	request.Var("id", id)
	request.Header.Set("Authorization", token)

	// Execute
	var response map[string]interface{}
	err := client.Run(context.Background(), request, &response)
	var person model.Person
	mapstructure.Decode(response["person"], &person)

	// Verify
	stepName := "getPerson"
	if expectedPerson != nil {
		assert.Nil(t, err, stepName+": should not error")
		assert.Equal(t, *expectedPerson, person, stepName+": should find the correct person")
	} else {
		assert.NotNil(t, err, stepName+": should not find person with id "+id)
		assert.Equal(t, "", person.Id, stepName+": should not find person with id "+id)
	}
}

func listPeople(t *testing.T) {
	// Setup
	request := graphql.NewRequest(`
		query {
			people (first: 1, after: "") {
				totalCount
				edges {
					node {
						id
					}
					cursor
				}
				pageInfo {
					endCursor
					hasNextPage
				}
			}
		}
	`)
	request.Header.Set("Authorization", token)

	// Execute
	var response map[string]interface{}
	err := client.Run(context.Background(), request, &response)
	var connection model.PersonConnection
	mapstructure.Decode(response["people"], &connection)

	// Verify
	stepName := "listPeople"
	assert.Nil(t, err, stepName+": should not error")
	if err != nil {
		return
	}
	assert.Equal(t, 1, len(connection.Edges), stepName+": should return 1 person")
	assert.GreaterOrEqual(t, connection.TotalCount, 2, stepName+": should have total count of at least 2")
	lastEdge := connection.Edges[len(connection.Edges)-1]
	assert.Equal(t, lastEdge.Cursor, connection.PageInfo.EndCursor, stepName+": should have correct end cursor")
	assert.Equal(t, true, connection.PageInfo.HasNextPage, connection.PageInfo.EndCursor, stepName+": should have next page")
}
