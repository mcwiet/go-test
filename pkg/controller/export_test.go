package controller_test

import "github.com/mcwiet/go-test/pkg/model"

type FakePetService struct {
	createOutput      model.Pet
	createErr         error
	deleteErr         error
	getByIdOutput     model.Pet
	getByIdErr        error
	listOutput        model.PetConnection
	listErr           error
	updateOwnerOutput model.Pet
	updateOwnerErr    error
}

func (s *FakePetService) Create(name string, age int, owner string) (model.Pet, error) {
	return s.createOutput, s.createErr
}
func (s *FakePetService) Delete(id string) error {
	return s.deleteErr
}
func (s *FakePetService) GetById(id string) (model.Pet, error) {
	return s.getByIdOutput, s.getByIdErr
}
func (s *FakePetService) List(first int, after string) (model.PetConnection, error) {
	return s.listOutput, s.listErr
}
func (s *FakePetService) UpdateOwner(id string, owner string) (model.Pet, error) {
	return s.updateOwnerOutput, s.updateOwnerErr
}
