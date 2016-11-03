import React, { Component } from 'react';
import { Navbar, Nav, NavItem, NavDropdown, MenuItem } from 'react-bootstrap';

class MainNavBar extends Component {

  constructor(props) {
    super(props);
  }

  render() {
    const useIcon = "<i class='fa fa-fw fa-question'></i>"
    return (
      <Navbar collapseOnSelect>
        <Navbar.Header>
          <Navbar.Brand>
            <a href="#">Umbrella Admin</a>
          </Navbar.Brand>
          <Navbar.Toggle />
        </Navbar.Header>
        <Navbar.Collapse>
          <Nav pullRight>
            <NavDropdown eventKey={1} title="User" id="basic-nav-dropdown">
              <MenuItem eventKey={1.1}>User Profile</MenuItem>
              <MenuItem divider />
              <MenuItem eventKey={1.3}>Log out</MenuItem>
            </NavDropdown>
          </Nav>
        </Navbar.Collapse>
      </Navbar>
    );
  }
}

export default MainNavBar;