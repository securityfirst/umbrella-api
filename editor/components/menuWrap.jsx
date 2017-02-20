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
        c.slug = c.name.replace(/\s/g, "-").toLowerCase();
        content.push({"id":c.slug, "icon": "icon-class-name", "label":c.name, "to": "/category/"+c.slug});
        if (c.subcategories.length>1) { // 1 means there is only the basic category info
          for (var k = c.subcategories.length - 1; k >= 0; k--) {
            var sc = c.subcategories[k];
            sc.slug = sc.name.replace(/\s/g, "-").toLowerCase();
            content.push({"id":sc.slug, "parentId": c.slug, "icon": "icon-class-name", "label":sc.name, "to": "/category/"+sc.slug});
          }
        }
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