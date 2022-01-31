package service

import (
	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/data"
	"github.com/mcwiet/go-test/pkg/model"
)

// Object containing data needed to use the Person service
type PersonService struct {
	personDao data.PersonDao
}

// Creates a Person service object
func NewPersonService(personDao data.PersonDao) PersonService {
	return PersonService{
		personDao: personDao,
	}
}

// Create a new person
func (s *PersonService) CreatePerson(name string, age int) (*model.Person, error) {
	person := model.Person{
		Id:   uuid.NewString(),
		Name: name,
		Age:  age,
	}
	err := s.personDao.AddPerson(&person)
	return &person, err
}

// Deletes a person
func (s *PersonService) DeletePerson(id string) error {
	err := s.personDao.DeletePerson(id)
	return err
}

// Gets a single person
func (s *PersonService) GetPerson(id string) (*model.Person, error) {
	person, err := s.personDao.GetPerson(id)
	return person, err
}

// Gets a list of people
func (s *PersonService) GetPeople() (*[]model.Person, error) {
	people, err := s.personDao.GetPeople()
	return people, err
}
