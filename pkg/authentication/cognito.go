package authentication

import (
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/jsii-runtime-go"
)

type CognitoAuthenticator struct {
	AppClientId string
	Provider    *cognito.CognitoIdentityProvider
}

// Login to the Cognito User Pool
func (a *CognitoAuthenticator) Login(email string, password string) (string, error) {
	authTry := &cognito.InitiateAuthInput{
		AuthFlow: jsii.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": &email,
			"PASSWORD": &password,
		},
		ClientId: jsii.String(a.AppClientId),
	}

	authResp, err := a.Provider.InitiateAuth(authTry)

	var token string
	if authResp != nil && authResp.AuthenticationResult != nil {
		token = *authResp.AuthenticationResult.AccessToken
	}

	return token, err
}
