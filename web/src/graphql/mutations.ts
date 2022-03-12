/* tslint:disable */
/* eslint-disable */
// this is an auto generated file. This will be overwritten

export const createPet = /* GraphQL */ `
  mutation CreatePet($input: CreatePetInput!) {
    createPet(input: $input) {
      pet {
        id
        name
        age
        owner
      }
    }
  }
`;
export const deletePet = /* GraphQL */ `
  mutation DeletePet($input: DeletePetInput!) {
    deletePet(input: $input) {
      id
    }
  }
`;
export const updatePetOwner = /* GraphQL */ `
  mutation UpdatePetOwner($input: UpdatePetOwnerInput!) {
    updatePetOwner(input: $input) {
      pet {
        id
        name
        age
        owner
      }
    }
  }
`;
