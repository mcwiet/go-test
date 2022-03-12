/* tslint:disable */
/* eslint-disable */
// this is an auto generated file. This will be overwritten

export const pet = /* GraphQL */ `
  query Pet($input: PetInput!) {
    pet(input: $input) {
      id
      name
      age
      owner
    }
  }
`;
export const pets = /* GraphQL */ `
  query Pets($input: PetsInput!) {
    pets(input: $input) {
      totalCount
      edges {
        node {
          id
          name
          age
          owner
        }
        cursor
      }
      pageInfo {
        endCursor
        hasNextPage
      }
    }
  }
`;
export const user = /* GraphQL */ `
  query User($input: UserInput!) {
    user(input: $input) {
      username
      email
      name
    }
  }
`;
export const users = /* GraphQL */ `
  query Users($input: UsersInput!) {
    users(input: $input) {
      totalCount
      edges {
        node {
          username
          email
          name
        }
      }
      pageInfo {
        endCursor
        hasNextPage
      }
    }
  }
`;
