package integration_test

import (
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/machinebox/graphql"
	"github.com/mcwiet/go-test/pkg/authentication"
	"github.com/mcwiet/go-test/test/integration"
)

var (
	Authenticator authentication.CognitoAuthenticator
	AuthToken     authentication.CognitoToken
	GraphQlClient *graphql.Client
	TestUserEmail string
)

func init() {
	// Setup the GraphQL client
	apiUrl := integration.GetRequiredEnv("API_URL")
	GraphQlClient = graphql.NewClient(apiUrl)

	// Get an access token for making API calls
	clientId := integration.GetRequiredEnv("USER_POOL_APP_CLIENT_ID")
	session, _ := session.NewSession()
	cognitoClient := cognito.New(session)
	Authenticator = authentication.NewCognitoAuthenticator(cognitoClient, clientId)
	TestUserEmail = integration.GetRequiredEnv("TEST_USER_EMAIL")
	password := integration.GetRequiredEnv("TEST_USER_PASSWORD")
	AuthToken, _ = Authenticator.Login(TestUserEmail, password)
}
