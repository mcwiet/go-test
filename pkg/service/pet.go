package service

import (
	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/model"
)

type PetDao interface {
	Delete(id string) error
	GetById(id string) (model.Pet, error)
	Insert(model.Pet) error
	GetTotalCount() (int, error)
	Query(count int, exclusiveStartId string) ([]model.Pet, bool, error)
}

type CursorEncoder interface {
	Encode(string) string
	Decode(string) (string, error)
}

// Object containing data needed to use the Pet service
type PetService struct {
	petDao  PetDao
	encoder CursorEncoder
}

// Creates a Pet service object
func NewPetService(petDao PetDao, encoder CursorEncoder) PetService {
	return PetService{
		petDao:  petDao,
		encoder: encoder,
	}
}

// Create a new pet
func (s *PetService) Create(name string, age int) (model.Pet, error) {
	pet := model.Pet{
		Id:   uuid.NewString(),
		Name: name,
		Age:  age,
	}
	err := s.petDao.Insert(pet)
	return pet, err
}

// Deletes a pet
func (s *PetService) Delete(id string) error {
	err := s.petDao.Delete(id)
	return err
}

// Gets a single pet
func (s *PetService) GetById(id string) (model.Pet, error) {
	pet, err := s.petDao.GetById(id)
	return pet, err
}

// Lists pets
func (s *PetService) List(first int, after string) (model.PetConnection, error) {
	exclusiveStartId, _ := s.encoder.Decode(after)
	pets, hasNextPage, err := s.petDao.Query(first, exclusiveStartId)
	if err != nil {
		return model.PetConnection{}, err
	}

	totalCount, err := s.petDao.GetTotalCount()
	if err != nil {
		return model.PetConnection{}, err
	}

	endCursor := ""
	if len(pets) > 0 {
		lastId := pets[len(pets)-1].Id
		endCursor = s.encoder.Encode(lastId)
	}

	connection := model.PetConnection{
		TotalCount: totalCount,
		Edges:      []model.PetEdge{},
		PageInfo: model.PageInfo{
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		},
	}
	for _, pet := range pets {
		connection.Edges = append(connection.Edges, model.PetEdge{
			Node:   pet,
			Cursor: s.encoder.Encode(pet.Id),
		})
	}

	return connection, err
}
