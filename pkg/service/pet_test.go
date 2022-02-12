package service_test

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/mcwiet/go-test/pkg/service"
	"github.com/stretchr/testify/assert"
)

// Define mocks / stubs
type fakePetDao struct {
	deleteErr          error
	getByIdPet         model.Pet
	getByIdErr         error
	getTotalCountValue int
	getTotalCountErr   error
	insertErr          error
	queryPets          []model.Pet
	queryHasNextPage   bool
	queryErr           error
}
type petDaoGetTotalCount = service.PetDao
type fakeEncoder struct{}

// Define mock / stub behavior
func (f fakePetDao) Delete(string) error {
	return f.deleteErr
}
func (f fakePetDao) GetById(string) (model.Pet, error) {
	return f.getByIdPet, f.getByIdErr
}
func (f fakePetDao) GetTotalCount() (int, error) {
	return f.getTotalCountValue, f.getTotalCountErr
}
func (f fakePetDao) Insert(model.Pet) error {
	return f.insertErr
}
func (f fakePetDao) Query(count int, exclusiveStartId string) ([]model.Pet, bool, error) {
	return f.queryPets, f.queryHasNextPage, f.queryErr
}

func (e *fakeEncoder) Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}
func (e *fakeEncoder) Decode(input string) (string, error) {
	valBytes, _ := base64.StdEncoding.DecodeString(input)
	return string(valBytes), nil
}

// Define common data
var (
	encoder    = fakeEncoder{}
	samplePet1 = model.Pet{
		Id:    uuid.NewString(),
		Name:  "pet 1",
		Age:   12,
		Owner: "User 1",
	}
	samplePet2 = model.Pet{
		Id:    uuid.NewString(),
		Name:  "pet 2",
		Age:   20,
		Owner: "User 2",
	}
)

func TestCreate(t *testing.T) {
	// Define test struct
	type Test struct {
		name      string
		petDao    fakePetDao
		petName   string
		petAge    int
		petOwner  string
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid create",
			petDao:    fakePetDao{},
			petName:   samplePet1.Name,
			petAge:    samplePet1.Age,
			petOwner:  samplePet1.Owner,
			expectErr: false,
		},
		{
			name:      "DAO insert error",
			petDao:    fakePetDao{insertErr: errors.New("dao error")},
			petName:   samplePet1.Name,
			petAge:    samplePet1.Age,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		service := service.NewPetService(test.petDao, &encoder)

		// Execute
		pet, err := service.Create(test.petName, test.petAge, test.petOwner)

		// Verify
		if !test.expectErr {
			_, uuidErr := uuid.Parse(pet.Id)
			assert.Equal(t, pet.Name, test.petName, test.name)
			assert.Equal(t, pet.Age, test.petAge, test.name)
			assert.Nil(t, uuidErr, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestGetById(t *testing.T) {
	// Define test struct
	type Test struct {
		name        string
		petDao      fakePetDao
		petId       string
		expectedPet model.Pet
		expectErr   bool
	}

	// Define tests
	tests := []Test{
		{
			name:        "valid get by id",
			petDao:      fakePetDao{getByIdPet: samplePet1},
			petId:       samplePet1.Id,
			expectedPet: samplePet1,
			expectErr:   false,
		},
		{
			name:        "DAO get error",
			petDao:      fakePetDao{getByIdErr: errors.New("dao error")},
			petId:       samplePet1.Id,
			expectedPet: samplePet1,
			expectErr:   true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		service := service.NewPetService(test.petDao, &encoder)

		// Execute
		pet, err := service.GetById(test.petId)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedPet, pet, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestDelete(t *testing.T) {
	// Define test struct
	type Test struct {
		name      string
		petDao    fakePetDao
		petId     string
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid delete",
			petDao:    fakePetDao{deleteErr: nil},
			petId:     samplePet1.Id,
			expectErr: false,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		service := service.NewPetService(test.petDao, &encoder)

		// Execute
		err := service.Delete(test.petId)

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestList(t *testing.T) {
	// Define test struct
	type Test struct {
		name               string
		petDao             fakePetDao
		first              int
		after              string
		expectedConnection model.PetConnection
		expectErr          bool
	}

	// Define tests
	tests := []Test{
		{
			name: "list all pets",
			petDao: fakePetDao{
				getTotalCountValue: 2,
				queryPets:          []model.Pet{samplePet1, samplePet2},
				queryHasNextPage:   false,
			},
			first: 10,
			after: "",
			expectedConnection: model.PetConnection{
				TotalCount: 2,
				Edges: []model.PetEdge{
					{
						Node:   samplePet1,
						Cursor: encoder.Encode(samplePet1.Id),
					},
					{
						Node:   samplePet2,
						Cursor: encoder.Encode(samplePet2.Id),
					},
				},
				PageInfo: model.PageInfo{
					EndCursor:   encoder.Encode(samplePet2.Id),
					HasNextPage: false,
				},
			},
			expectErr: false,
		},
		{
			name: "list first of two pets",
			petDao: fakePetDao{
				getTotalCountValue: 2,
				queryPets:          []model.Pet{samplePet1},
				queryHasNextPage:   true,
			},
			first: 1,
			after: "",
			expectedConnection: model.PetConnection{
				TotalCount: 2,
				Edges: []model.PetEdge{
					{
						Node:   samplePet1,
						Cursor: encoder.Encode(samplePet1.Id),
					},
				},
				PageInfo: model.PageInfo{
					EndCursor:   encoder.Encode(samplePet1.Id),
					HasNextPage: true,
				},
			},
			expectErr: false,
		},
		{
			name: "list second of two pets",
			petDao: fakePetDao{
				getTotalCountValue: 2,
				queryPets:          []model.Pet{samplePet2},
				queryHasNextPage:   false,
			},
			first: 1,
			after: encoder.Encode(samplePet1.Id),
			expectedConnection: model.PetConnection{
				TotalCount: 2,
				Edges: []model.PetEdge{
					{
						Node:   samplePet2,
						Cursor: encoder.Encode(samplePet2.Id),
					},
				},
				PageInfo: model.PageInfo{
					EndCursor:   encoder.Encode(samplePet2.Id),
					HasNextPage: false,
				},
			},
			expectErr: false,
		},
		{
			name: "DAO total count error",
			petDao: fakePetDao{
				getTotalCountErr: assert.AnError,
			},
			first:     1,
			after:     "",
			expectErr: true,
		},
		{
			name: "DAO query error",
			petDao: fakePetDao{
				getTotalCountValue: 2,
				queryErr:           assert.AnError,
			},
			first:     1,
			after:     "",
			expectErr: true,
		},
	}

	//Run tests
	for _, test := range tests {
		// Setup
		service := service.NewPetService(test.petDao, &encoder)

		// Execute
		pets, err := service.List(test.first, test.after)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedConnection, pets, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}
