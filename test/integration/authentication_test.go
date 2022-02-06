package integration_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/mcwiet/go-test/pkg/authentication"
	"github.com/mcwiet/go-test/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	// Setup
	clientId := integration.GetRequiredEnv("USER_POOL_APP_CLIENT_ID")
	session, _ := session.NewSession()
	cognitoClient := cognito.New(session)
	auth := authentication.NewCognitoAuthenticator(cognitoClient, clientId)
	email := integration.GetRequiredEnv("USER_EMAIL")
	password := integration.GetRequiredEnv("USER_PASSWORD")

	// Test
	token, err := auth.Login(email, password)

	// Verify
	assert.Nil(t, err)
	assert.NotNil(t, token)
}
