package controller_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/controller"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/stretchr/testify/assert"
)

var (
	SamplePet = model.Pet{
		Id:    uuid.NewString(),
		Name:  "Levi",
		Age:   1,
		Owner: "User",
	}
	SampleConnection = model.PetConnection{
		TotalCount: 1,
		Edges:      []model.PetEdge{{Node: SamplePet, Cursor: "cursor"}},
		PageInfo:   model.PageInfo{EndCursor: "cursor", HasNextPage: false},
	}
)

// Define test struct
type Test struct {
	name             string
	petService       FakePetService
	request          controller.Request
	expectedResponse controller.Response
	expectErr        bool
}

func TestHandleCreate(t *testing.T) {
	// Define tests
	tests := []Test{
		{
			name: "create with all args",
			petService: FakePetService{
				createOutput: SamplePet,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"name":  SamplePet.Name,
					"age":   float64(SamplePet.Age),
					"owner": SamplePet.Owner,
				}},
			},
			expectedResponse: controller.Response{
				Data: model.CreatePetPayload{
					Pet: SamplePet,
				}},
			expectErr: false,
		},
		{
			name: "create without optional args",
			petService: FakePetService{
				createOutput: SamplePet,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"name": SamplePet.Name,
					"age":  float64(SamplePet.Age),
				}},
			},
			expectedResponse: controller.Response{
				Data: model.CreatePetPayload{
					Pet: SamplePet,
				}},
			expectErr: false,
		},
		{
			name: "service create error",
			petService: FakePetService{
				createErr: assert.AnError,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"name":  SamplePet.Name,
					"age":   float64(SamplePet.Age),
					"owner": SamplePet.Owner,
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
			petService: FakePetService{},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"id": SamplePet.Id,
				}},
			},
			expectedResponse: controller.Response{
				Data: model.DeletePetPayload{
					Id: SamplePet.Id,
				},
			},
			expectErr: false,
		},
		{
			name: "service delete error",
			petService: FakePetService{
				deleteErr: assert.AnError,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"id": SamplePet.Id,
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
			petService: FakePetService{
				getByIdOutput: SamplePet,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"id": SamplePet.Id,
				}},
			},
			expectedResponse: controller.Response{
				Data: SamplePet,
			},
			expectErr: false,
		},
		{
			name: "service get error",
			petService: FakePetService{
				getByIdErr: assert.AnError,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"id": SamplePet.Id,
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
			petService: FakePetService{
				listOutput: SampleConnection,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{},
			},
			expectedResponse: controller.Response{
				Data: SampleConnection,
			},
			expectErr: false,
		},
		{
			name:       "list with input",
			petService: FakePetService{listOutput: SampleConnection},
			request: controller.Request{
				Arguments: map[string]interface{}{"input": map[string]interface{}{
					"first": float64(10),
					"after": "some cursor value",
				}},
			},
			expectedResponse: controller.Response{
				Data: SampleConnection,
			},
			expectErr: false,
		},
		{
			name:       "service list error",
			petService: FakePetService{listErr: assert.AnError},
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
			petService: FakePetService{
				updateOwnerOutput: SamplePet,
			},
			request: controller.Request{
				Arguments: map[string]interface{}{
					"id":    SamplePet.Id,
					"owner": SamplePet.Owner,
				},
			},
			expectedResponse: controller.Response{
				Data: model.UpdatePetOwnerPayload{
					Pet: SamplePet,
				},
			},
			expectErr: false,
		},
		{
			name:       "service update owner error",
			petService: FakePetService{updateOwnerErr: assert.AnError},
			request: controller.Request{
				Arguments: map[string]interface{}{
					"id":    SamplePet.Id,
					"owner": SamplePet.Owner,
				},
			},
			expectedResponse: controller.Response{
				Data: model.UpdatePetOwnerPayload{
					Pet: SamplePet,
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
