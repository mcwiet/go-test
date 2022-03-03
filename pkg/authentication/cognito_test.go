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
	SampleAccessTokenString  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1pa2UifQ.juIlHgvyROIutaxRV95WBIr9Cn4snjqWMHO385Io_uA"
	SampleIdTokenString      = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb2duaXRvOnVzZXJuYW1lIjoibWlrZSIsImNvZ25pdG86Z3JvdXBzIjpbImdyb3VwIl0sImVtYWlsIjoibWlrZUBlbWFpbC5jb20ifQ.A9WKwOIkoRuOYYMaY05uBYyNwHJ7mNONZXQCyqkoPpI"
	SampleRefreshTokenString = "refresh"
	SampleUsername           = "mike"            // same as in ID token string
	SampleEmail              = "mike@email.com"  // same as in ID token string
	SampleGroups             = []string{"group"} // same as in ID token string
	SampleInitiateAuthOutput = cognito.InitiateAuthOutput{
		AuthenticationResult: &cognito.AuthenticationResultType{
			AccessToken:  &SampleAccessTokenString,
			IdToken:      &SampleIdTokenString,
			RefreshToken: &SampleRefreshTokenString,
		},
	}
	SampleToken = authentication.UserToken{
		AccessTokenString:  *SampleInitiateAuthOutput.AuthenticationResult.AccessToken,
		IdTokenString:      *SampleInitiateAuthOutput.AuthenticationResult.IdToken,
		RefreshTokenString: *SampleInitiateAuthOutput.AuthenticationResult.RefreshToken,
		Username:           SampleUsername,
		Email:              SampleEmail,
		Groups:             SampleGroups,
	}
)

func TestLogin(t *testing.T) {
	// Define test struct
	type Test struct {
		name          string
		provider      fakeProvider
		expectedToken authentication.UserToken
		expectErr     bool
	}

	// Define tests
	tests := []Test{
		{
			name:          "valid input",
			provider:      fakeProvider{intiateAuthOutput: SampleInitiateAuthOutput},
			expectedToken: SampleToken,
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
			assert.Equal(t, *SampleInitiateAuthOutput.AuthenticationResult.AccessToken, token.AccessTokenString, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}
