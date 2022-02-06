package integration_test

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/machinebox/graphql"
	"github.com/mcwiet/go-test/pkg/authentication"
	"github.com/mcwiet/go-test/pkg/model"
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
	apiUrl = GetRequiredEnv("API_URL")
	client = graphql.NewClient(apiUrl)

	// Get an access token for making API calls
	clientId := GetRequiredEnv("USER_POOL_APP_CLIENT_ID")
	session, _ := session.NewSession()
	cognitoClient := cognito.New(session)
	app := authentication.CognitoAuthenticator{
		AppClientId: clientId,
		Provider:    cognitoClient,
	}
	email := GetRequiredEnv("USER_EMAIL")
	password := GetRequiredEnv("USER_PASSWORD")
	token, _ = app.Login(email, password)
}

// Attempt to load environment variable; panic if not found (fail fast)
func GetRequiredEnv(name string) string {
	val, exists := os.LookupEnv(name)
	if exists == false {
		panic("Could not load environment variable: " + name)
	}
	return val
}

func TestCreatePerson(t *testing.T) {
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
	assert.Nil(t, err, "should not error")
	assert.NotNil(t, person.Id, "id should exist")
	assert.Equal(t, personName, person.Name, "name should match")
	assert.Equal(t, personAge, person.Age, "age should match")
}
