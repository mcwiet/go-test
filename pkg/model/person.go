package model

type Person struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type PersonEdge struct {
	Node   Person `json:"Node"`
	Cursor string `json:"Cursor"`
}

type PersonConnection struct {
	TotalCount int          `json:"totalCount"`
	Edges      []PersonEdge `json:"edges"`
	PageInfo   PageInfo     `json:"pageInfo"`
}
