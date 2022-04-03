import "@aws-amplify/ui-react/styles.css";
import { Auth } from "../service";

function Logout() {
  const { signOut } = Auth.useAuth();
  signOut();

  return (
    <div className="Page">
      <h2>Logged out!</h2>
    </div>
  );
}

export default Logout;
