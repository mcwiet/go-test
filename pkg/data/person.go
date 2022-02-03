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
	personIdPrefix = "person-"
	personSort     = "person"
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
			"Id":   {S: jsii.String(personIdPrefix + id)},
			"Sort": {S: jsii.String(personSort)},
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
			"Id":   {S: jsii.String(personIdPrefix + id)},
			"Sort": {S: jsii.String(personSort)},
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
			"Id":   {S: jsii.String(personIdPrefix + person.Id)},
			"Sort": {S: jsii.String(personSort)},
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

// Lists persons from the data store
func (p *PersonDao) List() (*[]model.Person, error) {
	ret, err := p.client.Query(&dynamodb.QueryInput{
		TableName:              &p.tableName,
		IndexName:              jsii.String("sort-key-gsi"),
		KeyConditionExpression: jsii.String("Sort = :sortVal"),
		ProjectionExpression:   jsii.String("Id, #name, Age"),
		ExpressionAttributeNames: map[string]*string{
			"#name": jsii.String("Name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":sortVal": {S: jsii.String(personSort)},
		},
	})

	if ret == nil || err != nil {
		log.Println(err)
		return nil, errors.New("error retrieving people")
	}

	items := ret.Items
	arr := []model.Person{}
	for _, item := range items {
		age, _ := strconv.Atoi(*item["Age"].N)
		arr = append(arr, model.Person{
			Id:   *item["Id"].S,
			Name: *item["Name"].S,
			Age:  age,
		})
	}
	return &arr, nil
}
