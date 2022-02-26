package integration_test

import (
	"context"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/mcwiet/go-test/pkg/model"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

// Sequentially run functions involved for testing pet API operations
func TestPetApi(t *testing.T) {
	// Create some pets
	pet1 := createPet(t)
	pet2 := createPet(t)

	// List the pets
	listPets(t)

	// Get a pet
	getPet(t, pet1.Id, &pet1)

	// Update the pet owner
	updatePetOwner(t, pet1, UserToken.Username)

	// Delete the pets
	deletePet(t, &pet1)
	deletePet(t, &pet2)

	// Attempt to get a pet again
	getPet(t, pet1.Id, nil)
}

func createPet(t *testing.T) model.Pet {
	// Setup
	petName := "Integration Test"
	petAge := 10
	petOwner := ""
	request := graphql.NewRequest(`
		mutation ($name: String!, $age: Int!, $owner: String) {
			createPet (input: { name: $name, age: $age, owner: $owner }) {
				pet {
					id
					name
					age
					owner
				}
			}
		}
	`)
	request.Var("name", petName)
	request.Var("age", petAge)
	request.Var("owner", petOwner)
	request.Header.Set("Authorization", UserToken.AccessTokenString)

	// Execute
	var response map[string]interface{}
	err := GraphQlClient.Run(context.Background(), request, &response)
	var payload model.CreatePetPayload
	mapstructure.Decode(response["createPet"], &payload)

	// Verify
	stepName := "createPet"
	assert.Nil(t, err, stepName+": should not error")
	if err != nil {
		return model.Pet{}
	}
	pet := payload.Pet
	assert.NotNil(t, pet.Id, stepName+"id should exist")
	assert.Equal(t, petName, pet.Name, stepName+"name should match")
	assert.Equal(t, petAge, pet.Age, stepName+"age should match")
	assert.Equal(t, petOwner, pet.Owner, stepName+"owner should match")

	return pet
}

func deletePet(t *testing.T, pet *model.Pet) {
	// Setup
	request := graphql.NewRequest(`
		mutation ($id: ID!) {
			deletePet (input: { id: $id }) {
				id
			}
		}
	`)
	request.Var("id", pet.Id)
	request.Header.Set("Authorization", UserToken.AccessTokenString)

	// Execute
	var response map[string]interface{}
	err := GraphQlClient.Run(context.Background(), request, &response)

	// Verify
	stepName := "deletePet"
	assert.Nil(t, err, stepName+": should not error")
}

func getPet(t *testing.T, id string, expectedPet *model.Pet) {
	// Setup
	request := graphql.NewRequest(`
		query ($id: ID!) {
			pet (input: { id: $id }) {
				id
				name
				age
				owner
			}
		}
	`)
	request.Var("id", id)
	request.Header.Set("Authorization", UserToken.AccessTokenString)

	// Execute
	var response map[string]interface{}
	err := GraphQlClient.Run(context.Background(), request, &response)
	var pet model.Pet
	mapstructure.Decode(response["pet"], &pet)

	// Verify
	stepName := "getPet"
	if expectedPet != nil {
		assert.Nil(t, err, stepName+": should not error")
		assert.Equal(t, *expectedPet, pet, stepName+": should find the correct pet")
	} else {
		assert.NotNil(t, err, stepName+": should not find pet with id "+id)
		assert.Equal(t, "", pet.Id, stepName+": should not find pet with id "+id)
	}
}

func listPets(t *testing.T) {
	// Setup
	request := graphql.NewRequest(`
		query {
			pets (input: { first: 1, after: "" }) {
				totalCount
				edges {
					node {
						id
					}
					cursor
				}
				pageInfo {
					endCursor
					hasNextPage
				}
			}
		}
	`)
	request.Header.Set("Authorization", UserToken.AccessTokenString)

	// Execute
	var response map[string]interface{}
	err := GraphQlClient.Run(context.Background(), request, &response)
	var connection model.PetConnection
	mapstructure.Decode(response["pets"], &connection)

	// Verify
	stepName := "listPets"
	assert.Nil(t, err, stepName+": should not error")
	if err != nil {
		return
	}
	assert.Equal(t, 1, len(connection.Edges), stepName+": should return 1 pet")
	assert.GreaterOrEqual(t, connection.TotalCount, 2, stepName+": should have total count of at least 2")
	lastEdge := connection.Edges[len(connection.Edges)-1]
	assert.Equal(t, lastEdge.Cursor, connection.PageInfo.EndCursor, stepName+": should have correct end cursor")
	assert.Equal(t, true, connection.PageInfo.HasNextPage, stepName+": should have next page")
}

func updatePetOwner(t *testing.T, pet model.Pet, newOwner string) {
	// Setup
	request := graphql.NewRequest(`
		mutation ($id: ID!, $owner: String!) {
			updatePetOwner (input: { id: $id, owner: $owner }) {
				pet {
					id
					owner
				}
			}
		}
	`)
	request.Var("id", pet.Id)
	request.Var("owner", newOwner)
	request.Header.Set("Authorization", UserToken.AccessTokenString)

	// Execute
	var response map[string]interface{}
	err := GraphQlClient.Run(context.Background(), request, &response)
	var updatedPet model.UpdatePetOwnerPayload
	mapstructure.Decode(response["updatePetOwner"], &updatedPet)

	// Verify
	stepName := "updatePetOwner"
	assert.Nil(t, err, stepName+": should not error")
	assert.NotEqual(t, pet.Owner, newOwner, stepName+": new owner should not match current owner")
	assert.Equal(t, newOwner, updatedPet.Pet.Owner, stepName+": should should update the owner")
}
