package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	// Setup
	email := GetRequiredEnv("TEST_USER_EMAIL")
	password := GetRequiredEnv("TEST_USER_PASSWORD")

	// Test
	token, err := Authenticator.Login(email, password)

	// Verify
	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, email, token.Email, "email should match")
}
