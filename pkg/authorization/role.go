package authorization

type Role int

const (
	Undefined Role = iota
	Admin
)

func (r Role) String() string {
	switch r {
	case Admin:
		return "admin"
	}
	return "unknown"
}
