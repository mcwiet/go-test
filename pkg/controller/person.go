package controller

import (
	"github.com/mcwiet/go-test/pkg/service"
)

func GetPerson(request Request) Response {
	person, err := service.GetPerson(request.Arguments["id"])
	return Response{
		Data:  person,
		Error: err,
	}
}

func GetPeople(request Request) Response {
	people, err := service.GetPeople()
	return Response{
		Data:  people,
		Error: err,
	}
}
