package controller_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/controller"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/stretchr/testify/assert"
)

// Define mocks / stubs
type fakePetService struct {
	createOutput      model.Pet
	createErr         error
	deleteErr         error
	getByIdOutput     model.Pet
	getByIdErr        error
	listOutput        model.PetConnection
	listErr           error
	updateOwnerOutput model.Pet
	updateOwnerErr    error
}

// Define mock / stub behavior
func (s *fakePetService) Create(name string, age int, owner string) (model.Pet, error) {
	return s.createOutput, s.createErr
}
func (s *fakePetService) Delete(id string) error {
	return s.deleteErr
}
func (s *fakePetService) GetById(id string) (model.Pet, error) {
	return s.getByIdOutput, s.getByIdErr
}
func (s *fakePetService) List(first int, after string) (model.PetConnection, error) {
	return s.listOutput, s.listErr
}
func (s *fakePetService) UpdateOwner(id string, owner string) (model.Pet, error) {
	return s.updateOwnerOutput, s.updateOwnerErr
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
			name: "create with all args",
			petService: fakePetService{
				createOutput: samplePet,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"name":  samplePet.Name,
					"age":   float64(samplePet.Age),
					"owner": samplePet.Owner,
				}},
			},
			expectedResponse: controller.Response{
				Data: model.CreatePetPayload{
					Pet: samplePet,
				}},
			expectErr: false,
		},
		{
			name: "create without optional args",
			petService: fakePetService{
				createOutput: samplePet,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"name": samplePet.Name,
					"age":  float64(samplePet.Age),
				}},
			},
			expectedResponse: controller.Response{
				Data: model.CreatePetPayload{
					Pet: samplePet,
				}},
			expectErr: false,
		},
		{
			name: "service create error",
			petService: fakePetService{
				createErr: assert.AnError,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"name":  samplePet.Name,
					"age":   float64(samplePet.Age),
					"owner": samplePet.Owner,
				}},
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
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"id": samplePet.Id,
				}},
			},
			expectedResponse: controller.Response{
				Data: model.DeletePetPayload{
					Id: samplePet.Id,
				},
			},
			expectErr: false,
		},
		{
			name: "service delete error",
			petService: fakePetService{
				deleteErr: assert.AnError,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"id": samplePet.Id,
				}},
			},
			expectErr: true,
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
			name: "valid get",
			petService: fakePetService{
				getByIdOutput: samplePet,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"id": samplePet.Id,
				}},
			},
			expectedResponse: controller.Response{
				Data: samplePet,
			},
			expectErr: false,
		},
		{
			name: "service get error",
			petService: fakePetService{
				getByIdErr: assert.AnError,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"id": samplePet.Id,
				}},
			},
			expectErr: true,
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
			name: "list without input",
			petService: fakePetService{
				listOutput: sampleConnection,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{},
			},
			expectedResponse: controller.Response{
				Data: sampleConnection,
			},
			expectErr: false,
		},
		{
			name:       "list with input",
			petService: fakePetService{listOutput: sampleConnection},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"first": float64(10),
					"after": "some cursor value",
				}},
			},
			expectedResponse: controller.Response{
				Data: sampleConnection,
			},
			expectErr: false,
		},
		{
			name:       "service list error",
			petService: fakePetService{listErr: assert.AnError},
			request: controller.Request{
				Arguments: map[string]interface{}{},
			},
			expectErr: true,
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

func TestHandleUpdateOwner(t *testing.T) {
	// Define tests
	tests := []Test{
		{
			name: "valid update owner",
			petService: fakePetService{
				updateOwnerOutput: samplePet,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{
					"id":    samplePet.Id,
					"owner": samplePet.Owner,
				},
			},
			expectedResponse: controller.Response{
				Data: model.UpdatePetOwnerPayload{
					Pet: samplePet,
				},
			},
			expectErr: false,
		},
		{
			name:       "service update owner error",
			petService: fakePetService{updateOwnerErr: assert.AnError},
			request: controller.Request{
				Arguments: map[string]interface{}{
					"id":    samplePet.Id,
					"owner": samplePet.Owner,
				},
			},
			expectedResponse: controller.Response{
				Data: model.UpdatePetOwnerPayload{
					Pet: samplePet,
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
		response := controller.HandleUpdateOwner(test.request)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedResponse, response, test.name)
		} else {
			assert.NotNil(t, response.Error, test.name)
		}
	}
}
