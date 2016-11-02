import React from 'react';
import { Route, IndexRoute } from 'react-router';

import App from './components/app';
import NotFoundPage from './components/pages/not-found';
import HomePage from './components/pages/home';
import Category from './components/pages/category';

export default (
  <Route path="/" component={App}>
    <IndexRoute component={HomePage} />
    <Route name="categoryDetails" path="category/:categoryName" component={Category} />
    <Route path="*" component={NotFoundPage} />
  </Route>
);