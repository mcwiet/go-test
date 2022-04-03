import { GraphQLAPI, GraphQLResult } from "@aws-amplify/api-graphql";
import { GRAPHQL_AUTH_MODE } from "@aws-amplify/auth";
import { User } from "../model";

export async function Query<Output>(
  user: User | null | undefined,
  queryString: string,
  variables: object | undefined
): Promise<GraphQLResult<Output>> {
  let response = (await GraphQLAPI.graphql({
    query: queryString,
    variables: variables,
    authMode: GetAuthMode(user),
  })) as GraphQLResult<Output>;

  if (response.errors) {
    throw new Error(response.errors.toString());
  }

  return response;
}

export async function Mutate<Output>(
  user: User | null | undefined,
  mutationString: string,
  variables: object | undefined
): Promise<GraphQLResult<Output>> {
  let response = (await GraphQLAPI.graphql({
    query: mutationString,
    variables: variables,
    authMode: GetAuthMode(user),
  })) as GraphQLResult<Output>;

  if (response.errors) {
    throw new Error(response.errors.toString());
  }

  return response;
}

function GetAuthMode(user: User | null | undefined): GRAPHQL_AUTH_MODE {
  if (user === undefined) {
    throw new Error("User is not defined; should be 'null' if not authenticated or a valid user if authenticated");
  }
  return user ? GRAPHQL_AUTH_MODE.AMAZON_COGNITO_USER_POOLS : GRAPHQL_AUTH_MODE.AWS_IAM;
}
