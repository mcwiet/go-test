package controller

import (
	"github.com/mcwiet/go-test/pkg/service"
)

type PersonController struct {
	personService service.PersonService
}

func NewPerosnController(service service.PersonService) PersonController {
	return PersonController{
		personService: service,
	}
}

func (c PersonController) GetPerson(request Request) Response {
	person, err := c.personService.GetPerson(request.Arguments["id"])
	if err == nil {
		return Response{
			Data: *person,
		}
	} else {
		return Response{
			Error: err,
		}
	}
}

func (c PersonController) GetPeople(request Request) Response {
	people, err := c.personService.GetPeople()
	if err == nil {
		return Response{
			Data: *people,
		}
	} else {
		return Response{
			Error: err,
		}
	}
}
