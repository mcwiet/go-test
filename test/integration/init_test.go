package integration_test

import (
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/machinebox/graphql"
	"github.com/mcwiet/go-test/pkg/authentication"
)

var (
	Authenticator authentication.CognitoAuthenticator
	UserToken     authentication.UserToken
	GraphQlClient *graphql.Client
	TestUserEmail string
)

func init() {
	// Setup the GraphQL client
	apiUrl := GetRequiredEnv("API_URL")
	GraphQlClient = graphql.NewClient(apiUrl)

	// Get an access token for making API calls
	clientId := GetRequiredEnv("USER_POOL_APP_CLIENT_ID")
	session, _ := session.NewSession()
	cognitoClient := cognito.New(session)
	Authenticator = authentication.NewCognitoAuthenticator(cognitoClient, clientId)
	TestUserEmail = GetRequiredEnv("TEST_USER_EMAIL")
	password := GetRequiredEnv("TEST_USER_PASSWORD")
	UserToken, _ = Authenticator.Login(TestUserEmail, password)
}
