package data_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/jsii-runtime-go"
	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/data"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/stretchr/testify/assert"
)

// Define mocks / stubs
type fakeDynamoDbClient struct {
	returnedValue interface{}
	returnedErr   error
}

// Define mock / sub behavior
func (c *fakeDynamoDbClient) DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	ret, _ := c.returnedValue.(*dynamodb.DeleteItemOutput)
	return ret, c.returnedErr
}
func (c *fakeDynamoDbClient) GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	ret, _ := c.returnedValue.(*dynamodb.GetItemOutput)
	return ret, c.returnedErr
}
func (c *fakeDynamoDbClient) PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	ret, _ := c.returnedValue.(*dynamodb.PutItemOutput)
	return ret, c.returnedErr
}
func (c *fakeDynamoDbClient) Query(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	ret, _ := c.returnedValue.(*dynamodb.QueryOutput)
	return ret, c.returnedErr
}

// Define common data
var (
	tableName    = "table"
	samplePerson = model.Person{
		Id:   uuid.NewString(),
		Name: "dummy",
		Age:  5,
	}
	samplePersonItem = map[string]*dynamodb.AttributeValue{
		"Id":   {S: &samplePerson.Id},
		"Name": {S: &samplePerson.Name},
		"Age":  {N: jsii.String("5")},
	}
)

func TestDelete(t *testing.T) {
	// Define test struct
	type Test struct {
		name      string
		dbClient  fakeDynamoDbClient
		personId  string
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid delete",
			dbClient:  fakeDynamoDbClient{},
			personId:  samplePerson.Id,
			expectErr: false,
		},
		{
			name:      "db delete error",
			dbClient:  fakeDynamoDbClient{returnedErr: assert.AnError},
			personId:  samplePerson.Id,
			expectErr: true,
		},
		{
			name:      "db item not found error",
			dbClient:  fakeDynamoDbClient{returnedErr: &dynamodb.ConditionalCheckFailedException{}},
			personId:  samplePerson.Id,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPersonDao(&test.dbClient, tableName)

		// Execute
		err := dao.Delete(test.personId)

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestGetById(t *testing.T) {
	// Define test struct
	type Test struct {
		name           string
		dbClient       fakeDynamoDbClient
		personId       string
		expectedPerson *model.Person
		expectErr      bool
	}

	// Define tests
	tests := []Test{
		{
			name:           "valid get by id",
			dbClient:       fakeDynamoDbClient{returnedValue: &dynamodb.GetItemOutput{Item: samplePersonItem}},
			personId:       samplePerson.Id,
			expectedPerson: &samplePerson,
			expectErr:      false,
		},
		{
			name:      "db get error",
			dbClient:  fakeDynamoDbClient{returnedErr: assert.AnError},
			personId:  samplePerson.Id,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPersonDao(&test.dbClient, tableName)

		// Execute
		person, err := dao.GetById(test.personId)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedPerson, person, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestInsert(t *testing.T) {
	// Define test struct
	type Test struct {
		name      string
		dbClient  fakeDynamoDbClient
		person    model.Person
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid insert",
			dbClient:  fakeDynamoDbClient{},
			person:    samplePerson,
			expectErr: false,
		},
		{
			name:      "db put error",
			dbClient:  fakeDynamoDbClient{returnedErr: assert.AnError},
			person:    samplePerson,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPersonDao(&test.dbClient, tableName)

		// Execute
		err := dao.Insert(&test.person)

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
		dbClient           fakeDynamoDbClient
		expectedConnection model.PersonConnection
		expectErr          bool
	}

	// Define tests
	tests := []Test{
		{
			name: "valid list",
			dbClient: fakeDynamoDbClient{returnedValue: &dynamodb.QueryOutput{
				Count:            newInt64(1),
				Items:            []map[string]*dynamodb.AttributeValue{samplePersonItem},
				LastEvaluatedKey: map[string]*dynamodb.AttributeValue{},
			}},
			expectedConnection: model.PersonConnection{
				TotalCount: 1,
				Edges: []model.PersonEdge{
					{
						Node:   samplePerson,
						Cursor: samplePerson.Id,
					}},
				PageInfo: model.PageInfo{
					EndCursor:   samplePerson.Id,
					HasNextPage: false,
				},
			},
			expectErr: false,
		},
		{
			name:      "db query error",
			dbClient:  fakeDynamoDbClient{returnedErr: assert.AnError},
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPersonDao(&test.dbClient, tableName)

		// Execute
		connection, err := dao.List()

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedConnection, connection, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func newInt64(val int64) *int64 {
	return &val
}
