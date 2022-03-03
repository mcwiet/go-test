package authentication

import (
	"errors"
	"fmt"

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

type UserToken struct {
	AccessTokenString  string
	IdTokenString      string
	RefreshTokenString string
	Username           string
	Email              string
	Groups             []string
}

// Creates a new authenticator object
func NewCognitoAuthenticator(provider CognitoIdentityProvider, appClientId string) CognitoAuthenticator {
	return CognitoAuthenticator{
		provider:    provider,
		appClientId: appClientId,
	}
}

// Login to the Cognito User Pool
func (a *CognitoAuthenticator) Login(email string, password string) (UserToken, error) {
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
		return UserToken{}, err
	}

	var accessToken, idToken, refreshToken string
	if authResp != nil && authResp.AuthenticationResult != nil {
		accessToken = *authResp.AuthenticationResult.AccessToken
		idToken = *authResp.AuthenticationResult.IdToken
		refreshToken = *authResp.AuthenticationResult.RefreshToken
	}

	return buildUserToken(accessToken, idToken, refreshToken)
}

// Turn a token string into a token object (does not verify token!)
func buildUserToken(accessToken string, idToken string, refreshToken string) (UserToken, error) {
	idClaims, err := getClaims(idToken)
	if err != nil {
		return UserToken{}, errors.New("could not parse user's ID token")
	}

	groups := []string{}
	for _, group := range idClaims["cognito:groups"].([]interface{}) {
		groups = append(groups, fmt.Sprintf("%v", group))
	}

	token := UserToken{
		AccessTokenString:  accessToken,
		IdTokenString:      idToken,
		RefreshTokenString: refreshToken,
		Username:           idClaims["cognito:username"].(string),
		Email:              idClaims["email"].(string),
		Groups:             groups,
	}

	return token, err
}

func getClaims(token string) (jwt.MapClaims, error) {
	decoded, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return jwt.MapClaims{}, err
	}

	claims := decoded.Claims.(jwt.MapClaims)

	return claims, nil
}
