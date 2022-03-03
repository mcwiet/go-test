package model

type Identity struct {
	Username string
	Email    string
	Groups   map[string]bool
}
