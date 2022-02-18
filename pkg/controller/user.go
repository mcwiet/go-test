package controller

import (
	"encoding/json"
	"log"

	"github.com/mcwiet/go-test/pkg/model"
)

type UserService interface {
	GetByUsername(username string) (model.User, error)
	List(first int, after string) (model.UserConnection, error)
}

// Object containing data needed for the User controller
type UserController struct {
	userService UserService
}

// Creates a new user controller object
func NewUserController(service UserService) UserController {
	return UserController{
		userService: service,
	}
}

// Handles request for getting a specific user
func (c *UserController) HandleGet(request Request) Response {
	var input model.UserInput
	inputBytes, _ := json.Marshal(request.Arguments["input"])
	json.Unmarshal(inputBytes, &input)

	user, err := c.userService.GetByUsername(input.Username)

	if err == nil {
		return Response{Data: user}
	} else {
		return Response{Error: err}
	}
}

// Handles request for listing users
func (c *UserController) HandleList(request Request) Response {
	var input model.UsersInput
	inputBytes, _ := json.Marshal(request.Arguments["input"])
	json.Unmarshal(inputBytes, &input)

	connection, err := c.userService.List(input.First, input.After)

	log.Println(connection)

	if err == nil {
		return Response{Data: connection}
	} else {
		return Response{Error: err}
	}
}
