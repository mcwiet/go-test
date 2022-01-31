package data

import (
	"errors"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/jsii-runtime-go"
	"github.com/mcwiet/go-test/pkg/model"
)

// Object containing information needed to access the Person data store
type PersonDao struct {
	client    *dynamodb.DynamoDB
	tableName string
}

// Creates a Person data store access object
func NewPersonDao(client *dynamodb.DynamoDB, tableName string) PersonDao {
	return PersonDao{
		client:    client,
		tableName: tableName,
	}
}

// Adds a person to the data store
func (p *PersonDao) AddPerson(person *model.Person) error {
	age := strconv.Itoa(person.Age)
	_, err := p.client.PutItem(&dynamodb.PutItemInput{
		TableName: &p.tableName,
		Item: map[string]*dynamodb.AttributeValue{
			"Id":   {S: jsii.String("person-" + person.Id)},
			"Sort": {S: jsii.String("person")},
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

// Gets a person from the data store using the ID
func (p *PersonDao) GetPerson(id string) (*model.Person, error) {
	ret, err := p.client.GetItem(&dynamodb.GetItemInput{
		TableName: &p.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Id":   {S: jsii.String("person-" + id)},
			"Sort": {S: jsii.String("person")},
		},
		ProjectionExpression: jsii.String("#name, Age"),
		ExpressionAttributeNames: map[string]*string{
			"#name": jsii.String("Name"),
		},
	})

	if err != nil {
		log.Println(err)
		return nil, errors.New("error retrieving person")
	} else if ret.Item == nil {
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

// Get a list of persons from the data store
func (p *PersonDao) GetPeople() (*[]model.Person, error) {
	ret, err := p.client.Query(&dynamodb.QueryInput{
		TableName:              &p.tableName,
		IndexName:              jsii.String("sort-key-gsi"),
		KeyConditionExpression: jsii.String("Sort = :sortVal"),
		ProjectionExpression:   jsii.String("Id, #name, Age"),
		ExpressionAttributeNames: map[string]*string{
			"#name": jsii.String("Name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":sortVal": {S: jsii.String("person")},
		},
	})

	if err != nil {
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
