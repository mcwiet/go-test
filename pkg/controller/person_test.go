package controller

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/stretchr/testify/assert"
)

// Define mocks / stubs
type fakePersonService struct {
	returnedValue interface{}
	returnedErr   error
}

// Define mock / stub functionality
func (s *fakePersonService) Create(name string, age int) (*model.Person, error) {
	ret, _ := s.returnedValue.(*model.Person)
	return ret, s.returnedErr
}
func (s *fakePersonService) Delete(id string) error {
	return s.returnedErr
}
func (s *fakePersonService) GetById(id string) (*model.Person, error) {
	ret, _ := s.returnedValue.(*model.Person)
	return ret, s.returnedErr
}
func (s *fakePersonService) List() (*[]model.Person, error) {
	ret, _ := s.returnedValue.(*[]model.Person)
	return ret, s.returnedErr
}

// Define common data
var (
	samplePerson = model.Person{
		Id:   uuid.NewString(),
		Name: "dummy",
		Age:  1,
	}
)

// Define test struct
type Test struct {
	name             string
	personService    fakePersonService
	request          Request
	expectedResponse Response
	expectErr        bool
}

func TestHandleCreate(t *testing.T) {
	// Define tests
	tests := []Test{
		{
			name:          "valid input",
			personService: fakePersonService{returnedValue: &samplePerson},
			request: Request{
				Arguments: map[string]interface{}{
					"name": samplePerson.Name,
					"age":  float64(samplePerson.Age),
				},
			},
			expectedResponse: Response{Data: samplePerson},
			expectErr:        false,
		},
		{
			name:          "service create error",
			personService: fakePersonService{returnedErr: errors.New("create error")},
			request: Request{
				Arguments: map[string]interface{}{
					"name": samplePerson.Name,
					"age":  float64(samplePerson.Age),
				},
			},
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		controller := NewPersonController(&test.personService)

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
			name:          "valid input",
			personService: fakePersonService{},
			request: Request{
				Arguments: map[string]interface{}{
					"id": samplePerson.Id,
				},
			},
			expectedResponse: Response{},
			expectErr:        false,
		},
		{
			name:          "service delete error",
			personService: fakePersonService{returnedErr: errors.New("delete error")},
			request: Request{
				Arguments: map[string]interface{}{
					"id": samplePerson.Id,
				},
			},
			expectedResponse: Response{},
			expectErr:        true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		controller := NewPersonController(&test.personService)

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
			name:          "valid input",
			personService: fakePersonService{returnedValue: &samplePerson},
			request: Request{
				Arguments: map[string]interface{}{
					"id": samplePerson.Id,
				},
			},
			expectedResponse: Response{Data: samplePerson},
			expectErr:        false,
		},
		{
			name:          "service get error",
			personService: fakePersonService{returnedErr: errors.New("get error")},
			request: Request{
				Arguments: map[string]interface{}{
					"id": samplePerson.Id,
				},
			},
			expectedResponse: Response{},
			expectErr:        true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		controller := NewPersonController(&test.personService)

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
			name:             "valid input",
			personService:    fakePersonService{returnedValue: &[]model.Person{samplePerson}},
			request:          Request{},
			expectedResponse: Response{Data: []model.Person{samplePerson}},
			expectErr:        false,
		},
		{
			name:             "service list error",
			personService:    fakePersonService{returnedErr: errors.New("list error")},
			request:          Request{},
			expectedResponse: Response{},
			expectErr:        true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		controller := NewPersonController(&test.personService)

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
