import axios from 'axios';
import { browserHistory } from 'react-router';
import cookie from 'react-cookie';
import { AUTH_USER,
         ACTIVE_CAT,
         TREE,
         FETCH_USER,
         SEGMENTS,
         CHECK_ITEMS } from './types';

const API_URL = 'http://127.0.0.1:8080/v2';

axios.interceptors.request.use(function (config) {
  config.headers['Authorization'] = 'Bearer ' + cookie.load('github-email');
  return config;
});

export function errorHandler(dispatch, error, type) {
  let errorMessage = '';
  if(error.data.error) {
    errorMessage = error.data.error;
  } else if(error.data) {
    errorMessage = error.data;
  } else {
    errorMessage = error;
  }

  if(error.status === 401 || error.status === 403) {
    dispatch({
      type: type,
      payload: 'You are not authorized to do this. Please login and try again.'
    });
    browserHistory.push('/login');
  } else {
    dispatch({
      type: type,
      payload: errorMessage
    });
  }
}

export function getRepos() {
  return function(dispatch) {
    axios.get(`${API_URL}/api/tree`)
    .then(response => {
      dispatch({
        type: TREE,
        payload: response.data
      });
    })
  }
}

export function getInfo() {
  return function(dispatch) {
    axios.get(`${API_URL}/`)
    .then(response => {
      dispatch({
        type: FETCH_USER,
        payload: response.data
      });
      dispatch({
        type: AUTH_USER
      });
    })
    .catch(function (error) {
      return errorHandler(dispatch, error.response, UNAUTH_USER);
    });
  }
}

export function getCat(categoryName) {
  return function(dispatch) {
    if (typeof myVar == 'undefined') return
    axios.get(`${API_URL}/repo`+categoryName)
    .then(response => {
      dispatch({
        type: ACTIVE_CAT,
        payload: response.data
      });
      dispatch({
        type: AUTH_USER
      });
    })
    .catch(function (error) {
      return errorHandler(dispatch, error.response, UNAUTH_USER);
    });
  }
}