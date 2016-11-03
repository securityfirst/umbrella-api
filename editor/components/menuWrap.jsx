import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import MetisMenu from 'react-metismenu';

class MenuWrap extends Component {

    constructor(props) {
      super(props);
    }

    getCategoriesForMenu() {
      var content = [];
      for (var i = this.props.categories.length - 1; i >= 0; i--) {
        var c = this.props.categories[i];
        var nc = {"id":c.id, "icon": "icon-class-name", "label":c.category, "to": "/category/"+c.id};
        if (c.parent!=0) {
          nc.parentId = c.parent;
        }
        content.push(nc);
      }
      return content.reverse();
    }

    render() {
      var content = this.getCategoriesForMenu();
        return (
            <div>
                <MetisMenu content={content} ref="menu" activeLinkFromLocation />
            </div>
        );
    }
}

export default MenuWrap;