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

type DynamoItem = map[string]*dynamodb.AttributeValue

// Define mocks / stubs
type fakeDynamoDbClient struct {
	returnedValue interface{}
	returnedErr   error
}
type fakeEncoder struct{}

// Define mock / stub behavior
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

func (e *fakeEncoder) Encode(input string) string {
	return input
}
func (e *fakeEncoder) Decode(input string) (string, error) {
	return input, nil
}

// Define common data
var (
	encoder       = fakeEncoder{}
	tableName     = "table"
	samplePerson1 = model.Person{
		Id:   uuid.NewString(),
		Name: "person 1",
		Age:  10,
	}
	samplePerson2 = model.Person{
		Id:   uuid.NewString(),
		Name: "person 2",
		Age:  92,
	}
	samplePersonItem1 = DynamoItem{
		"Id":   {S: &samplePerson1.Id},
		"Name": {S: &samplePerson1.Name},
		"Age":  {N: jsii.String("10")},
	}
	samplePersonItem2 = DynamoItem{
		"Id":   {S: &samplePerson2.Id},
		"Name": {S: &samplePerson2.Name},
		"Age":  {N: jsii.String("92")},
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
			personId:  samplePerson1.Id,
			expectErr: false,
		},
		{
			name:      "db delete error",
			dbClient:  fakeDynamoDbClient{returnedErr: assert.AnError},
			personId:  samplePerson1.Id,
			expectErr: true,
		},
		{
			name:      "db item not found error",
			dbClient:  fakeDynamoDbClient{returnedErr: &dynamodb.ConditionalCheckFailedException{}},
			personId:  samplePerson1.Id,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPersonDao(&test.dbClient, &encoder, tableName)

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
		expectedPerson model.Person
		expectErr      bool
	}

	// Define tests
	tests := []Test{
		{
			name:           "valid get by id",
			dbClient:       fakeDynamoDbClient{returnedValue: &dynamodb.GetItemOutput{Item: samplePersonItem1}},
			personId:       samplePerson1.Id,
			expectedPerson: samplePerson1,
			expectErr:      false,
		},
		{
			name:      "db get error",
			dbClient:  fakeDynamoDbClient{returnedErr: assert.AnError},
			personId:  samplePerson1.Id,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPersonDao(&test.dbClient, &encoder, tableName)

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
			person:    samplePerson1,
			expectErr: false,
		},
		{
			name:      "db put error",
			dbClient:  fakeDynamoDbClient{returnedErr: assert.AnError},
			person:    samplePerson1,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPersonDao(&test.dbClient, &encoder, tableName)

		// Execute
		err := dao.Insert(test.person)

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
		first              int
		after              string
		expectedConnection model.PersonConnection
		expectErr          bool
	}

	// Define tests
	tests := []Test{
		{
			name: "request more items than in DB",
			dbClient: fakeDynamoDbClient{returnedValue: &dynamodb.QueryOutput{
				Count:            newInt64(2),
				Items:            []DynamoItem{samplePersonItem1, samplePersonItem2},
				LastEvaluatedKey: DynamoItem{},
			}},
			first: 3,
			after: "",
			expectedConnection: model.PersonConnection{
				TotalCount: 2,
				Edges: []model.PersonEdge{
					{
						Node:   samplePerson1,
						Cursor: samplePerson1.Id,
					},
					{
						Node:   samplePerson2,
						Cursor: samplePerson2.Id,
					},
				},
				PageInfo: model.PageInfo{
					EndCursor:   samplePerson2.Id,
					HasNextPage: false,
				},
			},
			expectErr: false,
		},
		{
			name: "request less items than in DB, beginning of list",
			dbClient: fakeDynamoDbClient{returnedValue: &dynamodb.QueryOutput{
				Count: newInt64(2),
				Items: []DynamoItem{samplePersonItem1},
				LastEvaluatedKey: DynamoItem{
					"Id": &dynamodb.AttributeValue{S: jsii.String(samplePerson1.Id)},
				},
			}},
			first: 1,
			after: "",
			expectedConnection: model.PersonConnection{
				TotalCount: 2,
				Edges: []model.PersonEdge{
					{
						Node:   samplePerson1,
						Cursor: samplePerson1.Id,
					},
				},
				PageInfo: model.PageInfo{
					EndCursor:   samplePerson1.Id,
					HasNextPage: true,
				},
			},
			expectErr: false,
		},
		{
			name: "request less items than in DB, end of list",
			dbClient: fakeDynamoDbClient{returnedValue: &dynamodb.QueryOutput{
				Count:            newInt64(2),
				Items:            []DynamoItem{samplePersonItem2},
				LastEvaluatedKey: DynamoItem{},
			}},
			first: 1,
			after: samplePerson1.Id,
			expectedConnection: model.PersonConnection{
				TotalCount: 2,
				Edges: []model.PersonEdge{
					{
						Node:   samplePerson2,
						Cursor: samplePerson2.Id,
					},
				},
				PageInfo: model.PageInfo{
					EndCursor:   samplePerson2.Id,
					HasNextPage: false,
				},
			},
			expectErr: false,
		},
		{
			name: "request 'first=0' but 'after' is not last person",
			dbClient: fakeDynamoDbClient{returnedValue: &dynamodb.QueryOutput{
				Count: newInt64(2),
				Items: []DynamoItem{},
				LastEvaluatedKey: DynamoItem{
					"Id": &dynamodb.AttributeValue{S: jsii.String(samplePerson2.Id)},
				},
			}},
			first: 0,
			after: samplePerson1.Id,
			expectedConnection: model.PersonConnection{
				TotalCount: 2,
				Edges:      []model.PersonEdge{},
				PageInfo: model.PageInfo{
					EndCursor:   "",
					HasNextPage: true,
				},
			},
		},
		{
			name: "request 'first=0' but 'after' is last person",
			dbClient: fakeDynamoDbClient{returnedValue: &dynamodb.QueryOutput{
				Count:            newInt64(2),
				Items:            []DynamoItem{},
				LastEvaluatedKey: DynamoItem{},
			}},
			first: 0,
			after: samplePerson2.Id,
			expectedConnection: model.PersonConnection{
				TotalCount: 2,
				Edges:      []model.PersonEdge{},
				PageInfo: model.PageInfo{
					EndCursor:   "",
					HasNextPage: false,
				},
			},
		},
		{
			name: "db query error",
			dbClient: fakeDynamoDbClient{
				returnedValue: &dynamodb.QueryOutput{},
				returnedErr:   assert.AnError,
			},
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPersonDao(&test.dbClient, &encoder, tableName)

		// Execute
		connection, err := dao.List(test.first, test.after)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedConnection, connection, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func newInt64(val int) *int64 {
	valInt64 := int64(val)
	return &valInt64
}
