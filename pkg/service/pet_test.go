package service_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/mcwiet/go-test/pkg/service"
	"github.com/stretchr/testify/assert"
)

var (
	SamplePet1 = model.Pet{
		Id:    uuid.NewString(),
		Name:  "pet 1",
		Age:   12,
		Owner: "User 1",
	}
	SamplePet2 = model.Pet{
		Id:    uuid.NewString(),
		Name:  "pet 2",
		Age:   20,
		Owner: "User 2",
	}
	SamplePet1Edge = model.PetEdge{
		Node:   SamplePet1,
		Cursor: SampleEncoder.Encode(SamplePet1.Id),
	}
	SamplePet2Edge = model.PetEdge{
		Node:   SamplePet2,
		Cursor: SampleEncoder.Encode(SamplePet2.Id),
	}
)

func TestPetCreate(t *testing.T) {
	// Define test struct
	type Test struct {
		name      string
		petDao    FakePetDao
		petName   string
		petAge    int
		petOwner  string
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid create",
			petDao:    FakePetDao{},
			petName:   SamplePet1.Name,
			petAge:    SamplePet1.Age,
			petOwner:  SamplePet1.Owner,
			expectErr: false,
		},
		{
			name:      "DAO insert error",
			petDao:    FakePetDao{insertErr: errors.New("dao error")},
			petName:   SamplePet1.Name,
			petAge:    SamplePet1.Age,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		service := service.NewPetService(&test.petDao, &SampleEncoder)

		// Execute
		pet, err := service.Create(test.petName, test.petAge, test.petOwner)

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
			_, uuidErr := uuid.Parse(pet.Id)
			assert.Nil(t, uuidErr, test.name)
			assert.Equal(t, pet.Name, test.petName, test.name)
			assert.Equal(t, pet.Age, test.petAge, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestPetGetById(t *testing.T) {
	// Define test struct
	type Test struct {
		name        string
		petDao      FakePetDao
		petId       string
		expectedPet model.Pet
		expectErr   bool
	}

	// Define tests
	tests := []Test{
		{
			name:        "valid get by id",
			petDao:      FakePetDao{getByIdPet: SamplePet1},
			petId:       SamplePet1.Id,
			expectedPet: SamplePet1,
			expectErr:   false,
		},
		{
			name:        "DAO get error",
			petDao:      FakePetDao{getByIdErr: errors.New("dao error")},
			petId:       SamplePet1.Id,
			expectedPet: SamplePet1,
			expectErr:   true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		service := service.NewPetService(&test.petDao, &SampleEncoder)

		// Execute
		pet, err := service.GetById(test.petId)

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
			assert.Equal(t, test.expectedPet, pet, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestPetDelete(t *testing.T) {
	// Define test struct
	type Test struct {
		name      string
		petDao    FakePetDao
		petId     string
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid delete",
			petDao:    FakePetDao{deleteErr: nil},
			petId:     SamplePet1.Id,
			expectErr: false,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		service := service.NewPetService(&test.petDao, &SampleEncoder)

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

func TestPetList(t *testing.T) {
	// Define test struct
	type Test struct {
		name               string
		petDao             FakePetDao
		encoder            FakeEncoder
		first              int
		after              string
		expectedConnection model.PetConnection
		expectErr          bool
	}

	// Define tests
	tests := []Test{
		{
			name: "list all pets",
			petDao: FakePetDao{
				getTotalCountValue: 2,
				queryPets:          []model.Pet{SamplePet1, SamplePet2},
				queryHasNextPage:   false,
			},
			encoder: SampleEncoder,
			first:   10,
			after:   "",
			expectedConnection: model.PetConnection{
				TotalCount: 2,
				Edges: []model.PetEdge{
					SamplePet1Edge,
					SamplePet2Edge,
				},
				PageInfo: model.PageInfo{
					EndCursor:   SampleEncoder.Encode(SamplePet2.Id),
					HasNextPage: false,
				},
			},
			expectErr: false,
		},
		{
			name: "list first of two pets",
			petDao: FakePetDao{
				getTotalCountValue: 2,
				queryPets:          []model.Pet{SamplePet1},
				queryHasNextPage:   true,
			},
			encoder: SampleEncoder,
			first:   1,
			after:   "",
			expectedConnection: model.PetConnection{
				TotalCount: 2,
				Edges: []model.PetEdge{
					SamplePet1Edge,
				},
				PageInfo: model.PageInfo{
					EndCursor:   SampleEncoder.Encode(SamplePet1.Id),
					HasNextPage: true,
				},
			},
			expectErr: false,
		},
		{
			name: "list second of two pets",
			petDao: FakePetDao{
				getTotalCountValue: 2,
				queryPets:          []model.Pet{SamplePet2},
				queryHasNextPage:   false,
			},
			encoder: SampleEncoder,
			first:   1,
			after:   SampleEncoder.Encode("token"),
			expectedConnection: model.PetConnection{
				TotalCount: 2,
				Edges: []model.PetEdge{
					SamplePet2Edge,
				},
				PageInfo: model.PageInfo{
					EndCursor:   SampleEncoder.Encode(SamplePet2.Id),
					HasNextPage: false,
				},
			},
			expectErr: false,
		},
		{
			name: "decode error",
			petDao: FakePetDao{
				getTotalCountValue: 2,
				queryPets:          []model.Pet{SamplePet1, SamplePet2},
				queryHasNextPage:   false,
			},
			encoder: FakeEncoder{
				decodeErr: assert.AnError,
			},
			first:     1,
			after:     "",
			expectErr: true,
		},
		{
			name: "DAO total count error",
			petDao: FakePetDao{
				getTotalCountErr: assert.AnError,
			},
			encoder:   SampleEncoder,
			first:     1,
			after:     "",
			expectErr: true,
		},
		{
			name: "DAO query error",
			petDao: FakePetDao{
				getTotalCountValue: 2,
				queryErr:           assert.AnError,
			},
			encoder:   SampleEncoder,
			first:     1,
			after:     "",
			expectErr: true,
		},
	}

	//Run tests
	for _, test := range tests {
		// Setup
		service := service.NewPetService(&test.petDao, &test.encoder)

		// Execute
		pets, err := service.List(test.first, test.after)

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
			assert.Equal(t, test.expectedConnection, pets, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestPetUpdateOwner(t *testing.T) {
	// Define test struct
	type Test struct {
		name      string
		petDao    FakePetDao
		petId     string
		petOwner  string
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name: "valid update owner",
			petDao: FakePetDao{
				getByIdPet: SamplePet1,
			},
			petId:     SamplePet1.Id,
			petOwner:  SamplePet2.Owner,
			expectErr: false,
		},
		{
			name: "DAO get error",
			petDao: FakePetDao{
				getByIdErr: assert.AnError,
			},
			petId:     SamplePet1.Id,
			petOwner:  SamplePet2.Owner,
			expectErr: true,
		},
		{
			name: "DAO update error",
			petDao: FakePetDao{
				updateErr: assert.AnError,
			},
			petId:     SamplePet1.Id,
			petOwner:  SamplePet2.Owner,
			expectErr: true,
		},
	}

	// Run
	for _, test := range tests {
		// Setup
		service := service.NewPetService(&test.petDao, &SampleEncoder)

		// Execute
		pet, err := service.UpdateOwner(test.petId, test.petOwner)

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
			assert.Equal(t, test.petOwner, pet.Owner)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}
