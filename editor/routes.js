import React from 'react';
import { Route, IndexRoute } from 'react-router';

import App from './components/app';
import NotFoundPage from './components/pages/not-found';
import HomePage from './components/pages/home';
import LoginPage from './components/auth/login';
import LogoutPage from './components/auth/logout';
import UserProfile from './components/pages/profile';
import Category from './components/pages/category';

export default (
  <Route path="/" component={App}>
    <IndexRoute component={HomePage} />
    <Route name="categoryDetails" path="category/:categoryName/:subCategoryName" component={Category} />
    <Route name="categoryDetails" path="category/:categoryName" component={Category} />
    <Route name="loginPage" path="login" component={LoginPage} />
    <Route name="logout" path="logout" component={LogoutPage} />
    <Route name="profilePage" path="profile" component={UserProfile} />
    <Route path="*" component={NotFoundPage} />
  </Route>
);