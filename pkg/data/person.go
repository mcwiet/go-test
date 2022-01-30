package data

import (
	"errors"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/jsii-runtime-go"
	"github.com/mcwiet/go-test/pkg/model"
)

type PersonDao struct {
	client    *dynamodb.DynamoDB
	tableName string
}

func NewPersonDao(client *dynamodb.DynamoDB, tableName string) PersonDao {
	return PersonDao{
		client:    client,
		tableName: tableName,
	}
}

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
		return nil, errors.New("internal error retrieving person")
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

func (p *PersonDao) GetPeople() (*[]model.Person, error) {
	ret, err := p.client.Query(&dynamodb.QueryInput{
		TableName: &p.tableName,
		// TODO index name
		IndexName:              jsii.String("sort-key-gsi"),
		KeyConditionExpression: jsii.String("Sort = :sortVal"),
		ProjectionExpression:   jsii.String("#name, Age"),
		ExpressionAttributeNames: map[string]*string{
			"#name": jsii.String("Name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":sortVal": {S: jsii.String("person")},
		},
	})

	if err != nil {
		log.Println(err)
		return nil, errors.New("could not retrieve people")
	}

	items := ret.Items
	arr := []model.Person{}
	for _, item := range items {
		age, _ := strconv.Atoi(*item["Age"].N)
		arr = append(arr, model.Person{
			Name: *item["Name"].S,
			Age:  age,
		})
	}
	return &arr, nil
}
