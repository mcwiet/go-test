package controller

import (
	"log"

	"github.com/mcwiet/go-test/pkg/service"
)

// Object containing data needed for the Person controller
type PersonController struct {
	personService *service.PersonService
}

// Creates a new Person controller object
func NewPersonController(service *service.PersonService) PersonController {
	return PersonController{
		personService: service,
	}
}

// Handles request for creating a person
func (c *PersonController) HandleCreate(request Request) Response {
	log.Println(request.Arguments["age"])
	person, err := c.personService.Create(
		request.Arguments["name"].(string),
		int(request.Arguments["age"].(float64)))
	if err == nil {
		return Response{Data: *person}
	} else {
		return Response{Error: err}
	}
}

// Handles request for deleting a person
func (c *PersonController) HandleDelete(request Request) Response {
	err := c.personService.Delete(request.Arguments["id"].(string))
	if err == nil {
		return Response{Data: true}
	} else {
		return Response{Error: err}
	}
}

// Handles request for getting a specific person
func (c *PersonController) HandleGet(request Request) Response {
	person, err := c.personService.GetById(request.Arguments["id"].(string))
	if err == nil {
		return Response{Data: *person}
	} else {
		return Response{Error: err}
	}
}

// Handles request for listing people
func (c *PersonController) HandleList(request Request) Response {
	people, err := c.personService.List()
	if err == nil {
		return Response{Data: *people}
	} else {
		return Response{Error: err}
	}
}
