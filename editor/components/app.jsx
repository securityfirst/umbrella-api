import React, { Component } from 'react';
import classNames from 'classnames';
import { connect } from 'react-redux';
import AppContent from './appContent';
import MenuWrap from './menuWrap';
import MainNavBar from './mainNavBar';
import axios from 'axios';
import { Link } from 'react-router';
import { getRepos, getInfo } from '../actions';
import { Button, Grid, Row, Col } from 'react-bootstrap';
import '../public/stylesheets/metisMenu.css';
require('font-awesome/css/font-awesome.css');

class App extends Component {

  constructor(props) {
    super(props);

    this.props.getRepos();
    this.props.getInfo();
  }

  render() {
    return (
      <Grid>
        <Row className="show-grid">
          <MainNavBar loggedIn1={this.props.loggedIn} />
        </Row>
        <Row className="show-grid" id="menu">
          <Col xs={3} md={3} lg={3}>
            { this.props.loggedIn && <MenuWrap/> }
          </Col>
          <Col xs={9} md={9} lg={9}>
            <AppContent children={this.props.children}/>
          </Col>
        </Row>
      </Grid>
    );
  }
}

function mapStateToProps(state) {
  return {
    loggedIn: state.auth.authenticated,
  };
}

function mapDispatchToProps(dispatch) {
  return {
    getInfo: function () {
      return dispatch(getInfo());
    },
    getRepos: function () {
      return dispatch(getRepos());
    }
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(App);
