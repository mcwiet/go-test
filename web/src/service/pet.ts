import { Graphql, SchemaMutation, SchemaQuery, SchemaType } from "../graphql";
import { Response, User } from "../model";

export async function listPets(
  user: User | null | undefined,
  input: SchemaType.PetsInput
): Promise<Response<SchemaType.PetEdge[]>> {
  const ret: Response<SchemaType.PetEdge[]> = {
    data: Array<SchemaType.PetEdge>(),
    error: null,
  };
  try {
    const response = await Graphql.Query<SchemaType.PetsQuery>(user, SchemaQuery.pets, { input: input });
    ret.data = response.data?.pets.edges as SchemaType.PetEdge[];
  } catch (e) {
    console.error(e);
    ret.error = "Error listing pets";
  }
  return ret;
}

export async function createPet(
  user: User | null | undefined,
  input: SchemaType.CreatePetInput
): Promise<Response<SchemaType.Pet>> {
  const ret: Response<SchemaType.Pet> = {
    data: null as unknown as SchemaType.Pet,
    error: null,
  };
  try {
    const response = await Graphql.Mutate<SchemaType.CreatePetMutation>(user, SchemaMutation.createPet, {
      input: input,
    });
    ret.data = response.data?.createPet.pet as SchemaType.Pet;
  } catch (e) {
    console.error(e);
    ret.error = "Error creating pet";
  }
  return ret;
}

export async function updatePetOwner(
  user: User | null | undefined,
  input: SchemaType.UpdatePetOwnerInput
): Promise<Response<SchemaType.Pet>> {
  const ret: Response<SchemaType.Pet> = {
    data: null as unknown as SchemaType.Pet,
    error: null,
  };
  try {
    const response = await Graphql.Mutate<SchemaType.UpdatePetOwnerMutation>(user, SchemaMutation.updatePetOwner, {
      input: input,
    });
    ret.data = response.data?.updatePetOwner.pet as SchemaType.Pet;
  } catch (e) {
    console.error(e);
    ret.error = "Error updating pet owner";
  }
  return ret;
}
