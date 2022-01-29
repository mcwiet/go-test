package service

import (
	"github.com/mcwiet/go-test/pkg/data"
	"github.com/mcwiet/go-test/pkg/model"
)

func GetPerson(id string) model.Person {
	return data.GetPerson(id)
}

func GetPeople() *[]model.Person {
	return data.GetPeople()
}
