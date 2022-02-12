package authentication_test

import (
	"errors"
	"testing"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/mcwiet/go-test/pkg/authentication"
	"github.com/stretchr/testify/assert"
)

// Define mocks / stubs
type fakeProvider struct {
	intiateAuthOutput cognito.InitiateAuthOutput
	intiateAuthErr    error
}

// Define mock / stub behavior
func (p *fakeProvider) InitiateAuth(input *cognito.InitiateAuthInput) (*cognito.InitiateAuthOutput, error) {
	return &p.intiateAuthOutput, p.intiateAuthErr
}

// Define common data
var (
	tokenString              = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJvcmlnaW5fanRpIjoianRpIiwic3ViIjoic3ViIiwiZXZlbnRfaWQiOiJpZCIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiYXdzLmNvZ25pdG8uc2lnbmluLnVzZXIuYWRtaW4iLCJhdXRoX3RpbWUiOjEsImlzcyI6Imh0dHBzOi8vY29nbml0by1pZHAudXMtZWFzdC0xLmFtYXpvbmF3cy5jb20vdXMtZWFzdC0xX1hYWFhYWCIsImV4cCI6MiwiaWF0IjoxLCJqdGkiOiJqdGkiLCJjbGllbnRfaWQiOiJpZCIsInVzZXJuYW1lIjoidXNlcm5hbWUifQ.fHDYAob2l4ibGkRcA7IVfvU6bMuNns4gPwX5oAWTmN0"
	sampleInitiateAuthOutput = cognito.InitiateAuthOutput{
		AuthenticationResult: &cognito.AuthenticationResultType{
			AccessToken: &tokenString,
		},
	}
	sampleToken = authentication.CognitoToken{
		String: *sampleInitiateAuthOutput.AuthenticationResult.AccessToken,
		Payload: authentication.CognitoTokenPayload{
			Username: "mike",
		},
	}
)

func TestLogin(t *testing.T) {
	// Define test struct
	type Test struct {
		name          string
		provider      fakeProvider
		expectedToken string
		expectErr     bool
	}

	// Define tests
	tests := []Test{
		{
			name:          "valid input",
			provider:      fakeProvider{intiateAuthOutput: sampleInitiateAuthOutput},
			expectedToken: *sampleInitiateAuthOutput.AuthenticationResult.AccessToken,
			expectErr:     false,
		},
		{
			name:      "provider error",
			provider:  fakeProvider{intiateAuthErr: errors.New("auth error")},
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		auth := authentication.NewCognitoAuthenticator(&test.provider, "app client")

		// Execute
		token, err := auth.Login("username", "password")

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
			assert.Equal(t, *sampleInitiateAuthOutput.AuthenticationResult.AccessToken, token.String, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}
