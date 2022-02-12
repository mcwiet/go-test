package controller

import (
	"github.com/mcwiet/go-test/pkg/model"
)

type PetService interface {
	Create(name string, age int, owner string) (model.Pet, error)
	Delete(id string) error
	GetById(id string) (model.Pet, error)
	List(first int, after string) (model.PetConnection, error)
}

// Object containing data needed for the Pet controller
type PetController struct {
	petService PetService
}

// Creates a new pet controller object
func NewPetController(service PetService) PetController {
	return PetController{
		petService: service,
	}
}

// Handles request for creating a pet
func (c *PetController) HandleCreate(request Request) Response {
	name := request.Arguments["name"].(string)
	age := int(request.Arguments["age"].(float64))
	owner := ""
	if request.Arguments["owner"] != nil {
		owner = request.Arguments["owner"].(string)
	}
	pet, err := c.petService.Create(name, age, owner)
	if err == nil {
		return Response{Data: pet}
	} else {
		return Response{Error: err}
	}
}

// Handles request for deleting a pet
func (c *PetController) HandleDelete(request Request) Response {
	id := request.Arguments["id"].(string)
	err := c.petService.Delete(id)
	if err == nil {
		return Response{}
	} else {
		return Response{Error: err}
	}
}

// Handles request for getting a specific pet
func (c *PetController) HandleGet(request Request) Response {
	id := request.Arguments["id"].(string)
	pet, err := c.petService.GetById(id)
	if err == nil {
		return Response{Data: pet}
	} else {
		return Response{Error: err}
	}
}

// Handles request for listing pets
func (c *PetController) HandleList(request Request) Response {
	first := 0
	if request.Arguments["first"] != nil {
		first = int(request.Arguments["first"].(float64))
	}
	after := ""
	if request.Arguments["after"] != nil {
		after = request.Arguments["after"].(string)
	}
	connection, err := c.petService.List(first, after)
	if err == nil {
		return Response{Data: connection}
	} else {
		return Response{Error: err}
	}
}
