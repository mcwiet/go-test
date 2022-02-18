package model

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `name:"name,omitempty"`
}

type UserEdge struct {
	Node User `json:"node"`
}

type UserConnection struct {
	TotalCount int        `json:"totalCount"`
	Edges      []UserEdge `json:"edges"`
	PageInfo   PageInfo   `json:"pageInfo"`
}

type UserInput struct {
	Username string `json:"username"`
}

type UsersInput struct {
	First int    `json:"first"`
	After string `json:"after"`
}
