package controller_test

import (
	"testing"

	"github.com/mcwiet/go-test/pkg/controller"
	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	// Setup
	args := map[string]interface{}{
		"key1": "val1",
	}
	field := "field"
	parentType := "parentType"
	input := map[string]interface{}{
		"arguments": args,
		"info": map[string]interface{}{
			"fieldName":      field,
			"parentTypeName": parentType,
		},
	}

	// Execute
	request := controller.NewRequest(input)

	// Verify
	assert.Equal(t, args, request.Arguments, "arguments match")
	assert.Equal(t, field, request.Info.FieldName, "field name matches")
	assert.Equal(t, parentType, request.Info.ParentTypeName, "parent type matches")
}
