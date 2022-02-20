package controller_test

import "github.com/mcwiet/go-test/pkg/model"

type FakePetService struct {
	createPet      model.Pet
	createErr      error
	deleteErr      error
	getByIdUser    model.Pet
	getByIdErr     error
	listConnection model.PetConnection
	listErr        error
	updateOwnerPet model.Pet
	updateOwnerErr error
}

func (s *FakePetService) Create(name string, age int, owner string) (model.Pet, error) {
	return s.createPet, s.createErr
}
func (s *FakePetService) Delete(id string) error {
	return s.deleteErr
}
func (s *FakePetService) GetById(id string) (model.Pet, error) {
	return s.getByIdUser, s.getByIdErr
}
func (s *FakePetService) List(first int, after string) (model.PetConnection, error) {
	return s.listConnection, s.listErr
}
func (s *FakePetService) UpdateOwner(id string, owner string) (model.Pet, error) {
	return s.updateOwnerPet, s.updateOwnerErr
}

type FakeUserService struct {
	getByUsernameUser model.User
	getByUsernameErr  error
	listConnection    model.UserConnection
	listErr           error
}

func (s *FakeUserService) GetByUsername(username string) (model.User, error) {
	return s.getByUsernameUser, s.getByUsernameErr
}
func (s *FakeUserService) List(first int, after string) (model.UserConnection, error) {
	return s.listConnection, s.listErr
}
