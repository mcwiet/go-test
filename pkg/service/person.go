package service

import (
	"cwietnie/go-test/pkg/data"
	"cwietnie/go-test/pkg/model"
)

func GetPerson(id string) model.Person {
	return data.GetPerson(id)
}

func GetPeople() *[]model.Person {
	return data.GetPeople()
}
