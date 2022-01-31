package controller

import (
	"log"

	"github.com/mcwiet/go-test/pkg/service"
)

// Object containing data needed for the Person controller
type PersonController struct {
	personService service.PersonService
}

// Creates a new Person controller object
func NewPerosnController(service service.PersonService) PersonController {
	return PersonController{
		personService: service,
	}
}

// Handles request for creating a person
func (c PersonController) HandleCreatePerson(request Request) Response {
	log.Println(request.Arguments["age"])
	person, err := c.personService.CreatePerson(
		request.Arguments["name"].(string),
		int(request.Arguments["age"].(float64)))
	if err == nil {
		return Response{Data: *person}
	} else {
		return Response{Error: err}
	}
}

// Handles request for deleting a person
func (c PersonController) HandleDeletePerson(request Request) Response {
	err := c.personService.DeletePerson(request.Arguments["id"].(string))
	if err == nil {
		return Response{Data: true}
	} else {
		return Response{Error: err}
	}
}

// Handles request for getting a specific person
func (c PersonController) HandleGetPerson(request Request) Response {
	person, err := c.personService.GetPerson(request.Arguments["id"].(string))
	if err == nil {
		return Response{Data: *person}
	} else {
		return Response{Error: err}
	}
}

// Handles request for getting a list of people
func (c PersonController) HandleGetPeople(request Request) Response {
	people, err := c.personService.GetPeople()
	if err == nil {
		return Response{Data: *people}
	} else {
		return Response{Error: err}
	}
}
