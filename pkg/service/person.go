package service

import (
	"github.com/mcwiet/go-test/pkg/data"
	"github.com/mcwiet/go-test/pkg/model"
)

func GetPerson(id string) (model.Person, error) {
	return data.GetPerson(id)
}

func GetPeople() ([]model.Person, error) {
	return data.GetPeople()
}
