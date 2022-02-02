package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/stretchr/testify/assert"
)

type mockPersonDao struct {
	valReturn interface{}
	errReturn error
}

func (m mockPersonDao) Delete(string) error {
	return m.errReturn
}
func (m mockPersonDao) GetById(string) (*model.Person, error) {
	ret, _ := m.valReturn.(*model.Person)
	return ret, m.errReturn
}
func (m mockPersonDao) Insert(*model.Person) error {
	return m.errReturn
}
func (m mockPersonDao) List() (*[]model.Person, error) {
	ret, _ := m.valReturn.(*[]model.Person)
	return ret, m.errReturn
}

var (
	samplePerson = model.Person{
		Id:   uuid.NewString(),
		Name: "dummy",
		Age:  12,
	}
)

func TestCreate(t *testing.T) {
	// Define test struct
	type Test struct {
		name       string
		personDao  mockPersonDao
		personName string
		personAge  int
		expectErr  bool
	}

	// Define tests
	tests := []Test{
		{
			name:       "valid input",
			personDao:  mockPersonDao{},
			personName: samplePerson.Name,
			personAge:  samplePerson.Age,
			expectErr:  false,
		},
		{
			name:       "DAO insert error",
			personDao:  mockPersonDao{errReturn: errors.New("dao error")},
			personName: samplePerson.Name,
			personAge:  samplePerson.Age,
			expectErr:  true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		service := NewPersonService(test.personDao)

		// Execute
		person, err := service.Create(test.personName, test.personAge)

		// Verify
		if !test.expectErr {
			_, uuidErr := uuid.Parse(person.Id)
			assert.Equal(t, person.Name, test.personName, test.name)
			assert.Equal(t, person.Age, test.personAge, test.name)
			assert.Nil(t, uuidErr, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestGetById(t *testing.T) {
	// Define test struct
	type Test struct {
		name           string
		personDao      mockPersonDao
		personId       string
		expectedPerson *model.Person
		expectErr      bool
	}

	// Define tests
	tests := []Test{
		{
			name:           "valid input",
			personDao:      mockPersonDao{valReturn: &samplePerson},
			personId:       samplePerson.Id,
			expectedPerson: &samplePerson,
			expectErr:      false,
		},
		{
			name:           "DAO get error",
			personDao:      mockPersonDao{errReturn: errors.New("dao error")},
			personId:       samplePerson.Id,
			expectedPerson: &samplePerson,
			expectErr:      true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		service := NewPersonService(test.personDao)

		// Execute
		person, err := service.GetById(test.personId)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedPerson, person, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestDelete(t *testing.T) {
	// Define test struct
	type Test struct {
		name      string
		personDao mockPersonDao
		personId  string
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid input",
			personDao: mockPersonDao{errReturn: nil},
			personId:  samplePerson.Id,
			expectErr: false,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		service := NewPersonService(test.personDao)

		// Execute
		err := service.Delete(test.personId)

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
		name           string
		personDao      mockPersonDao
		expectedPeople *[]model.Person
		expectErr      bool
	}

	// Define tests
	tests := []Test{
		{
			name:           "valid input",
			personDao:      mockPersonDao{valReturn: &[]model.Person{samplePerson}},
			expectedPeople: &[]model.Person{samplePerson},
			expectErr:      false,
		},
	}

	//Run tests
	for _, test := range tests {
		// Setup
		service := NewPersonService(test.personDao)

		// Execute
		people, err := service.List()

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedPeople, people, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}
