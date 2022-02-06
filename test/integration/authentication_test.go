package integration_test

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/mcwiet/go-test/pkg/authentication"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	clientId := os.Getenv("USER_POOL_APP_CLIENT_ID")
	session, _ := session.NewSession()
	cognitoClient := cognito.New(session)

	app := authentication.CognitoAuthenticator{
		AppClientId: clientId,
		Provider:    cognitoClient,
	}

	email := os.Getenv("USER_EMAIL")
	password := os.Getenv("USER_PASSWORD")

	token, err := app.Login(email, password)

	assert.Nil(t, err)
	assert.NotNil(t, token)
}
