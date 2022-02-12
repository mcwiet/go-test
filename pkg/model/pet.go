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
	Edges      []PetEdge `json:"edges,omitempty"`
	PageInfo   PageInfo  `json:"pageInfo"`
}
