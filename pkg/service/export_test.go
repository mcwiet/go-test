package service_test

import (
	"encoding/base64"

	"github.com/mcwiet/go-test/pkg/model"
)

var (
	SampleEncoder = FakeEncoder{}
)

type FakeEncoder struct{}

func (e *FakeEncoder) Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}
func (e *FakeEncoder) Decode(input string) (string, error) {
	valBytes, _ := base64.StdEncoding.DecodeString(input)
	return string(valBytes), nil
}

type FakePetDao struct {
	deleteErr          error
	getByIdPet         model.Pet
	getByIdErr         error
	getTotalCountValue int
	getTotalCountErr   error
	insertErr          error
	queryPets          []model.Pet
	queryHasNextPage   bool
	queryErr           error
	updateErr          error
}

func (f FakePetDao) Delete(string) error {
	return f.deleteErr
}
func (f FakePetDao) GetById(string) (model.Pet, error) {
	return f.getByIdPet, f.getByIdErr
}
func (f FakePetDao) GetTotalCount() (int, error) {
	return f.getTotalCountValue, f.getTotalCountErr
}
func (f FakePetDao) Insert(model.Pet) error {
	return f.insertErr
}
func (f FakePetDao) Query(count int, exclusiveStartId string) ([]model.Pet, bool, error) {
	return f.queryPets, f.queryHasNextPage, f.queryErr
}
func (f FakePetDao) Update(pet model.Pet) error {
	return f.updateErr
}

type FakeUserDao struct {
	GetByUsernameOutput model.User
	GetByUsernameErr    error
	GetTotalCountOutput int
	GetTotalCountErr    error
	ListOutput          []model.User
	ListToken           string
	ListErr             error
}

func (u *FakeUserDao) GetByUsername(string) (model.User, error) {
	return u.GetByUsernameOutput, u.GetByUsernameErr
}
func (u *FakeUserDao) GetTotalCount() (int, error) {
	return u.GetTotalCountOutput, u.GetTotalCountErr
}
func (u *FakeUserDao) List(int, string) ([]model.User, string, error) {
	return u.ListOutput, u.ListToken, u.ListErr
}
