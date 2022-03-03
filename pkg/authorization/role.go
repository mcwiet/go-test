package authorization

type Role int

const (
	RoleUndefined Role = iota
	RoleAdmin
)

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	}
	return "unknown"
}
