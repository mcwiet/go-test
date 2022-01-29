package data

import "github.com/mcwiet/go-test/pkg/model"

var (
	people = map[string]model.Person{
		"6e483041-002a-4942-bc18-5605e5826078": {
			Id:   "6e483041-002a-4942-bc18-5605e5826078",
			Name: "Mike",
			Age:  28,
		},
		"0231d150-fee7-4158-a0c2-95831b152062": {
			Id:   "0231d150-fee7-4158-a0c2-95831b152062",
			Name: "Katherine",
			Age:  28,
		},
		"43186a74-db2f-463c-a1d9-d4adea731cfc": {
			Id:   "43186a74-db2f-463c-a1d9-d4adea731cfc",
			Name: "Levi",
			Age:  1,
		},
	}
)

func GetPerson(id string) model.Person {
	return people[id]
}

func GetPeople() *[]model.Person {
	arr := []model.Person{}
	for _, value := range people {
		arr = append(arr, value)
	}
	return &arr
}
