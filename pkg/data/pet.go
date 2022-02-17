package data

import (
	"errors"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/jsii-runtime-go"
	"github.com/mcwiet/go-test/pkg/model"
)

type DynamoItem = map[string]*dynamodb.AttributeValue

type DynamoDbClient interface {
	DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
	GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	Query(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
}

// Object containing information needed to access the pet data store
type PetDao struct {
	client    DynamoDbClient
	tableName string
}

const (
	petSortLabel = "pet"
)

// Creates a pet data store access object
func NewPetDao(client DynamoDbClient, tableName string) PetDao {
	return PetDao{
		client:    client,
		tableName: tableName,
	}
}

// Deletes a pet from the data store
func (p *PetDao) Delete(id string) error {
	_, err := p.client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &p.tableName,
		Key: DynamoItem{
			"Id":   {S: jsii.String(id)},
			"Sort": {S: jsii.String(petSortLabel)},
		},
		ConditionExpression: jsii.String("attribute_exists(Id)"),
	})

	if err != nil {
		log.Println(err)
		var notFoundError *dynamodb.ConditionalCheckFailedException
		if errors.As(err, &notFoundError) {
			return errors.New("could not delete pet; pet not found")
		} else {
			return errors.New("error deleting pet")
		}
	}

	return nil
}

// Gets a pet from the data store using the ID
func (p *PetDao) GetById(id string) (model.Pet, error) {
	ret, err := p.client.GetItem(&dynamodb.GetItemInput{
		TableName: &p.tableName,
		Key: DynamoItem{
			"Id":   {S: jsii.String(id)},
			"Sort": {S: jsii.String(petSortLabel)},
		},
		ProjectionExpression: jsii.String("Id, #name, Age, #owner"),
		ExpressionAttributeNames: map[string]*string{
			"#name":  jsii.String("Name"),
			"#owner": jsii.String("Owner"),
		},
	})

	if err != nil {
		log.Println(err)
		return model.Pet{}, errors.New("error retrieving pet")
	} else if ret == nil || ret.Item == nil {
		return model.Pet{}, errors.New("pet not found")
	}

	pet := convertItemToPet(ret.Item)

	return pet, err
}

// Inserts a pet to the data store
func (p *PetDao) Insert(pet model.Pet) error {
	_, err := p.client.PutItem(&dynamodb.PutItemInput{
		TableName: &p.tableName,
		Item:      convertPetToItem(pet),
	})

	if err != nil {
		log.Println(err)
		return errors.New("error adding pet")
	}

	return nil
}

// Query for a set of pets (first n pets after the exclusive start value)
func (p *PetDao) Query(count int, exclusiveStartId string) ([]model.Pet, bool, error) {
	queryInput := buildQueryInput(p.tableName, count, exclusiveStartId)

	ret, err := p.client.Query(&queryInput)

	if err != nil {
		log.Println(err)
		return []model.Pet{}, false, errors.New("error retrieving pets")
	}

	hasNextPage := len(ret.LastEvaluatedKey) != 0

	// If count is zero, double check value for hasNextPage and ensure un-requested items are not returned
	if count == 0 {
		if len(ret.Items) > 0 {
			hasNextPage = true
		}
		ret.Items = []DynamoItem{}
	}

	// Convert items to pets
	pets := []model.Pet{}
	for _, item := range ret.Items {
		pets = append(pets, convertItemToPet(item))
	}

	return pets, hasNextPage, err
}

// Get the total count of pets
func (p *PetDao) GetTotalCount() (int, error) {
	ret, err := p.client.Query(&dynamodb.QueryInput{
		TableName:              &p.tableName,
		IndexName:              jsii.String("sort-key-gsi"),
		KeyConditionExpression: jsii.String("Sort = :sortVal"),
		ExpressionAttributeValues: DynamoItem{
			":sortVal": {S: jsii.String(petSortLabel)},
		},
	})

	if err != nil {
		log.Println(err)
		return 0, errors.New("error getting total pets count")
	}

	count := int(*ret.Count)

	return count, err
}

// Updates a pet in the data store by performing a full replace
func (p *PetDao) Update(pet model.Pet) error {
	_, err := p.client.PutItem(&dynamodb.PutItemInput{
		TableName: &p.tableName,
		Item:      convertPetToItem(pet),
	})

	if err != nil {
		log.Println(err)
		return errors.New("error updating pet")
	}

	return nil
}

// Convert a DynamoDB item to a pet
func convertItemToPet(item DynamoItem) model.Pet {
	age, _ := strconv.Atoi(*item["Age"].N)
	owner := ""
	if item["Owner"] != nil {
		owner = *item["Owner"].S
	}
	return model.Pet{
		Id:    *item["Id"].S,
		Name:  *item["Name"].S,
		Age:   age,
		Owner: owner,
	}
}

// Convert a pet to a DynamoDB item
func convertPetToItem(pet model.Pet) DynamoItem {
	age := strconv.Itoa(pet.Age)
	return DynamoItem{
		"Id":    {S: jsii.String(pet.Id)},
		"Sort":  {S: jsii.String(petSortLabel)},
		"Name":  {S: &pet.Name},
		"Age":   {N: &age},
		"Owner": {S: &pet.Owner},
	}
}

// Build input to query for pets
func buildQueryInput(tableName string, count int, exclusiveStartId string) dynamodb.QueryInput {
	limit := int64(count)
	if count == 0 {
		// Dynamo minimum limit is 1
		limit = 1
	}
	exclusiveStartKey := DynamoItem{
		"Id":   {S: &exclusiveStartId},
		"Sort": {S: jsii.String(petSortLabel)},
	}
	if exclusiveStartId == "" {
		exclusiveStartKey = nil
	}

	return dynamodb.QueryInput{
		TableName:              &tableName,
		IndexName:              jsii.String("sort-key-gsi"),
		KeyConditionExpression: jsii.String("Sort = :sortVal"),
		ProjectionExpression:   jsii.String("Id, #name, Age, #owner"),
		ExpressionAttributeNames: map[string]*string{
			"#name":  jsii.String("Name"),
			"#owner": jsii.String("Owner"),
		},
		ExpressionAttributeValues: DynamoItem{
			":sortVal": {S: jsii.String(petSortLabel)},
		},
		ExclusiveStartKey: exclusiveStartKey,
		Limit:             &limit,
	}
}
