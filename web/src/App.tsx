import "./App.css";
import logo from "./logo.svg";
import { Amplify } from "aws-amplify";

const myAppConfig = {
  aws_appsync_graphqlEndpoint: process.env["REACT_APP_API_URL"],
  aws_appsync_region: process.env["REACT_APP_AWS_REGION"],
  aws_appsync_authenticationType: "AWS_IAM",
  Auth: {
    identityPoolId: process.env["REACT_APP_IDENTITY_POOL_ID"],
    region: process.env["REACT_APP_AWS_REGION"],
  },
};

Amplify.configure(myAppConfig);

function App() {
  return (
    <div className="Page">
      <h2>Home</h2>
      <img src={logo} className="App-logo" alt="logo" />
    </div>
  );
}

export default App;
