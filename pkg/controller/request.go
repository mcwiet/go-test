package controller

import (
	"github.com/mcwiet/go-test/pkg/model"
)

// Standard request format
type Request struct {
	Arguments      map[string]interface{}
	FieldName      string
	ParentTypeName string
	Identity       model.Identity
}
