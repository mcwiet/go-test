import { Authenticator, components, Flex } from "@aws-amplify/ui-react";
import "@aws-amplify/ui-react/styles.css";
import "../App.css";

function Login() {
  return (
    <Flex className="Page" justifyContent="center">
      <Authenticator components={components}>
        {({ user }) => (
          <div className="Page">
            <h4>{user.attributes?.email} is now logged in!</h4>
          </div>
        )}
      </Authenticator>
    </Flex>
  );
}

export default Login;
