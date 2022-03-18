import "../App.css";
import { API, graphqlOperation } from "aws-amplify";
import { pets as petsQuery } from "../graphql/queries";
import { useEffect, useState } from "react";
import { PetEdge } from "../api";
import { Collapse, List, ListItemButton, ListItemText } from "@mui/material";
import { ExpandLess, ExpandMore } from "@mui/icons-material";

function ViewPets() {
  const [pets, setPets] = useState<PetEdge[]>([]);
  const [showPetDetails, setShowPetDetails] = useState<Map<string, boolean>>(new Map<string, boolean>());

  useEffect(() => {
    getData();

    async function getData() {
      try {
        const response = await API.graphql(graphqlOperation(petsQuery, { input: { first: "100" } }));
        const data = (response as any).data.pets.edges as PetEdge[];
        var details = new Map<string, boolean>();
        data.forEach((edge) => {
          details.set(edge.node.id, false);
        });

        setPets(data);
        setShowPetDetails(details);
      } catch (e) {
        console.warn(e);
      }
    }
  }, []);

  const handleClick = (id: string) => {
    const updated = showPetDetails.set(id, !showPetDetails.get(id));
    setShowPetDetails(new Map<string, boolean>(updated));
  };

  return (
    <div className="Page">
      <h2>Pets</h2>
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
    </div>
  );
}

export default ViewPets;
