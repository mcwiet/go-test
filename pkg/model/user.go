package model

type User struct {
	Username string `json:"id"`
	Email    string `json:"email"`
	Name     string `name:"name,omitempty"`
}

type UserEdge struct {
	Node   User   `json:"user"`
	Cursor string `json:"cursor"`
}

type UserConnection struct {
	TotalCount int        `json:"totalCount"`
	Edges      []UserEdge `json:"edges"`
	PageInfo   PageInfo   `json:"pageInfo"`
}
