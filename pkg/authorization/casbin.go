package authorization

import "errors"

type CasbinEnforcer interface {
	AddPolicy(...interface{}) (bool, error)
	AddRoleForUser(user string, role string, domain ...string) (bool, error)
	DeleteRoleForUser(user string, role string, domain ...string) (bool, error)
	Enforce(...interface{}) (bool, error)
	RemovePolicy(...interface{}) (bool, error)
}

type CasbinAuthorizer struct {
	enforcer CasbinEnforcer
}

func NewCasbinAuthorizer(enforcer CasbinEnforcer) CasbinAuthorizer {
	return CasbinAuthorizer{
		enforcer: enforcer,
	}
}

func (a *CasbinAuthorizer) AddPermission(subject string, object string, action string) (bool, error) {
	return a.enforcer.AddPolicy(subject, object, action)
}

func (a *CasbinAuthorizer) AddRoleForUser(user string, role Role) (bool, error) {
	if role == Undefined {
		return false, errors.New("cannot add undefined role to user")
	}
	return a.enforcer.AddRoleForUser(user, role.String())
}

func (a *CasbinAuthorizer) IsAuthorized(subject string, object string, action string) (bool, error) {
	return a.enforcer.Enforce(subject, object, action)
}

func (a *CasbinAuthorizer) RemovePermission(subject string, object string, action string) (bool, error) {
	return a.enforcer.RemovePolicy(subject, object, action)
}

func (a *CasbinAuthorizer) RemoveRoleForUser(user string, role Role) (bool, error) {
	if role == Undefined {
		return false, errors.New("cannot remove undefined role from user")
	}
	return a.enforcer.DeleteRoleForUser(user, role.String())
}
