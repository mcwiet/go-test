package main

import (
	"encoding/json"

	"github.com/mcwiet/go-test/pkg/controller"
	"github.com/mcwiet/go-test/pkg/model"
)

type AppSyncRequest struct {
	Arguments map[string]interface{} `json:"arguments"`
	Info      struct {
		FieldName      string `json:"fieldName"`
		ParentTypeName string `json:"parentTypeName"`
	}
	Identity struct {
		Claims struct {
			Username string   `json:"cognito:username"`
			Email    string   `json:"email"`
			Groups   []string `json:"cognito:groups"`
		}
	}
}

// Takes an arbitrary object (whose shape should match an AppSyncRequest) and converts it into a standardized request
func NewRequest(req interface{}) controller.Request {
	appsync := newAppSyncRequest(req)
	groups := convertToSet(appsync.Identity.Claims.Groups)
	return controller.Request{
		Arguments:      appsync.Arguments,
		FieldName:      appsync.Info.FieldName,
		ParentTypeName: appsync.Info.ParentTypeName,
		Identity: model.Identity{
			Username: appsync.Identity.Claims.Username,
			Email:    appsync.Identity.Claims.Email,
			Groups:   groups,
		},
	}
}

func newAppSyncRequest(req interface{}) AppSyncRequest {
	var appsync AppSyncRequest
	reqBytes, _ := json.Marshal(req)
	json.Unmarshal(reqBytes, &appsync)
	return appsync
}

func convertToSet(arr []string) map[string]bool {
	set := map[string]bool{}
	for _, item := range arr {
		set[item] = true
	}
	return set
}
