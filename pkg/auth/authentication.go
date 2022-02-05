package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/jsii-runtime-go"
)

type (
	App struct {
		CognitoClient   *cognito.CognitoIdentityProvider
		UserPoolID      string
		AppClientID     string
		AppClientSecret string
	}

	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	Response struct {
		Error error `json:"error"`
	}
)

func (a *App) Login(email string, password string) (string, error) {

	params := map[string]*string{
		"USERNAME": jsii.String(email),
		"PASSWORD": jsii.String(password),
	}

	secretHash := computeSecretHash(a.AppClientSecret, email, a.AppClientID)
	params["SECRET_HASH"] = jsii.String(secretHash)

	authTry := &cognito.InitiateAuthInput{
		AuthFlow: jsii.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME":    jsii.String(*params["USERNAME"]),
			"PASSWORD":    jsii.String(*params["PASSWORD"]),
			"SECRET_HASH": jsii.String(*params["SECRET_HASH"]),
		},
		ClientId: jsii.String(a.AppClientID), // this is the app client ID
	}

	authResp, err := a.CognitoClient.InitiateAuth(authTry)
	if err != nil {
		log.Println(err)
	}

	log.Println(*authResp)

	token := *authResp.AuthenticationResult.AccessToken
	return token, err
}

func computeSecretHash(clientSecret string, username string, clientId string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
