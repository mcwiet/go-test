import logo from "./logo.svg";
import "./App.css";
import { Amplify, API, graphqlOperation } from "aws-amplify";
import { pets as petsQuery } from "./graphql/queries";
import { useEffect, useState } from "react";
import { PetEdge } from "./api";
import { Collapse, List, ListItemButton, ListItemText } from "@mui/material";
import { ExpandLess, ExpandMore } from "@mui/icons-material";

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
  const [pets, setPets] = useState<PetEdge[]>([]);
  const [showPetDetails, setShowPetDetails] = useState<Map<string, boolean>>(new Map<string, boolean>());

  useEffect(() => {
    getData();

    async function getData() {
      const response = await API.graphql(graphqlOperation(petsQuery, { input: { first: "100" } }));
      const data = (response as any).data.pets.edges as PetEdge[];
      var details = new Map<string, boolean>();
      data.forEach((edge) => {
        details.set(edge.node.id, false);
      });

      setPets(data);
      setShowPetDetails(details);
    }
  }, []);

  const handleClick = (id: string) => {
    const updated = showPetDetails.set(id, !showPetDetails.get(id));
    setShowPetDetails(new Map<string, boolean>(updated));
  };

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        {pets && (
          <List sx={{ width: "100%", maxWidth: 720, bgcolor: "transparent" }}>
            {pets.map((pet) => (
              <div key={pet.node.id}>
                <ListItemButton onClick={() => handleClick(pet.node.id)}>
                  {pet.node.name}
                  {showPetDetails.get(pet.node.id) ? <ExpandLess /> : <ExpandMore />}
                </ListItemButton>
                <Collapse in={showPetDetails.get(pet.node.id)} timeout="auto" unmountOnExit>
                  <List sx={{ textAlign: "left", paddingTop: 0, paddingLeft: 5 }}>
                    <ListItemText
                      primary={
                        <div>
                          <b>ID:</b> <i>{pet.node.id}</i> <br />
                          <b>Age:</b> <i>{pet.node.age}</i> <br />
                        </div>
                      }
                    />
                  </List>
                </Collapse>
              </div>
            ))}
          </List>
        )}
      </header>
    </div>
  );
}

export default App;
