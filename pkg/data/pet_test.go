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

var (
	SamplePet1 = model.Pet{
		Id:    uuid.NewString(),
		Name:  "pet 1",
		Age:   10,
		Owner: "User1",
	}
	SamplePet2 = model.Pet{
		Id:    uuid.NewString(),
		Name:  "pet 2",
		Age:   92,
		Owner: "User2",
	}
	SamplePet1Item = data.DynamoItem{
		"Id":    {S: &SamplePet1.Id},
		"Name":  {S: &SamplePet1.Name},
		"Age":   {N: jsii.String("10")},
		"Owner": {S: &SamplePet1.Owner},
	}
	SamplePet2Item = data.DynamoItem{
		"Id":    {S: &SamplePet2.Id},
		"Name":  {S: &SamplePet2.Name},
		"Age":   {N: jsii.String("92")},
		"Owner": {S: &SamplePet2.Owner},
	}
)

func TestPetDelete(t *testing.T) {
	// Define test struct
	type Test struct {
		name      string
		dbClient  FakeDynamoDbClient
		petId     string
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid delete",
			dbClient:  FakeDynamoDbClient{},
			petId:     SamplePet1.Id,
			expectErr: false,
		},
		{
			name: "db delete error",
			dbClient: FakeDynamoDbClient{
				deleteItemErr: assert.AnError,
			},
			petId:     SamplePet1.Id,
			expectErr: true,
		},
		{
			name: "db item not found error",
			dbClient: FakeDynamoDbClient{
				deleteItemErr: &dynamodb.ConditionalCheckFailedException{},
			},
			petId:     SamplePet1.Id,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPetDao(&test.dbClient, SampleTableName)

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

func TestPetGetById(t *testing.T) {
	// Define test struct
	type Test struct {
		name        string
		dbClient    FakeDynamoDbClient
		petId       string
		expectedPet model.Pet
		expectErr   bool
	}

	// Define tests
	tests := []Test{
		{
			name: "pet found",
			dbClient: FakeDynamoDbClient{
				getItemOutput: &dynamodb.GetItemOutput{Item: SamplePet1Item},
			},
			petId:       SamplePet1.Id,
			expectedPet: SamplePet1,
			expectErr:   false,
		},
		{
			name: "pet not found",
			dbClient: FakeDynamoDbClient{
				getItemOutput: &dynamodb.GetItemOutput{Item: nil},
			},
			petId:     SamplePet1.Id,
			expectErr: true,
		},
		{
			name: "db get error",
			dbClient: FakeDynamoDbClient{
				getItemErr: assert.AnError,
			},
			petId:     SamplePet1.Id,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPetDao(&test.dbClient, SampleTableName)

		// Execute
		pet, err := dao.GetById(test.petId)

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
			assert.Equal(t, test.expectedPet, pet, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestPetInsert(t *testing.T) {
	// Define test struct
	type Test struct {
		name      string
		dbClient  FakeDynamoDbClient
		pet       model.Pet
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid insert",
			dbClient:  FakeDynamoDbClient{},
			pet:       SamplePet1,
			expectErr: false,
		},
		{
			name: "db put error",
			dbClient: FakeDynamoDbClient{
				putItemErr: assert.AnError,
			},
			pet:       SamplePet1,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPetDao(&test.dbClient, SampleTableName)

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

func TestPetQuery(t *testing.T) {
	// Define test struct
	type Test struct {
		name                string
		dbClient            FakeDynamoDbClient
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
			dbClient: FakeDynamoDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            pointy.Int64(2),
					Items:            []data.DynamoItem{SamplePet1Item, SamplePet2Item},
					LastEvaluatedKey: data.DynamoItem{},
				}},
			count:            3,
			exclusiveStartId: "",
			expectedPets: []model.Pet{
				SamplePet1,
				SamplePet2,
			},
			expectedHasNextPage: false,
		},
		{
			name: "request less items than in DB (start at beginning of list)",
			dbClient: FakeDynamoDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            pointy.Int64(2),
					Items:            []data.DynamoItem{SamplePet1Item},
					LastEvaluatedKey: SamplePet1Item,
				}},
			count: 1,
			expectedPets: []model.Pet{
				SamplePet1,
			},
			expectedHasNextPage: true,
		},
		{
			name: "request less items than in DB (reach end of list)",
			dbClient: FakeDynamoDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            pointy.Int64(2),
					Items:            []data.DynamoItem{SamplePet2Item},
					LastEvaluatedKey: data.DynamoItem{},
				}},
			count:            1,
			exclusiveStartId: SamplePet1.Id,
			expectedPets: []model.Pet{
				SamplePet2,
			},
			expectedHasNextPage: false,
		},
		{
			name: "request 'count=0' but 'exclusiveStartId' is not last pet",
			dbClient: FakeDynamoDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            pointy.Int64(2),
					Items:            []data.DynamoItem{},
					LastEvaluatedKey: SamplePet2Item,
				}},
			count:               0,
			exclusiveStartId:    SamplePet1.Id,
			expectedPets:        []model.Pet{},
			expectedHasNextPage: true,
		},
		{
			name: "request 'count=0' but 'exclusiveStartId' is last pet",
			dbClient: FakeDynamoDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            pointy.Int64(2),
					Items:            []data.DynamoItem{},
					LastEvaluatedKey: data.DynamoItem{},
				}},
			count:               0,
			exclusiveStartId:    SamplePet2.Id,
			expectedPets:        []model.Pet{},
			expectedHasNextPage: false,
		},
		{
			name: "no pets returned when 'count=0'",
			dbClient: FakeDynamoDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count:            pointy.Int64(1),
					Items:            []data.DynamoItem{SamplePet1Item},
					LastEvaluatedKey: SamplePet1Item,
				},
			},
			expectedPets:        []model.Pet{},
			expectedHasNextPage: true,
		},
		{
			name: "db query error",
			dbClient: FakeDynamoDbClient{
				queryErr: assert.AnError,
			},
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPetDao(&test.dbClient, SampleTableName)

		// Execute
		pets, hasNextPage, err := dao.Query(test.count, test.exclusiveStartId)

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
			assert.Equal(t, test.expectedPets, pets, test.name)
			assert.Equal(t, test.expectedHasNextPage, hasNextPage, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestPetGetTotalCount(t *testing.T) {
	// Define test struct
	type Test struct {
		name          string
		dbClient      FakeDynamoDbClient
		expectedCount int
		expectErr     bool
	}

	// Define tests
	tests := []Test{
		{
			name: "valid get total count",
			dbClient: FakeDynamoDbClient{
				queryOutput: &dynamodb.QueryOutput{
					Count: pointy.Int64(10),
				},
			},
			expectedCount: 10,
		},
		{
			name: "db query error",
			dbClient: FakeDynamoDbClient{
				queryErr: assert.AnError,
			},
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPetDao(&test.dbClient, SampleTableName)

		// Execute
		count, err := dao.GetTotalCount()

		// Verify
		if !test.expectErr {
			assert.Nil(t, err, test.name)
			assert.Equal(t, test.expectedCount, count, test.name)
		} else {
			assert.NotNil(t, err, test.name)
		}
	}
}

func TestPetUpdate(t *testing.T) {
	// Define test struct
	type Test struct {
		name      string
		dbClient  FakeDynamoDbClient
		pet       model.Pet
		expectErr bool
	}

	// Define tests
	tests := []Test{
		{
			name:      "valid update",
			dbClient:  FakeDynamoDbClient{},
			pet:       SamplePet1,
			expectErr: false,
		},
		{
			name: "db put error",
			dbClient: FakeDynamoDbClient{
				putItemErr: assert.AnError,
			},
			pet:       SamplePet1,
			expectErr: true,
		},
	}

	// Run tests
	for _, test := range tests {
		// Setup
		dao := data.NewPetDao(&test.dbClient, SampleTableName)

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
