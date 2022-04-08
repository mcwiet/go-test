import { ExpandLess, ExpandMore } from "@mui/icons-material";
import { Button, Collapse, Grid, List, ListItemButton, ListItemText } from "@mui/material";
import { useEffect, useState } from "react";
import "../App.css";
import { SchemaType } from "../graphql";
import { PageProps } from "../model";
import { Pet } from "../service";

function ViewPets(props: PageProps) {
  const [allPets, setAllPets] = useState<SchemaType.PetEdge[]>([]);
  const [userPets, setUserPets] = useState<SchemaType.PetEdge[]>([]);
  const [unownedPets, setUnownedPets] = useState<SchemaType.PetEdge[]>([]);
  const [otherUserPets, setOtherUserPets] = useState<SchemaType.PetEdge[]>([]);
  const [showPetDetails, setShowPetDetails] = useState<Map<string, boolean>>(new Map<string, boolean>());

  useEffect(() => {
    const getPets = async () => {
      const response = await Pet.listPets(props.user, { first: 100 });
      if (!response.data) return;

      const details = new Map<string, boolean>();
      response.data.forEach((edge) => {
        details.set(edge.node.id, false);
      });

      setAllPets(response.data);
      setUnownedPets(response.data.filter((edge) => !edge.node.owner));
      setUserPets(response.data.filter((edge) => edge.node.owner === props.user?.username));
      setOtherUserPets(response.data.filter((edge) => edge.node.owner && edge.node.owner !== props.user?.username));

      setShowPetDetails(details);
    };

    getPets();
  }, [props.user]);

  const handleClick = (petId: string) => {
    const updated = showPetDetails.set(petId, !showPetDetails.get(petId));
    setShowPetDetails(new Map<string, boolean>(updated));
  };

  const claimPet = async (petId: string) => {
    await Pet.updatePetOwner(props.user, { id: petId, owner: props.user!.username });
    const response = await Pet.listPets(props.user, { first: 100 });
    if (!response.data) return;

    const details = new Map<string, boolean>();
    response.data.forEach((edge) => {
      details.set(edge.node.id, false);
    });

    setAllPets(response.data);
    setUnownedPets(response.data.filter((edge) => !edge.node.owner));
    setUserPets(response.data.filter((edge) => edge.node.owner === props.user?.username));
    setOtherUserPets(response.data.filter((edge) => edge.node.owner && edge.node.owner !== props.user?.username));

    setShowPetDetails(details);
  };

  const releasePet = async (petId: string) => {
    await Pet.updatePetOwner(props.user, { id: petId, owner: "" });
    const response = await Pet.listPets(props.user, { first: 100 });
    if (!response.data) return;

    const details = new Map<string, boolean>();
    response.data.forEach((edge) => {
      details.set(edge.node.id, false);
    });

    setAllPets(response.data);
    setUnownedPets(response.data.filter((edge) => !edge.node.owner));
    setUserPets(response.data.filter((edge) => edge.node.owner === props.user?.username));
    setOtherUserPets(response.data.filter((edge) => edge.node.owner && edge.node.owner !== props.user?.username));

    setShowPetDetails(details);
  };

  return props.user ? (
    <div className="Page">
      <h2>Pets</h2>
      <h4>Your Pets</h4>
      {userPets && (
        <List sx={{ width: "100%", maxWidth: 720, bgcolor: "transparent" }}>
          {userPets.map((pet) => (
            <div key={pet.node.id}>
              <Grid container>
                <Grid item xs={9}>
                  <ListItemButton onClick={() => handleClick(pet.node.id)}>
                    {pet.node.name}
                    {showPetDetails.get(pet.node.id) ? <ExpandLess /> : <ExpandMore />}
                  </ListItemButton>
                </Grid>
                <Grid item xs={3}>
                  <Button onClick={() => releasePet(pet.node.id)} color="error">
                    Release
                  </Button>
                </Grid>
              </Grid>
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
      <h4>Unowned Pets</h4>
      {unownedPets && (
        <List sx={{ width: "100%", maxWidth: 720, bgcolor: "transparent" }}>
          {unownedPets.map((pet) => (
            <div key={pet.node.id}>
              <Grid container>
                <Grid item xs={9}>
                  <ListItemButton onClick={() => handleClick(pet.node.id)}>
                    {pet.node.name}
                    {showPetDetails.get(pet.node.id) ? <ExpandLess /> : <ExpandMore />}
                  </ListItemButton>
                </Grid>
                <Grid item xs={3}>
                  <Button onClick={() => claimPet(pet.node.id)}>Claim</Button>
                </Grid>
              </Grid>
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
      <h4>Other Users' Pets</h4>
      {otherUserPets && (
        <List sx={{ width: "100%", maxWidth: 720, bgcolor: "transparent" }}>
          {otherUserPets.map((pet) => (
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
  ) : (
    <div className="Page">
      <h2>Pets</h2>
      {allPets && (
        <List sx={{ width: "100%", maxWidth: 720, bgcolor: "transparent" }}>
          {allPets.map((pet) => (
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
