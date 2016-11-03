import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware } from 'redux';
import { Router, browserHistory } from 'react-router';
import reduxThunk from 'redux-thunk';
import reducers from './reducers/index';
import routes from './routes';
import './public/stylesheets/base.css';

const composer = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose
const createStoreWithMiddleware = composer(applyMiddleware(reduxThunk))(createStore);
const store = createStoreWithMiddleware(reducers);

ReactDOM.render(
  <Provider store={store}>
    <Router history={browserHistory} routes={routes} />
  </Provider>,
  document.querySelector('#wrapper'));