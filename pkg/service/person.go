package service

import (
	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/model"
)

type PersonDao interface {
	Delete(id string) error
	GetById(id string) (model.Person, error)
	Insert(model.Person) error
	GetTotalCount() (int, error)
	Query(count int, exclusiveStartId string) ([]model.Person, bool, error)
}

type CursorEncoder interface {
	Encode(string) string
	Decode(string) (string, error)
}

// Object containing data needed to use the Person service
type PersonService struct {
	personDao PersonDao
	encoder   CursorEncoder
}

// Creates a Person service object
func NewPersonService(personDao PersonDao, encoder CursorEncoder) PersonService {
	return PersonService{
		personDao: personDao,
		encoder:   encoder,
	}
}

// Create a new person
func (s *PersonService) Create(name string, age int) (model.Person, error) {
	person := model.Person{
		Id:   uuid.NewString(),
		Name: name,
		Age:  age,
	}
	err := s.personDao.Insert(person)
	return person, err
}

// Deletes a person
func (s *PersonService) Delete(id string) error {
	err := s.personDao.Delete(id)
	return err
}

// Gets a single person
func (s *PersonService) GetById(id string) (model.Person, error) {
	person, err := s.personDao.GetById(id)
	return person, err
}

// Lists people
func (s *PersonService) List(first int, after string) (model.PersonConnection, error) {
	exclusiveStartId, _ := s.encoder.Decode(after)
	people, hasNextPage, err := s.personDao.Query(first, exclusiveStartId)
	if err != nil {
		return model.PersonConnection{}, err
	}

	totalCount, err := s.personDao.GetTotalCount()
	if err != nil {
		return model.PersonConnection{}, err
	}

	endCursor := ""
	if len(people) > 0 {
		lastId := people[len(people)-1].Id
		endCursor = s.encoder.Encode(lastId)
	}

	connection := model.PersonConnection{
		TotalCount: totalCount,
		Edges:      []model.PersonEdge{},
		PageInfo: model.PageInfo{
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		},
	}
	for _, person := range people {
		connection.Edges = append(connection.Edges, model.PersonEdge{
			Node:   person,
			Cursor: s.encoder.Encode(person.Id),
		})
	}

	return connection, err
}
