import React, { Component } from 'react';
import { LinkContainer, IndexLinkContainer } from 'react-router-bootstrap';
import { connect } from 'react-redux';
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
          <IndexLinkContainer to="/">
            <a href="#">
              <Navbar.Brand>
                Umbrella Admin
              </Navbar.Brand>
            </a>
          </IndexLinkContainer>
          <Navbar.Toggle />
        </Navbar.Header>
        { this.props.loggedIn &&
          <Navbar.Collapse>
            <Nav pullRight>
              <NavDropdown eventKey={1} title="User" id="basic-nav-dropdown">
                <LinkContainer to="/profile">
                  <MenuItem eventKey={1.1}>User Profile</MenuItem>
                </LinkContainer>
                <MenuItem divider />
                <LinkContainer to="/logout">
                <MenuItem eventKey={1.3}>Log out</MenuItem>
                </LinkContainer>
              </NavDropdown>
            </Nav>
          </Navbar.Collapse>
        }
      </Navbar>
    );
  }

}

function mapStateToProps(state) {
  return {
    loggedIn: state.auth.authenticated,
  };
}

export default connect(mapStateToProps)(MainNavBar);