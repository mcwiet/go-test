package data

import (
	"errors"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/jsii-runtime-go"
	"github.com/mcwiet/go-test/pkg/model"
)

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
		Key: map[string]*dynamodb.AttributeValue{
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
func (p *PersonDao) GetById(id string) (*model.Person, error) {
	ret, err := p.client.GetItem(&dynamodb.GetItemInput{
		TableName: &p.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id":   {S: jsii.String(id)},
			"Sort": {S: jsii.String(personSortLabel)},
		},
		ProjectionExpression: jsii.String("#name, Age"),
		ExpressionAttributeNames: map[string]*string{
			"#name": jsii.String("Name"),
		},
	})

	if err != nil {
		log.Println(err)
		return nil, errors.New("error retrieving person")
	} else if ret == nil || ret.Item == nil {
		return nil, errors.New("person not found")
	}

	item := ret.Item
	age, err := strconv.Atoi(*item["Age"].N)
	person := model.Person{
		Id:   id,
		Name: *item["Name"].S,
		Age:  age,
	}

	return &person, err
}

// Inserts a person to the data store
func (p *PersonDao) Insert(person *model.Person) error {
	age := strconv.Itoa(person.Age)
	_, err := p.client.PutItem(&dynamodb.PutItemInput{
		TableName: &p.tableName,
		Item: map[string]*dynamodb.AttributeValue{
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

// Lists people from the data store
func (p *PersonDao) List(first int, after string) (model.PersonConnection, error) {
	queryRet, err := p.queryPeople(first, after)

	if err != nil {
		log.Println(err)
		return model.PersonConnection{}, errors.New("error retrieving people")
	}

	totalCount, err := p.getTotalCount()

	if err != nil {
		log.Println(err)
		return model.PersonConnection{}, errors.New("error getting total people count")
	}

	connection := model.PersonConnection{
		TotalCount: totalCount,
		Edges:      []model.PersonEdge{},
	}

	for _, item := range queryRet.Items {
		age, _ := strconv.Atoi(*item["Age"].N)
		connection.Edges = append(connection.Edges, model.PersonEdge{
			Node: model.Person{
				Id:   *item["Id"].S,
				Name: *item["Name"].S,
				Age:  age,
			},
			Cursor: *item["Id"].S,
		})
	}

	endCursor := ""
	if len(queryRet.Items) > 0 {
		endCursor = *queryRet.Items[len(queryRet.Items)-1]["Id"].S
	}

	connection.PageInfo = model.PageInfo{
		EndCursor:   endCursor,
		HasNextPage: len(queryRet.LastEvaluatedKey) != 0,
	}

	return connection, nil
}

// Get the total count of people
func (p *PersonDao) getTotalCount() (int, error) {
	ret, err := p.client.Query(&dynamodb.QueryInput{
		TableName:              &p.tableName,
		IndexName:              jsii.String("sort-key-gsi"),
		KeyConditionExpression: jsii.String("Sort = :sortVal"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":sortVal": {S: jsii.String(personSortLabel)},
		},
	})

	return int(*ret.Count), err
}

// Query for a set of people (first n people after the exclusive start value)
func (p *PersonDao) queryPeople(count int, exclusiveStartId string) (dynamodb.QueryOutput, error) {
	if count < 1 {
		return dynamodb.QueryOutput{}, nil
	}

	limit := int64(count)
	exclusiveStartKey := map[string]*dynamodb.AttributeValue{
		"Id":   {S: &exclusiveStartId},
		"Sort": {S: jsii.String(personSortLabel)},
	}
	if exclusiveStartId == "" {
		exclusiveStartKey = nil
	}

	ret, err := p.client.Query(&dynamodb.QueryInput{
		TableName:              &p.tableName,
		IndexName:              jsii.String("sort-key-gsi"),
		KeyConditionExpression: jsii.String("Sort = :sortVal"),
		ProjectionExpression:   jsii.String("Id, #name, Age"),
		ExpressionAttributeNames: map[string]*string{
			"#name": jsii.String("Name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":sortVal": {S: jsii.String(personSortLabel)},
		},
		ExclusiveStartKey: exclusiveStartKey,
		Limit:             &limit,
	})

	return *ret, err
}
