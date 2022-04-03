import { Auth } from "@aws-amplify/auth";
import { Hub } from "@aws-amplify/core";
import { CognitoUser } from "amazon-cognito-identity-js";
import { useEffect, useState } from "react";
import { promisify } from "util";
import { User } from "../model";

export interface UseAuthHookResponse {
  currentUser: User | null | undefined;
  signOut: () => void;
}

const getCurrentUser = async (): Promise<User | null> => {
  try {
    let cognitoUser: CognitoUser = await Auth.currentAuthenticatedUser();
    const getAttributes = promisify(cognitoUser.getUserAttributes).bind(cognitoUser);
    const attributes = await getAttributes();
    return {
      username: cognitoUser.getUsername(),
      email: attributes?.find((attr) => attr.Name === "email")?.Value ?? "",
      name: attributes?.find((attr) => attr.Name === "name")?.Value ?? "",
    };
  } catch {
    return null;
  }
};

export const useAuth = (): UseAuthHookResponse => {
  const [currentUser, setCurrentUser] = useState<User | null | undefined>(undefined);

  useEffect(() => {
    const updateUser = async () => {
      let user = await getCurrentUser();
      setCurrentUser(user);
    };
    Hub.listen("auth", updateUser); // listen for login/signup events
    updateUser(); // check manually the first time because we won't get a Hub event
    return () => Hub.remove("auth", updateUser);
  }, []);

  const signOut = () => Auth.signOut();

  return { currentUser, signOut };
};
