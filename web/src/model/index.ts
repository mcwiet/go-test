export { Observable } from "../../node_modules/zen-observable-ts";

export interface User {
  username: string;
  email: string;
  name: string;
}

export interface PageProps {
  user: User | null | undefined;
}

export interface Response<T> {
  data: T;
  error: string | null;
}
