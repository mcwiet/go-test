import { render, screen } from "@testing-library/react";
import { graphqlOperation } from "aws-amplify";
import { pets as petsQuery } from "./graphql/queries";
import App from "./App";

jest.mock("aws-amplify");

test("renders pets header", () => {
  render(<App />);
  const header = screen.getByText("Pets");
  expect(header).toBeInTheDocument();
});

test("queries for pets", () => {
  render(<App />);
  expect(graphqlOperation).toHaveBeenLastCalledWith(petsQuery, expect.anything());
});
