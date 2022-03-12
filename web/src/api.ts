/* tslint:disable */
/* eslint-disable */
//  This file was automatically generated and should not be edited.

export type CreatePetInput = {
  name: string,
  age: number,
  owner?: string | null,
};

export type CreatePetPayload = {
  __typename: "CreatePetPayload",
  pet: Pet,
};

export type Pet = {
  __typename: "Pet",
  id: string,
  name: string,
  age: number,
  owner?: string | null,
};

export type DeletePetInput = {
  id: string,
};

export type DeletePetPayload = {
  __typename: "DeletePetPayload",
  id: string,
};

export type UpdatePetOwnerInput = {
  id: string,
  owner: string,
};

export type UpdatePetOwnerPayload = {
  __typename: "UpdatePetOwnerPayload",
  pet: Pet,
};

export type PetInput = {
  id: string,
};

export type PetsInput = {
  first: number,
  after?: string | null,
};

export type PetConnection = {
  __typename: "PetConnection",
  totalCount: number,
  edges?:  Array<PetEdge > | null,
  pageInfo: PageInfo,
};

export type PetEdge = {
  __typename: "PetEdge",
  node: Pet,
  cursor: string,
};

export type PageInfo = {
  __typename: "PageInfo",
  endCursor?: string | null,
  hasNextPage: boolean,
};

export type UserInput = {
  username: string,
};

export type User = {
  __typename: "User",
  username: string,
  email?: string | null,
  name?: string | null,
};

export type UsersInput = {
  first: number,
  after?: string | null,
};

export type UserConnection = {
  __typename: "UserConnection",
  totalCount: number,
  edges?:  Array<UserEdge > | null,
  pageInfo: PageInfo,
};

export type UserEdge = {
  __typename: "UserEdge",
  node: User,
};

export type CreatePetMutationVariables = {
  input: CreatePetInput,
};

export type CreatePetMutation = {
  createPet:  {
    __typename: "CreatePetPayload",
    pet:  {
      __typename: "Pet",
      id: string,
      name: string,
      age: number,
      owner?: string | null,
    },
  },
};

export type DeletePetMutationVariables = {
  input: DeletePetInput,
};

export type DeletePetMutation = {
  deletePet:  {
    __typename: "DeletePetPayload",
    id: string,
  },
};

export type UpdatePetOwnerMutationVariables = {
  input: UpdatePetOwnerInput,
};

export type UpdatePetOwnerMutation = {
  updatePetOwner:  {
    __typename: "UpdatePetOwnerPayload",
    pet:  {
      __typename: "Pet",
      id: string,
      name: string,
      age: number,
      owner?: string | null,
    },
  },
};

export type PetQueryVariables = {
  input: PetInput,
};

export type PetQuery = {
  pet:  {
    __typename: "Pet",
    id: string,
    name: string,
    age: number,
    owner?: string | null,
  },
};

export type PetsQueryVariables = {
  input: PetsInput,
};

export type PetsQuery = {
  pets:  {
    __typename: "PetConnection",
    totalCount: number,
    edges?:  Array< {
      __typename: "PetEdge",
      node:  {
        __typename: "Pet",
        id: string,
        name: string,
        age: number,
        owner?: string | null,
      },
      cursor: string,
    } > | null,
    pageInfo:  {
      __typename: "PageInfo",
      endCursor?: string | null,
      hasNextPage: boolean,
    },
  },
};

export type UserQueryVariables = {
  input: UserInput,
};

export type UserQuery = {
  user:  {
    __typename: "User",
    username: string,
    email?: string | null,
    name?: string | null,
  },
};

export type UsersQueryVariables = {
  input: UsersInput,
};

export type UsersQuery = {
  users:  {
    __typename: "UserConnection",
    totalCount: number,
    edges?:  Array< {
      __typename: "UserEdge",
      node:  {
        __typename: "User",
        username: string,
        email?: string | null,
        name?: string | null,
      },
    } > | null,
    pageInfo:  {
      __typename: "PageInfo",
      endCursor?: string | null,
      hasNextPage: boolean,
    },
  },
};
