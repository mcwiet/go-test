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
type fakePersonDao struct {
	deleteErr          error
	getByIdPerson      model.Person
	getByIdErr         error
	getTotalCountValue int
	getTotalCountErr   error
	insertErr          error
	queryPeople        []model.Person
	queryHasNextPage   bool
	queryErr           error
}
type personDaoGetTotalCount = service.PersonDao
type fakeEncoder struct{}

// Define mock / stub behavior
func (f fakePersonDao) Delete(string) error {
	return f.deleteErr
}
func (f fakePersonDao) GetById(string) (model.Person, error) {
	return f.getByIdPerson, f.getByIdErr
}
func (f fakePersonDao) GetTotalCount() (int, error) {
	return f.getTotalCountValue, f.getTotalCountErr
}
func (f fakePersonDao) Insert(model.Person) error {
	return f.insertErr
}
func (f fakePersonDao) Query(count int, exclusiveStartId string) ([]model.Person, bool, error) {
	return f.queryPeople, f.queryHasNextPage, f.queryErr
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
	encoder       = fakeEncoder{}
	samplePerson1 = model.Person{
		Id:   uuid.NewString(),
		Name: "person 1",
		Age:  12,
	}
	samplePerson2 = model.Person{
		Id:   uuid.NewString(),
		Name: "person 2",
		Age:  20,
	}
)

func TestCreate(t *testing.T) {
	// Define test struct
	type Test struct {
		name       string
		personDao  fakePersonDao
		personName string
		personAge  int
		expectErr  bool
	}

	// Define tests
	tests := []Test{
		{
			name:       "valid create",
			personDao:  fakePersonDao{},
			personName: samplePerson1.Name,
			personAge:  samplePerson1.Age,
			expectErr:  false,
		},
		{
			name:       "DAO insert error",
			personDao:  fakePersonDao{insertErr: errors.New("dao error")},
			personName: samplePerson1.Name,
			personAge:  samplePerson1.Age,
			expectErr:  true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		service := service.NewPersonService(test.personDao, &encoder)

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
		personDao      fakePersonDao
		personId       string
		expectedPerson model.Person
		expectErr      bool
	}

	// Define tests
	tests := []Test{
		{
			name:           "valid get by id",
			personDao:      fakePersonDao{getByIdPerson: samplePerson1},
			personId:       samplePerson1.Id,
			expectedPerson: samplePerson1,
			expectErr:      false,
		},
		{
			name:           "DAO get error",
			personDao:      fakePersonDao{getByIdErr: errors.New("dao error")},
			personId:       samplePerson1.Id,
			expectedPerson: samplePerson1,
			expectErr:      true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		service := service.NewPersonService(test.personDao, &encoder)

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
		personDao fakePersonDao
		personId  string
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid delete",
			personDao: fakePersonDao{deleteErr: nil},
			personId:  samplePerson1.Id,
			expectErr: false,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		service := service.NewPersonService(test.personDao, &encoder)

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
		name               string
		personDao          fakePersonDao
		first              int
		after              string
		expectedConnection model.PersonConnection
		expectErr          bool
	}

	// Define tests
	tests := []Test{
		{
			name: "list all people",
			personDao: fakePersonDao{
				getTotalCountValue: 2,
				queryPeople:        []model.Person{samplePerson1, samplePerson2},
				queryHasNextPage:   false,
			},
			first: 10,
			after: "",
			expectedConnection: model.PersonConnection{
				TotalCount: 2,
				Edges: []model.PersonEdge{
					{
						Node:   samplePerson1,
						Cursor: encoder.Encode(samplePerson1.Id),
					},
					{
						Node:   samplePerson2,
						Cursor: encoder.Encode(samplePerson2.Id),
					},
				},
				PageInfo: model.PageInfo{
					EndCursor:   encoder.Encode(samplePerson2.Id),
					HasNextPage: false,
				},
			},
			expectErr: false,
		},
		{
			name: "list first of two people",
			personDao: fakePersonDao{
				getTotalCountValue: 2,
				queryPeople:        []model.Person{samplePerson1},
				queryHasNextPage:   true,
			},
			first: 1,
			after: "",
			expectedConnection: model.PersonConnection{
				TotalCount: 2,
				Edges: []model.PersonEdge{
					{
						Node:   samplePerson1,
						Cursor: encoder.Encode(samplePerson1.Id),
					},
				},
				PageInfo: model.PageInfo{
					EndCursor:   encoder.Encode(samplePerson1.Id),
					HasNextPage: true,
				},
			},
			expectErr: false,
		},
		{
			name: "list second of two people",
			personDao: fakePersonDao{
				getTotalCountValue: 2,
				queryPeople:        []model.Person{samplePerson2},
				queryHasNextPage:   false,
			},
			first: 1,
			after: encoder.Encode(samplePerson1.Id),
			expectedConnection: model.PersonConnection{
				TotalCount: 2,
				Edges: []model.PersonEdge{
					{
						Node:   samplePerson2,
						Cursor: encoder.Encode(samplePerson2.Id),
					},
				},
				PageInfo: model.PageInfo{
					EndCursor:   encoder.Encode(samplePerson2.Id),
					HasNextPage: false,
				},
			},
			expectErr: false,
		},
		{
			name: "DAO total count error",
			personDao: fakePersonDao{
				getTotalCountErr: assert.AnError,
			},
			first:     1,
			after:     "",
			expectErr: true,
		},
		{
			name: "DAO query error",
			personDao: fakePersonDao{
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
		service := service.NewPersonService(test.personDao, &encoder)

		// Execute
		people, err := service.List(test.first, test.after)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedConnection, people, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}
