import React, { Component } from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router'

class LoginPage extends Component {

  constructor(props) {
    super(props);
  }

  render() {
    if (this.props.loggedIn) { browserHistory.push('/'); }
    return (
      <div>This is the login page route.</div>
    );
  }
}

function mapStateToProps(state) {
  return {
    loggedIn: state.auth.authenticated,
  };
}

export default connect(mapStateToProps)(LoginPage)