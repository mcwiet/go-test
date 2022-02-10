package controller

import (
	"github.com/mcwiet/go-test/pkg/model"
)

type PersonService interface {
	Create(name string, age int) (*model.Person, error)
	Delete(id string) error
	GetById(id string) (*model.Person, error)
	List(first int, after string) (model.PersonConnection, error)
}

// Object containing data needed for the Person controller
type PersonController struct {
	personService PersonService
}

// Creates a new person controller object
func NewPersonController(service PersonService) PersonController {
	return PersonController{
		personService: service,
	}
}

// Handles request for creating a person
func (c *PersonController) HandleCreate(request Request) Response {
	name := request.Arguments["name"].(string)
	age := int(request.Arguments["age"].(float64))
	person, err := c.personService.Create(name, age)
	if err == nil {
		return Response{Data: *person}
	} else {
		return Response{Error: err}
	}
}

// Handles request for deleting a person
func (c *PersonController) HandleDelete(request Request) Response {
	id := request.Arguments["id"].(string)
	err := c.personService.Delete(id)
	if err == nil {
		return Response{}
	} else {
		return Response{Error: err}
	}
}

// Handles request for getting a specific person
func (c *PersonController) HandleGet(request Request) Response {
	id := request.Arguments["id"].(string)
	person, err := c.personService.GetById(id)
	if err == nil {
		return Response{Data: *person}
	} else {
		return Response{Error: err}
	}
}

// Handles request for listing people
func (c *PersonController) HandleList(request Request) Response {
	first := 0
	if request.Arguments["first"] != nil {
		first = int(request.Arguments["first"].(float64))
	}
	after := ""
	if request.Arguments["after"] != nil {
		after = request.Arguments["after"].(string)
	}
	connection, err := c.personService.List(first, after)
	if err == nil {
		return Response{Data: connection}
	} else {
		return Response{Error: err}
	}
}
