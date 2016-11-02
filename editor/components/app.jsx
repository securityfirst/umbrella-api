import React, { Component } from 'react';
import BurgerMenu from 'react-burger-menu';
import classNames from 'classnames';
import { connect } from 'react-redux';
import MenuWrap from './menuWrap';
import AppContent from './appContent';
import axios from 'axios';
import { Link } from 'react-router';
import { getRepos } from '../actions';
import '../public/stylesheets/menu.scss';

class App extends Component {

  constructor(props) {
    super(props);

    this.props.getRepos();
  }

  getCategories() {
    var categories = [];
    if(this.props.categories) {
      for (var i = this.props.categories.length - 1; i >= 0; i--) {
        var r = this.props.categories[i];
        if (r.parent == 0) {
          categories.push(<span><Link key="{r.id}" to={`/category/${r.id}`}><i className="fa fa-fw fa-bell-o"></i>{r.category}</Link></span>);
        };
      }
    }
    return categories;
  }

  getMenu() {
    const Menu = BurgerMenu['slide'];
    return (
      <MenuWrap wait={20}>
        <Menu id={ "slide" } pageWrapId={"page-wrap"} outerContainerId={"outer-container"}>
          {this.getCategories()}
        </Menu>
      </MenuWrap>
    );
  }

  render() {
    return (
      <div id="outer-container" style={{height: '100%'}}>
        {this.getMenu()}
        <main id="page-wrap">
          <AppContent children={this.props.children}/>
        </main>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return { categories: state.categoryReducer.categories };
}

function mapDispatchToProps(dispatch) {
  return {
    getRepos: function () {
      return dispatch(getRepos());
    }
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(App);
