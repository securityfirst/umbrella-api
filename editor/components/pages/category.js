import React, { Component } from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router'
import { getCat } from '../../actions';

class Item extends Component {

  constructor(props) {
    super(props);

    this.props.getCat();
  }

  render() {
    return (
		<div>
			<p>This is the category {this.props.category}.</p>
      		{this.props.category !== '' && <div><pre>{JSON.stringify(this.props.category, null, 2) }</pre></div>}
		</div>
    );
  }
}

function mapStateToProps(state, ownProps) {
  return {
    category: state.routing.locationBeforeTransitions.pathname,
  };
}

function mapDispatchToProps(dispatch) {
  return {
    getCat: function () {
      return dispatch(getCat());
    }
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Item);