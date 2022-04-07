import AppBar from "@mui/material/AppBar";
import Avatar from "@mui/material/Avatar";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Container from "@mui/material/Container";
import IconButton from "@mui/material/IconButton";
import Menu from "@mui/material/Menu";
import Toolbar from "@mui/material/Toolbar";
import Tooltip from "@mui/material/Tooltip";
import Typography from "@mui/material/Typography";
import * as React from "react";
import { Link } from "react-router-dom";
import { PageProps } from "../model";

const pages = new Map<string, string>([
  ["Home", "/"],
  ["View", "/pets"],
  ["Add", "/pet/add"],
]);
const unauthSettings = new Map<string, string>([["Login", "/login"]]);
const authSettings = new Map<string, string>([["Logout", "/logout"]]);

const Nav = (props: PageProps) => {
  const [anchorElUser, setAnchorElUser] = React.useState<null | HTMLElement>(null);

  const handleOpenUserMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorElUser(event.currentTarget);
  };

  const handleCloseUserMenu = () => {
    setAnchorElUser(null);
  };

  return (
    <AppBar position="static">
      <Container maxWidth="xl">
        <Toolbar disableGutters>
          <Typography variant="h6" noWrap component="div" sx={{ mr: 2, display: { xs: "none", md: "flex" } }}>
            PETS!
          </Typography>

          <Box sx={{ flexGrow: 1, display: { xs: "none", md: "flex" } }}>
            {Array.from(pages.entries()).map((page) => (
              <Link to={page[1]} key={page[0]} style={{ textDecoration: "none" }}>
                <Button key={page[0]} sx={{ my: 2, color: "white", display: "block" }}>
                  {page[0]}
                </Button>
              </Link>
            ))}
          </Box>

          <Box>
            <Tooltip title="Open settings">
              <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
                <Avatar>{props.user ? props.user.email[0].toUpperCase() : null}</Avatar>
              </IconButton>
            </Tooltip>
            <Menu
              sx={{ mt: "45px", width: "450px" }}
              id="menu-appbar"
              anchorEl={anchorElUser}
              anchorOrigin={{
                vertical: "top",
                horizontal: "right",
              }}
              keepMounted
              transformOrigin={{
                vertical: "top",
                horizontal: "right",
              }}
              open={Boolean(anchorElUser)}
              onClose={handleCloseUserMenu}
            >
              {Array.from((props.user ? authSettings : unauthSettings).entries()).map((setting) => (
                <Link to={setting[1]} key={setting[0]} style={{ textDecoration: "none" }}>
                  <Button key={setting[0]} sx={{ color: "black", display: "block" }}>
                    {setting[0]}
                  </Button>
                </Link>
              ))}
            </Menu>
          </Box>
        </Toolbar>
      </Container>
    </AppBar>
  );
};
export default Nav;
