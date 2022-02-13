package controller

import (
	"encoding/json"

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
	var input model.CreatePetInput
	inputBytes, _ := json.Marshal(request.Arguments["input"])
	json.Unmarshal(inputBytes, &input)

	pet, err := c.petService.Create(input.Name, input.Age, input.Owner)

	if err == nil {
		return Response{Data: model.CreatePetPayload{Pet: pet}}
	} else {
		return Response{Error: err}
	}
}

// Handles request for deleting a pet
func (c *PetController) HandleDelete(request Request) Response {
	var input model.DeletePetInput
	inputBytes, _ := json.Marshal(request.Arguments["input"])
	json.Unmarshal(inputBytes, &input)

	err := c.petService.Delete(input.Id)

	if err == nil {
		//lint:ignore S1016 Input and payload happen to look similar
		return Response{Data: model.DeletePetPayload{Id: input.Id}}
	} else {
		return Response{Error: err}
	}
}

// Handles request for getting a specific pet
func (c *PetController) HandleGet(request Request) Response {
	var input model.PetInput
	inputBytes, _ := json.Marshal(request.Arguments["input"])
	json.Unmarshal(inputBytes, &input)

	pet, err := c.petService.GetById(input.Id)

	if err == nil {
		return Response{Data: pet}
	} else {
		return Response{Error: err}
	}
}

// Handles request for listing pets
func (c *PetController) HandleList(request Request) Response {
	var input model.PetsInput
	inputBytes, _ := json.Marshal(request.Arguments["input"])
	json.Unmarshal(inputBytes, &input)

	connection, err := c.petService.List(input.First, input.After)

	if err == nil {
		return Response{Data: connection}
	} else {
		return Response{Error: err}
	}
}
