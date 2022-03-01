package controller

import (
	"encoding/json"
)

type Request struct {
	Arguments map[string]interface{} `json:"arguments"`
	Info      struct {
		FieldName      string `json:"fieldName"`
		ParentTypeName string `json:"parentTypeName"`
	}
	Identity struct {
		Claims struct {
			Username string `json:"cognito:username"`
			Email    string `json:"email"`
		}
	}
}

// Creates a request object given an arbitrary JSON input object
func NewRequest(request interface{}) Request {
	var convertedRequest Request
	requestBytes, _ := json.Marshal(request)
	json.Unmarshal(requestBytes, &convertedRequest)
	return convertedRequest
}
