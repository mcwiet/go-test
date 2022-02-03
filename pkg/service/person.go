package service

import (
	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/model"
)

type PersonDao interface {
	Delete(string) error
	GetById(string) (*model.Person, error)
	Insert(*model.Person) error
	List() (*[]model.Person, error)
}

// Object containing data needed to use the Person service
type PersonService struct {
	personDao PersonDao
}

// Creates a Person service object
func NewPersonService(personDao PersonDao) PersonService {
	return PersonService{
		personDao: personDao,
	}
}

// Create a new person
func (s *PersonService) Create(name string, age int) (*model.Person, error) {
	person := model.Person{
		Id:   uuid.NewString(),
		Name: name,
		Age:  age,
	}
	err := s.personDao.Insert(&person)
	return &person, err
}

// Deletes a person
func (s *PersonService) Delete(id string) error {
	err := s.personDao.Delete(id)
	return err
}

// Gets a single person
func (s *PersonService) GetById(id string) (*model.Person, error) {
	person, err := s.personDao.GetById(id)
	return person, err
}

// Lists people
func (s *PersonService) List() (*[]model.Person, error) {
	people, err := s.personDao.List()
	return people, err
}
