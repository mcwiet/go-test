package service

import (
	"github.com/mcwiet/go-test/pkg/data"
	"github.com/mcwiet/go-test/pkg/model"
)

type PersonService struct {
	personDao data.PersonDao
}

func NewPersonService(personDao data.PersonDao) PersonService {
	return PersonService{
		personDao: personDao,
	}
}

func (s PersonService) GetPerson(id string) (*model.Person, error) {
	person, err := s.personDao.GetPerson(id)
	return person, err
}

func (s PersonService) GetPeople() (*[]model.Person, error) {
	people, err := s.personDao.GetPeople()
	return people, err
}
