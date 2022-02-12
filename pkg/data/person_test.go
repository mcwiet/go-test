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
type fakeDbClient struct {
	deleteItemOutput *dynamodb.DeleteItemOutput
	deleteItemErr    error
	getItemOutput    *dynamodb.GetItemOutput
	getItemErr       error
	putItemOutput    *dynamodb.PutItemOutput
	putItemErr       error
	queryOutput      *dynamodb.QueryOutput
	queryErr         error
}

// Define mock / stub behavior
func (f *fakeDbClient) DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return f.deleteItemOutput, f.deleteItemErr
}
func (f *fakeDbClient) GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return f.getItemOutput, f.getItemErr
}
func (f *fakeDbClient) PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return f.putItemOutput, f.putItemErr
}
func (f *fakeDbClient) Query(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return f.queryOutput, f.queryErr
}

// Define common data
var (
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
		dbClient  fakeDbClient
		personId  string
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid delete",
			dbClient:  fakeDbClient{},
			personId:  samplePerson1.Id,
			expectErr: false,
		},
		{
			name: "db delete error",
			dbClient: fakeDbClient{
				deleteItemErr: assert.AnError,
			},
			personId:  samplePerson1.Id,
			expectErr: true,
		},
		{
			name: "db item not found error",
			dbClient: fakeDbClient{
				deleteItemErr: &dynamodb.ConditionalCheckFailedException{},
			},
			personId:  samplePerson1.Id,
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
		dbClient       fakeDbClient
		personId       string
		expectedPerson model.Person
		expectErr      bool
	}

	// Define tests
	tests := []Test{
		{
			name: "person found",
			dbClient: fakeDbClient{
				getItemOutput: &dynamodb.GetItemOutput{Item: samplePersonItem1},
			},
			personId:       samplePerson1.Id,
			expectedPerson: samplePerson1,
			expectErr:      false,
		},
		{
			name: "person not found",
			dbClient: fakeDbClient{
				getItemOutput: &dynamodb.GetItemOutput{Item: nil},
			},
			personId:  samplePerson1.Id,
			expectErr: true,
		},
		{
			name: "db get error",
			dbClient: fakeDbClient{
				getItemErr: assert.AnError,
			},
			personId:  samplePerson1.Id,
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
		dbClient  fakeDbClient
		person    model.Person
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid insert",
			dbClient:  fakeDbClient{},
			person:    samplePerson1,
			expectErr: false,
		},
		{
			name: "db put error",
			dbClient: fakeDbClient{
				putItemErr: assert.AnError,
			},
			person:    samplePerson1,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPersonDao(&test.dbClient, tableName)

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

func TestQuery(t *testing.T) {
	// Define test struct
	type Test struct {
		name                string
		dbClient            fakeDbClient
		count               int
		exclusiveStartId    string
		expectedPeople      []model.Person
		expectedHasNextPage bool
		expectErr           bool
	}

	// Define tests
	tests := []Test{
		{
			name: "request more items than in DB",
			dbClient: fakeDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            newInt64(2),
					Items:            []DynamoItem{samplePersonItem1, samplePersonItem2},
					LastEvaluatedKey: DynamoItem{},
				}},
			count:            3,
			exclusiveStartId: "",
			expectedPeople: []model.Person{
				samplePerson1,
				samplePerson2,
			},
			expectedHasNextPage: false,
		},
		{
			name: "request less items than in DB (start at beginning of list)",
			dbClient: fakeDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            newInt64(2),
					Items:            []DynamoItem{samplePersonItem1},
					LastEvaluatedKey: samplePersonItem1,
				}},
			count: 1,
			expectedPeople: []model.Person{
				samplePerson1,
			},
			expectedHasNextPage: true,
		},
		{
			name: "request less items than in DB (reach end of list)",
			dbClient: fakeDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            newInt64(2),
					Items:            []DynamoItem{samplePersonItem2},
					LastEvaluatedKey: DynamoItem{},
				}},
			count:            1,
			exclusiveStartId: samplePerson1.Id,
			expectedPeople: []model.Person{
				samplePerson2,
			},
			expectedHasNextPage: false,
		},
		{
			name: "request 'count=0' but 'exclusiveStartId' is not last person",
			dbClient: fakeDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            newInt64(2),
					Items:            []DynamoItem{},
					LastEvaluatedKey: samplePersonItem2,
				}},
			count:               0,
			exclusiveStartId:    samplePerson1.Id,
			expectedPeople:      []model.Person{},
			expectedHasNextPage: true,
		},
		{
			name: "request 'count=0' but 'exclusiveStartId' is last person",
			dbClient: fakeDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            newInt64(2),
					Items:            []DynamoItem{},
					LastEvaluatedKey: DynamoItem{},
				}},
			count:               0,
			exclusiveStartId:    samplePerson2.Id,
			expectedPeople:      []model.Person{},
			expectedHasNextPage: false,
		},
		{
			name: "no people returned when 'count=0'",
			dbClient: fakeDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            newInt64(1),
					Items:            []DynamoItem{samplePersonItem1},
					LastEvaluatedKey: samplePersonItem1,
				},
			},
			expectedPeople:      []model.Person{},
			expectedHasNextPage: true,
		},
		{
			name: "db query error",
			dbClient: fakeDbClient{
				queryErr: assert.AnError,
			},
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPersonDao(&test.dbClient, tableName)

		// Execute
		people, hasNextPage, err := dao.Query(test.count, test.exclusiveStartId)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedPeople, people, test.name)
			assert.Equal(t, test.expectedHasNextPage, hasNextPage, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestGetTotalCount(t *testing.T) {
	// Define test struct
	type Test struct {
		name          string
		dbClient      fakeDbClient
		expectedCount int
		expectErr     bool
	}

	// Define tests
	tests := []Test{
		{
			name: "valid get total count",
			dbClient: fakeDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count: newInt64(10),
				},
			},
			expectedCount: 10,
		},
		{
			name: "db query error",
			dbClient: fakeDbClient{
				queryErr: assert.AnError,
			},
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPersonDao(&test.dbClient, tableName)

		// Execute
		count, err := dao.GetTotalCount()

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedCount, count, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func newInt64(val int) *int64 {
	valInt64 := int64(val)
	return &valInt64
}
