import "@aws-amplify/ui-react/styles.css";
import { Amplify, Auth as AmplifyAuth } from "aws-amplify";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import "./App.css";
import Nav from "./nav/Nav";
import { Home, Login, Logout, ViewPets } from "./page";
import AddPet from "./page/AddPet";
import { Auth } from "./service";

const authConfig = {
  aws_appsync_graphqlEndpoint: process.env["REACT_APP_API_URL"],
  aws_appsync_region: process.env["REACT_APP_AWS_REGION"],
  aws_appsync_authenticationType: "AMAZON_COGNITO_USER_POOL",
  Auth: {
    identityPoolId: process.env["REACT_APP_IDENTITY_POOL_ID"],
    region: process.env["REACT_APP_AWS_REGION"],
    userPoolId: process.env["REACT_APP_USER_POOL_ID"],
    userPoolWebClientId: process.env["REACT_APP_USER_POOL_APP_CLIENT_ID"],
    jwtToken: async () => {
      try {
        return (await AmplifyAuth.currentSession()).getIdToken().getJwtToken();
      } catch {
        return null;
      }
    },
  },
};

Amplify.configure(authConfig);

function App() {
  const { currentUser } = Auth.useAuth();

  return (
    <BrowserRouter>
      <Nav user={currentUser} />
      <Routes>
        <Route path="/" element={<Home user={currentUser} />} />
        <Route path="pets" element={<ViewPets user={currentUser} />} />
        <Route path="pet/add" element={<AddPet user={currentUser} />} />
        <Route path="login" element={<Login />} />
        <Route path="logout" element={<Logout />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
