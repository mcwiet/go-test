package model

type Pet struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Owner string `json:"owner,omitempty"`
}

type PetEdge struct {
	Node   Pet    `json:"node"`
	Cursor string `json:"cursor"`
}

type PetConnection struct {
	TotalCount int       `json:"totalCount"`
	Edges      []PetEdge `json:"edges"`
	PageInfo   PageInfo  `json:"pageInfo"`
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

type UpdatePetOwnerInput struct {
	Id    string `json:"id"`
	Owner string `json:"owner"`
}

type UpdatePetOwnerPayload struct {
	Pet Pet `json:"pet"`
}
