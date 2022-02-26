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
	accessTokenString        = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1pa2UifQ.juIlHgvyROIutaxRV95WBIr9Cn4snjqWMHO385Io_uA"
	idTokenString            = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1pa2VAZW1haWwuY29tIn0.AtNfbHDL86yYD3MulZO5jzAsZRRSmgMjNmY5EmSLFO4"
	refreshTokenString       = ""
	sampleInitiateAuthOutput = cognito.InitiateAuthOutput{
		AuthenticationResult: &cognito.AuthenticationResultType{
			AccessToken:  &accessTokenString,
			IdToken:      &idTokenString,
			RefreshToken: &refreshTokenString,
		},
	}
	sampleToken = authentication.UserToken{
		AccessTokenString: *sampleInitiateAuthOutput.AuthenticationResult.AccessToken,
		Username:          "mike",
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
			assert.Equal(t, *sampleInitiateAuthOutput.AuthenticationResult.AccessToken, token.AccessTokenString, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}
