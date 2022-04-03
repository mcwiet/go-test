import "@aws-amplify/ui-react/styles.css";
import "../App.css";
import logo from "../logo.svg";
import { PageProps } from "../model";

function Home(props: PageProps) {
  return (
    <div className="Page">
      <h2>Home</h2>
      <h4>{props.user ? `Welcome, ${props.user?.email}!` : ""}</h4>
      <img src={logo} className="App-logo" alt="logo" />
    </div>
  );
}

export default Home;
