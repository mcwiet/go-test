package model

type Pet struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Owner string `json:"owner,omitempty"`
}

type CreatePetInput struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Owner string `json:"owner,omitempty"`
}

type CreatePetPayload struct {
	Pet Pet `json:"pet"`
}

type DeletePetInput struct {
	Id string `json:"id"`
}

type DeletePetPayload struct {
	Id string `json:"id"`
}

type PetInput struct {
	Id string `json:"id"`
}

type PetsInput struct {
	First int    `json:"first"`
	After string `json:"after"`
}

type UpdatePetInput struct {
	Name  string `json:"name,omitempty"`
	Age   int    `json:"age,omitempty"`
	Owner string `json:"owner,omitempty"`
}

type UpdatePetPayload struct {
	Pet Pet `json:"pet"`
}

type PetEdge struct {
	Node   Pet    `json:"node"`
	Cursor string `json:"cursor"`
}

type PetConnection struct {
	TotalCount int       `json:"totalCount"`
	Edges      []PetEdge `json:"edges,omitempty"`
	PageInfo   PageInfo  `json:"pageInfo"`
}
