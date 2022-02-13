package controller_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/controller"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/stretchr/testify/assert"
)

// Define mocks / stubs
type fakePetService struct {
	returnedValue interface{}
	returnedErr   error
}

// Define mock / stub behavior
func (s *fakePetService) Create(name string, age int, owner string) (model.Pet, error) {
	ret, _ := s.returnedValue.(model.Pet)
	return ret, s.returnedErr
}
func (s *fakePetService) Delete(id string) error {
	return s.returnedErr
}
func (s *fakePetService) GetById(id string) (model.Pet, error) {
	ret, _ := s.returnedValue.(model.Pet)
	return ret, s.returnedErr
}
func (s *fakePetService) List(first int, after string) (model.PetConnection, error) {
	ret, _ := s.returnedValue.(model.PetConnection)
	return ret, s.returnedErr
}

// Define common data
var (
	samplePet = model.Pet{
		Id:    uuid.NewString(),
		Name:  "Levi",
		Age:   1,
		Owner: "User",
	}
	sampleConnection = model.PetConnection{
		TotalCount: 1,
		Edges:      []model.PetEdge{{Node: samplePet, Cursor: "cursor"}},
		PageInfo:   model.PageInfo{EndCursor: "cursor", HasNextPage: false},
	}
)

// Define test struct
type Test struct {
	name             string
	petService       fakePetService
	request          controller.Request
	expectedResponse controller.Response
	expectErr        bool
}

func TestHandleCreate(t *testing.T) {
	// Define tests
	tests := []Test{
		{
			name:       "create with all args",
			petService: fakePetService{returnedValue: samplePet},
			request: controller.Request{
				Arguments: map[string]interface{}{
					"name":  samplePet.Name,
					"age":   float64(samplePet.Age),
					"owner": samplePet.Owner,
				},
			},
			expectedResponse: controller.Response{Data: model.CreatePetPayload{Pet: samplePet}},
			expectErr:        false,
		},
		{
			name:       "create without optional args",
			petService: fakePetService{returnedValue: samplePet},
			request: controller.Request{
				Arguments: map[string]interface{}{
					"name": samplePet.Name,
					"age":  float64(samplePet.Age),
				},
			},
			expectedResponse: controller.Response{Data: model.CreatePetPayload{Pet: samplePet}},
			expectErr:        false,
		},
		{
			name:       "service create error",
			petService: fakePetService{returnedErr: errors.New("create error")},
			request: controller.Request{
				Arguments: map[string]interface{}{
					"name":  samplePet.Name,
					"age":   float64(samplePet.Age),
					"owner": samplePet.Owner,
				},
			},
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		controller := controller.NewPetController(&test.petService)

		// Execute
		response := controller.HandleCreate(test.request)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedResponse, response, test.name)
		} else {
			assert.NotNil(t, response.Error, test.name)
		}
	}
}

func TestHandleDelete(t *testing.T) {
	// Define tests
	tests := []Test{
		{
			name:       "valid delete",
			petService: fakePetService{},
			request: controller.Request{
				Arguments: map[string]interface{}{
					"id": samplePet.Id,
				},
			},
			expectedResponse: controller.Response{Data: model.DeletePetPayload{Id: samplePet.Id}},
			expectErr:        false,
		},
		{
			name:       "service delete error",
			petService: fakePetService{returnedErr: errors.New("delete error")},
			request: controller.Request{
				Arguments: map[string]interface{}{
					"id": samplePet.Id,
				},
			},
			expectedResponse: controller.Response{},
			expectErr:        true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		controller := controller.NewPetController(&test.petService)

		// Execute
		response := controller.HandleDelete(test.request)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedResponse, response, test.name)
		} else {
			assert.NotNil(t, response.Error, test.name)
		}
	}
}

func TestHandleGet(t *testing.T) {
	// Define tests
	tests := []Test{
		{
			name:       "valid get",
			petService: fakePetService{returnedValue: samplePet},
			request: controller.Request{
				Arguments: map[string]interface{}{
					"id": samplePet.Id,
				},
			},
			expectedResponse: controller.Response{Data: samplePet},
			expectErr:        false,
		},
		{
			name:       "service get error",
			petService: fakePetService{returnedErr: errors.New("get error")},
			request: controller.Request{
				Arguments: map[string]interface{}{
					"id": samplePet.Id,
				},
			},
			expectedResponse: controller.Response{},
			expectErr:        true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		controller := controller.NewPetController(&test.petService)

		// Execute
		response := controller.HandleGet(test.request)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedResponse, response, test.name)
		} else {
			assert.NotNil(t, response.Error, test.name)
		}
	}
}

func TestHandleList(t *testing.T) {
	// Define tests
	tests := []Test{
		{
			name:       "list without input",
			petService: fakePetService{returnedValue: sampleConnection},
			request: controller.Request{
				Arguments: map[string]interface{}{},
			},
			expectedResponse: controller.Response{Data: sampleConnection},
			expectErr:        false,
		},
		{
			name:       "list with input",
			petService: fakePetService{returnedValue: sampleConnection},
			request: controller.Request{
				Arguments: map[string]interface{}{
					"first": float64(10),
					"after": "some cursor value",
				},
			},
			expectedResponse: controller.Response{Data: sampleConnection},
			expectErr:        false,
		},
		{
			name:       "service list error",
			petService: fakePetService{returnedErr: errors.New("list error")},
			request: controller.Request{
				Arguments: map[string]interface{}{},
			},
			expectedResponse: controller.Response{},
			expectErr:        true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		controller := controller.NewPetController(&test.petService)

		// Execute
		response := controller.HandleList(test.request)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedResponse, response, test.name)
		} else {
			assert.NotNil(t, response.Error, test.name)
		}
	}
}
