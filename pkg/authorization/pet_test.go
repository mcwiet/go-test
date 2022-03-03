package authorization_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mcwiet/go-test/pkg/authorization"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/mcwiet/go-test/pkg/service"
	"github.com/stretchr/testify/assert"
)

var (
	SamplePet = model.Pet{
		Id:    uuid.NewString(),
		Name:  "Levi",
		Age:   1,
		Owner: SampleUsername,
	}
)

func TestIsAuthorized(t *testing.T) {
	type Test struct {
		name           string
		identity       model.Identity
		pet            model.Pet
		action         service.PetAction
		expectedResult bool
	}

	tests := []Test{
		{
			name: "update pet owner - not authorized",
			identity: model.Identity{
				Username: "unexpected",
			},
			pet:            SamplePet,
			action:         service.PetActionUpdateOwner,
			expectedResult: false,
		},
		{
			name: "update pet owner - admin",
			identity: model.Identity{
				Username: "unexpected",
				Groups:   map[string]bool{authorization.RoleAdmin.String(): true},
			},
			pet:            SamplePet,
			action:         service.PetActionUpdateOwner,
			expectedResult: true,
		},
		{
			name: "update pet owner - user is owner",
			identity: model.Identity{
				Username: SamplePet.Owner,
			},
			pet:            SamplePet,
			action:         service.PetActionUpdateOwner,
			expectedResult: true,
		},
		{
			name: "undefined action",
			identity: model.Identity{
				Username: SamplePet.Owner,
			},
			pet:            SamplePet,
			action:         service.PetActionUndefined,
			expectedResult: false,
		},
	}

	for _, test := range tests {
		authorizer := authorization.NewPetAuthorizer()

		authorized := authorizer.IsAuthorized(test.identity, test.pet, test.action)

		assert.Equal(t, test.expectedResult, authorized, test.name)
	}
}
