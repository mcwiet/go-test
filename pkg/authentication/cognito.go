package authentication

import (
	"errors"

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
	accessClaims, err := getClaims(accessToken)
	if err != nil {
		return UserToken{}, errors.New("could not parse user's access token")
	}
	idClaims, err := getClaims(idToken)
	if err != nil {
		return UserToken{}, errors.New("could not parse user's ID token")
	}

	token := UserToken{
		AccessTokenString:  accessToken,
		IdTokenString:      idToken,
		RefreshTokenString: refreshToken,
		Username:           accessClaims["username"].(string),
		Email:              idClaims["email"].(string),
	}

	return token, err
}

func getClaims(token string) (jwt.MapClaims, error) {
	decoded, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return jwt.MapClaims{}, err
	}

	claims, ok := decoded.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}, err
	}

	return claims, nil
}
