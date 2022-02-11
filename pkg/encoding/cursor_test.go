package encoding_test

import (
	"testing"

	"github.com/mcwiet/go-test/pkg/encoding"
	"github.com/stretchr/testify/assert"
)

func TestCursorEncoder(t *testing.T) {
	// Setup
	encoder := encoding.NewCursorEncoder()
	original := "the original string"

	// Execute
	encoded := encoder.Encode(original)
	decoded, err := encoder.Decode(encoded)

	// Verify
	assert.Nil(t, err, "no error")
	assert.Equal(t, original, decoded, "decoded matches original")
	assert.NotEqual(t, original, encoded, "encoded does not match original")
}
