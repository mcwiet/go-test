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

// Object containing information needed to access the person data store
type PersonDao struct {
	client    DynamoDbClient
	tableName string
}

var (
	personSortLabel = "person"
)

// Creates a person data store access object
func NewPersonDao(client DynamoDbClient, tableName string) PersonDao {
	return PersonDao{
		client:    client,
		tableName: tableName,
	}
}

// Deletes a person from the data store
func (p *PersonDao) Delete(id string) error {
	_, err := p.client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &p.tableName,
		Key: DynamoItem{
			"Id":   {S: jsii.String(id)},
			"Sort": {S: jsii.String(personSortLabel)},
		},
		ConditionExpression: jsii.String("attribute_exists(Id)"),
	})

	if err != nil {
		log.Println(err)
		var notFoundError *dynamodb.ConditionalCheckFailedException
		if errors.As(err, &notFoundError) {
			return errors.New("could not delete person; person not found")
		} else {
			return errors.New("error deleting person")
		}
	}

	return nil
}

// Gets a person from the data store using the ID
func (p *PersonDao) GetById(id string) (model.Person, error) {
	ret, err := p.client.GetItem(&dynamodb.GetItemInput{
		TableName: &p.tableName,
		Key: DynamoItem{
			"Id":   {S: jsii.String(id)},
			"Sort": {S: jsii.String(personSortLabel)},
		},
		ProjectionExpression: jsii.String("Id, #name, Age"),
		ExpressionAttributeNames: map[string]*string{
			"#name": jsii.String("Name"),
		},
	})

	if err != nil {
		log.Println(err)
		return model.Person{}, errors.New("error retrieving person")
	} else if ret == nil || ret.Item == nil {
		return model.Person{}, errors.New("person not found")
	}

	person := convertItemToPerson(ret.Item)

	return person, err
}

// Inserts a person to the data store
func (p *PersonDao) Insert(person model.Person) error {
	age := strconv.Itoa(person.Age)
	_, err := p.client.PutItem(&dynamodb.PutItemInput{
		TableName: &p.tableName,
		Item: DynamoItem{
			"Id":   {S: jsii.String(person.Id)},
			"Sort": {S: jsii.String(personSortLabel)},
			"Name": {S: &person.Name},
			"Age":  {N: &age},
		},
	})

	if err != nil {
		log.Println(err)
		return errors.New("error adding person")
	}

	return nil
}

// Query for a set of people (first n people after the exclusive start value)
func (p *PersonDao) Query(count int, exclusiveStartId string) ([]model.Person, bool, error) {
	queryInput := buildQueryInput(p.tableName, count, exclusiveStartId)

	ret, err := p.client.Query(&queryInput)

	if err != nil {
		log.Println(err)
		return []model.Person{}, false, errors.New("error retrieving people")
	}

	hasNextPage := len(ret.LastEvaluatedKey) != 0

	// If count is zero, double check value for hasNextPage and ensure un-requested items are not returned
	if count == 0 {
		if len(ret.Items) > 0 {
			hasNextPage = true
		}
		ret.Items = []DynamoItem{}
	}

	// Convert items to people
	people := []model.Person{}
	for _, item := range ret.Items {
		people = append(people, convertItemToPerson(item))
	}

	return people, hasNextPage, err
}

// Get the total count of people
func (p *PersonDao) GetTotalCount() (int, error) {
	ret, err := p.client.Query(&dynamodb.QueryInput{
		TableName:              &p.tableName,
		IndexName:              jsii.String("sort-key-gsi"),
		KeyConditionExpression: jsii.String("Sort = :sortVal"),
		ExpressionAttributeValues: DynamoItem{
			":sortVal": {S: jsii.String(personSortLabel)},
		},
	})

	if err != nil {
		log.Println(err)
		return 0, errors.New("error getting total people count")
	}

	count := int(*ret.Count)

	return count, err
}

// Convert a DynamoDB item to a person
func convertItemToPerson(item DynamoItem) model.Person {
	age, _ := strconv.Atoi(*item["Age"].N)
	return model.Person{
		Id:   *item["Id"].S,
		Name: *item["Name"].S,
		Age:  age,
	}
}

// Build input to query for people
func buildQueryInput(tableName string, count int, exclusiveStartId string) dynamodb.QueryInput {
	limit := int64(count)
	if count == 0 {
		// Dynamo minimum limit is 1
		limit = 1
	}
	exclusiveStartKey := DynamoItem{
		"Id":   {S: &exclusiveStartId},
		"Sort": {S: jsii.String(personSortLabel)},
	}
	if exclusiveStartId == "" {
		exclusiveStartKey = nil
	}

	return dynamodb.QueryInput{
		TableName:              &tableName,
		IndexName:              jsii.String("sort-key-gsi"),
		KeyConditionExpression: jsii.String("Sort = :sortVal"),
		ProjectionExpression:   jsii.String("Id, #name, Age"),
		ExpressionAttributeNames: map[string]*string{
			"#name": jsii.String("Name"),
		},
		ExpressionAttributeValues: DynamoItem{
			":sortVal": {S: jsii.String(personSortLabel)},
		},
		ExclusiveStartKey: exclusiveStartKey,
		Limit:             &limit,
	}
}
