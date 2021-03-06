package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/model"
)

type PetDao interface {
	Delete(id string) error
	GetById(id string) (model.Pet, error)
	GetTotalCount() (int, error)
	Insert(model.Pet) error
	Query(count int, exclusiveStartId string) ([]model.Pet, bool, error)
	Update(model.Pet) error
}

type Authorizer interface {
	IsAuthorized(model.Identity, model.Pet, PetAction) bool
}

type CursorEncoder interface {
	Encode(string) string
	Decode(string) (string, error)
}

// Object containing data needed to use the Pet service
type PetService struct {
	authorizer Authorizer
	petDao     PetDao
	userDao    UserDao
	encoder    CursorEncoder
}

// Permissible pet actions
type PetAction int

const (
	PetActionUndefined PetAction = iota
	PetActionUpdateOwner
)

// Creates a Pet service object
func NewPetService(petDao PetDao, userDao UserDao, authorizer Authorizer, encoder CursorEncoder) PetService {
	return PetService{
		authorizer: authorizer,
		petDao:     petDao,
		userDao:    userDao,
		encoder:    encoder,
	}
}

// Create a new pet
func (s *PetService) Create(name string, age int, owner string) (model.Pet, error) {
	pet := model.Pet{
		Id:    uuid.NewString(),
		Name:  name,
		Age:   age,
		Owner: owner,
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
	exclusiveStartId, err := s.encoder.Decode(after)
	if err != nil {
		return model.PetConnection{}, err
	}

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

// Updates the owner of a pet
func (s *PetService) UpdateOwner(requestor model.Identity, id string, owner string) (model.Pet, error) {
	pet, err := s.petDao.GetById(id)
	if err != nil {
		return model.Pet{}, errors.New("could not find pet ID " + id)
	}

	authorized := s.authorizer.IsAuthorized(requestor, pet, PetActionUpdateOwner)
	if !authorized {
		return model.Pet{}, errors.New("not authorized to update the owner on this pet")
	}

	if owner != "" {
		_, err = s.userDao.GetByUsername(owner)
		if err != nil {
			return model.Pet{}, errors.New(owner + " is not a valid user")
		}
	}

	pet.Owner = owner
	err = s.petDao.Update(pet)

	return pet, err
}
