import React, { Component } from 'react';

class MenuWrap extends Component {
  constructor() {
    super();
    this.hidden = false;
  }

  render() {
    let style;

    if (this.hidden) {
      style = {display: 'none'};
    }

    return (
      <div style={style}>
        {this.props.children}
      </div>
    );
  }
}

export default MenuWrap;