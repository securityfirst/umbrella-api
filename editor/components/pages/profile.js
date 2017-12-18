import React, { Component } from 'react';
import { connect } from 'react-redux';

class UserProfile extends Component {
  render() {
    return (
      <div>
      	<p>This is the user profile page.</p>
      	<div><pre>{JSON.stringify(this.props.user, null, 2) }</pre></div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    user: state.categories.user,
  };
}

export default connect(mapStateToProps)(UserProfile);