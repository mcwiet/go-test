package authentication

import (
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/jsii-runtime-go"
)

type CognitoIdentityProvider interface {
	InitiateAuth(*cognito.InitiateAuthInput) (*cognito.InitiateAuthOutput, error)
}

type CognitoAuthenticator struct {
	provider    CognitoIdentityProvider
	appClientId string
}

// Creates a new authenticator object
func NewCognitoAuthenticator(provider CognitoIdentityProvider, appClientId string) CognitoAuthenticator {
	return CognitoAuthenticator{
		provider:    provider,
		appClientId: appClientId,
	}
}

// Login to the Cognito User Pool
func (a *CognitoAuthenticator) Login(email string, password string) (string, error) {
	authTry := &cognito.InitiateAuthInput{
		AuthFlow: jsii.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": &email,
			"PASSWORD": &password,
		},
		ClientId: jsii.String(a.appClientId),
	}

	authResp, err := a.provider.InitiateAuth(authTry)

	var token string
	if authResp != nil && authResp.AuthenticationResult != nil {
		token = *authResp.AuthenticationResult.AccessToken
	}

	return token, err
}
