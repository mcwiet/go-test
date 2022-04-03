import { ExpandLess, ExpandMore } from "@mui/icons-material";
import { Collapse, List, ListItemButton, ListItemText } from "@mui/material";
import { useEffect, useState } from "react";
import "../App.css";
import { SchemaType } from "../graphql";
import { PageProps } from "../model";
import { Pet } from "../service";

function ViewPets(props: PageProps) {
  const [pets, setPets] = useState<SchemaType.PetEdge[]>([]);
  const [showPetDetails, setShowPetDetails] = useState<Map<string, boolean>>(new Map<string, boolean>());

  useEffect(() => {
    getData();

    async function getData() {
      try {
        const response = await Pet.listPets(props.user, { first: 100 });

        var details = new Map<string, boolean>();
        response.data.forEach((edge) => {
          details.set(edge.node.id, false);
        });

        setPets(response.data);
        setShowPetDetails(details);
      } catch (e) {
        console.warn(e);
      }
    }
  }, [props.user]);

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
