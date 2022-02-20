package data_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/jsii-runtime-go"
	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/data"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/assert"
)

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
	tableName  = "table"
	samplePet1 = model.Pet{
		Id:    uuid.NewString(),
		Name:  "pet 1",
		Age:   10,
		Owner: "User1",
	}
	samplePet2 = model.Pet{
		Id:    uuid.NewString(),
		Name:  "pet 2",
		Age:   92,
		Owner: "User2",
	}
	samplePetItem1 = data.DynamoItem{
		"Id":    {S: &samplePet1.Id},
		"Name":  {S: &samplePet1.Name},
		"Age":   {N: jsii.String("10")},
		"Owner": {S: &samplePet1.Owner},
	}
	samplePetItem2 = data.DynamoItem{
		"Id":    {S: &samplePet2.Id},
		"Name":  {S: &samplePet2.Name},
		"Age":   {N: jsii.String("92")},
		"Owner": {S: &samplePet2.Owner},
	}
)

func TestDeletePet(t *testing.T) {
	// Define test struct
	type Test struct {
		name      string
		dbClient  fakeDbClient
		petId     string
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid delete",
			dbClient:  fakeDbClient{},
			petId:     samplePet1.Id,
			expectErr: jsii,
		},
		{
			name: "db delete error",
			dbClient: fakeDbClient{
				deleteItemErr: assert.AnError,
			},
			petId:     samplePet1.Id,
			expectErr: true,
		},
		{
			name: "db item not found error",
			dbClient: fakeDbClient{
				deleteItemErr: &dynamodb.ConditionalCheckFailedException{},
			},
			petId:     samplePet1.Id,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPetDao(&test.dbClient, tableName)

		// Execute
		err := dao.Delete(test.petId)

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestGetByIdPet(t *testing.T) {
	// Define test struct
	type Test struct {
		name        string
		dbClient    fakeDbClient
		petId       string
		expectedPet model.Pet
		expectErr   bool
	}

	// Define tests
	tests := []Test{
		{
			name: "pet found",
			dbClient: fakeDbClient{
				getItemOutput: &dynamodb.GetItemOutput{Item: samplePetItem1},
			},
			petId:       samplePet1.Id,
			expectedPet: samplePet1,
			expectErr:   false,
		},
		{
			name: "pet not found",
			dbClient: fakeDbClient{
				getItemOutput: &dynamodb.GetItemOutput{Item: nil},
			},
			petId:     samplePet1.Id,
			expectErr: true,
		},
		{
			name: "db get error",
			dbClient: fakeDbClient{
				getItemErr: assert.AnError,
			},
			petId:     samplePet1.Id,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPetDao(&test.dbClient, tableName)

		// Execute
		pet, err := dao.GetById(test.petId)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedPet, pet, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestInsertPet(t *testing.T) {
	// Define test struct
	type Test struct {
		name      string
		dbClient  fakeDbClient
		pet       model.Pet
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid insert",
			dbClient:  fakeDbClient{},
			pet:       samplePet1,
			expectErr: false,
		},
		{
			name: "db put error",
			dbClient: fakeDbClient{
				putItemErr: assert.AnError,
			},
			pet:       samplePet1,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPetDao(&test.dbClient, tableName)

		// Execute
		err := dao.Insert(test.pet)

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestQueryPet(t *testing.T) {
	// Define test struct
	type Test struct {
		name                string
		dbClient            fakeDbClient
		count               int
		exclusiveStartId    string
		expectedPets        []model.Pet
		expectedHasNextPage bool
		expectErr           bool
	}

	// Define tests
	tests := []Test{
		{
			name: "request more items than in DB",
			dbClient: fakeDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            pointy.Int64(2),
					Items:            []data.DynamoItem{samplePetItem1, samplePetItem2},
					LastEvaluatedKey: data.DynamoItem{},
				}},
			count:            3,
			exclusiveStartId: "",
			expectedPets: []model.Pet{
				samplePet1,
				samplePet2,
			},
			expectedHasNextPage: false,
		},
		{
			name: "request less items than in DB (start at beginning of list)",
			dbClient: fakeDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            pointy.Int64(2),
					Items:            []data.DynamoItem{samplePetItem1},
					LastEvaluatedKey: samplePetItem1,
				}},
			count: 1,
			expectedPets: []model.Pet{
				samplePet1,
			},
			expectedHasNextPage: true,
		},
		{
			name: "request less items than in DB (reach end of list)",
			dbClient: fakeDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            pointy.Int64(2),
					Items:            []data.DynamoItem{samplePetItem2},
					LastEvaluatedKey: data.DynamoItem{},
				}},
			count:            1,
			exclusiveStartId: samplePet1.Id,
			expectedPets: []model.Pet{
				samplePet2,
			},
			expectedHasNextPage: false,
		},
		{
			name: "request 'count=0' but 'exclusiveStartId' is not last pet",
			dbClient: fakeDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            pointy.Int64(2),
					Items:            []data.DynamoItem{},
					LastEvaluatedKey: samplePetItem2,
				}},
			count:               0,
			exclusiveStartId:    samplePet1.Id,
			expectedPets:        []model.Pet{},
			expectedHasNextPage: true,
		},
		{
			name: "request 'count=0' but 'exclusiveStartId' is last pet",
			dbClient: fakeDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            pointy.Int64(2),
					Items:            []data.DynamoItem{},
					LastEvaluatedKey: data.DynamoItem{},
				}},
			count:               0,
			exclusiveStartId:    samplePet2.Id,
			expectedPets:        []model.Pet{},
			expectedHasNextPage: false,
		},
		{
			name: "no pets returned when 'count=0'",
			dbClient: fakeDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            pointy.Int64(1),
					Items:            []data.DynamoItem{samplePetItem1},
					LastEvaluatedKey: samplePetItem1,
				},
			},
			expectedPets:        []model.Pet{},
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
		dao := data.NewPetDao(&test.dbClient, tableName)

		// Execute
		pets, hasNextPage, err := dao.Query(test.count, test.exclusiveStartId)

		// Verify
		if !test.expectErr {
			assert.Equal(t, test.expectedPets, pets, test.name)
			assert.Equal(t, test.expectedHasNextPage, hasNextPage, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestGetTotalCountPet(t *testing.T) {
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
					Count: pointy.Int64(10),
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
		dao := data.NewPetDao(&test.dbClient, tableName)

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

func TestUpdatePet(t *testing.T) {
	// Define test struct
	type Test struct {
		name      string
		dbClient  fakeDbClient
		pet       model.Pet
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid update",
			dbClient:  fakeDbClient{},
			pet:       samplePet1,
			expectErr: false,
		},
		{
			name: "db put error",
			dbClient: fakeDbClient{
				putItemErr: assert.AnError,
			},
			pet:       samplePet1,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPetDao(&test.dbClient, tableName)

		// Execute
		err := dao.Update(test.pet)

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}
