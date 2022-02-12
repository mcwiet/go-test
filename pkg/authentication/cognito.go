package authentication

import (
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/jsii-runtime-go"
	"github.com/dgrijalva/jwt-go"
)

type CognitoIdentityProvider interface {
	InitiateAuth(*cognito.InitiateAuthInput) (*cognito.InitiateAuthOutput, error)
}

type CognitoAuthenticator struct {
	provider    CognitoIdentityProvider
	appClientId string
}

type CognitoTokenPayload struct {
	Username string
}

type CognitoToken struct {
	String  string
	Payload CognitoTokenPayload
}

// Creates a new authenticator object
func NewCognitoAuthenticator(provider CognitoIdentityProvider, appClientId string) CognitoAuthenticator {
	return CognitoAuthenticator{
		provider:    provider,
		appClientId: appClientId,
	}
}

// Login to the Cognito User Pool
func (a *CognitoAuthenticator) Login(email string, password string) (CognitoToken, error) {
	authTry := &cognito.InitiateAuthInput{
		AuthFlow: jsii.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": &email,
			"PASSWORD": &password,
		},
		ClientId: jsii.String(a.appClientId),
	}

	authResp, err := a.provider.InitiateAuth(authTry)
	if err != nil {
		return CognitoToken{}, err
	}

	var tokenStr string
	if authResp != nil && authResp.AuthenticationResult != nil {
		tokenStr = *authResp.AuthenticationResult.AccessToken
	}

	return getTokenFromString(tokenStr)
}

// Turn a token string into a token object (does not verify token!)
func getTokenFromString(tokenStr string) (CognitoToken, error) {
	decodedToken, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return CognitoToken{}, err
	}

	claims, ok := decodedToken.Claims.(jwt.MapClaims)
	if !ok {
		return CognitoToken{}, err
	}

	token := CognitoToken{
		String: tokenStr,
		Payload: CognitoTokenPayload{
			Username: claims["username"].(string),
		},
	}

	return token, err
}
