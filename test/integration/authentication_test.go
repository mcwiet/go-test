package integration_test

import (
	"testing"

	"github.com/mcwiet/go-test/test/integration"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	// Setup
	email := integration.GetRequiredEnv("TEST_USER_EMAIL")
	password := integration.GetRequiredEnv("TEST_USER_PASSWORD")

	// Test
	token, err := Authenticator.Login(email, password)

	// Verify
	assert.Nil(t, err)
	assert.NotNil(t, token)
}
