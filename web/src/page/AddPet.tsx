import "@aws-amplify/ui-react/styles.css";
import { Box, Button, TextField } from "@mui/material";
import { useForm } from "react-hook-form";
import "../App.css";
import { CreatePetInput } from "../graphql/api";
import { PageProps } from "../model";
import { Pet } from "../service";

function AddPet(props: PageProps) {
  const { register, handleSubmit } = useForm();
  const onSubmit = (data: any) => {
    create(data as CreatePetInput);

    async function create(pet: CreatePetInput) {
      await Pet.createPet(props.user, pet);
      console.log("created!");
    }
  };

  return props.user ? (
    <div className="Page">
      <h2>Add a Pet</h2>
      <Box
        textAlign={"center"}
        component="form"
        sx={{
          "& .MuiTextField-root": { m: 1, width: "25ch" },
          backgroundColor: "primary.light",
          borderRadius: "16px",
        }}
        autoComplete="off"
        onSubmit={handleSubmit(onSubmit)}
      >
        <div>
          <h4>Pet Info</h4>
        </div>
        <div>
          <TextField
            {...register("name")}
            label="Name"
            sx={{ input: { color: "white" } }}
            InputLabelProps={{ style: { color: "white" } }}
            required
            id="outlined-required"
            defaultValue="Levi"
            variant="filled"
          />
        </div>
        <div>
          <TextField
            {...register("age")}
            label="Age"
            sx={{ input: { color: "white" } }}
            InputLabelProps={{ style: { color: "white" } }}
            inputProps={{ inputMode: "numeric", pattern: "[0-9]*" }}
            required
            id="outlined-required"
            defaultValue="2"
            variant="filled"
          />
        </div>
        <div>
          <Button sx={{ marginTop: 3, marginBottom: 1 }} type="submit" variant="contained">
            Submit
          </Button>
        </div>
      </Box>
    </div>
  ) : (
    <div className="Page">
      <h4>You must be signed in to view this page</h4>
    </div>
  );
}

export default AddPet;
