package authorization

import (
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/mcwiet/go-test/pkg/service"
)

type PetAuthorizer struct{}

func NewPetAuthorizer() PetAuthorizer {
	return PetAuthorizer{}
}

func (a *PetAuthorizer) IsAuthorized(identity model.Identity, pet model.Pet, action service.PetAction) bool {
	if identity.Groups[RoleAdmin.String()] {
		return true
	}

	switch action {
	case service.PetActionUpdateOwner:
		return canUpdatePetOwner(identity, pet)
	default:
		return false
	}
}

func canUpdatePetOwner(identity model.Identity, pet model.Pet) bool {
	return identity.Username == pet.Owner
}
