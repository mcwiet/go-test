package authentication_test

import (
	"errors"
	"testing"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/jsii-runtime-go"
	"github.com/mcwiet/go-test/pkg/authentication"
	"github.com/stretchr/testify/assert"
)

// Define mocks / stubs
type fakeProvider struct {
	returnedValue interface{}
	returnedErr   error
}

// Define mock / stub behavior
func (p *fakeProvider) InitiateAuth(input *cognito.InitiateAuthInput) (*cognito.InitiateAuthOutput, error) {
	ret, _ := p.returnedValue.(*cognito.InitiateAuthOutput)
	return ret, p.returnedErr
}

// Define common data
var (
	sampleInitiateAuthOutput = cognito.InitiateAuthOutput{
		AuthenticationResult: &cognito.AuthenticationResultType{
			AccessToken: jsii.String("access token"),
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
			provider:      fakeProvider{returnedValue: &sampleInitiateAuthOutput},
			expectedToken: *sampleInitiateAuthOutput.AuthenticationResult.AccessToken,
			expectErr:     false,
		},
		{
			name:      "provider error",
			provider:  fakeProvider{returnedErr: errors.New("auth error")},
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
			assert.Equal(t, *sampleInitiateAuthOutput.AuthenticationResult.AccessToken, token, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}
