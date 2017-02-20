import axios from 'axios';
import { browserHistory } from 'react-router';
import cookie from 'react-cookie';
import { API_URL, errorHandler } from './index';
import { AUTH_USER, AUTH_ERROR, UNAUTH_USER } from './types';

// TO-DO: Add expiration to cookie
export function loginUser({ email, password }) {
  return function (dispatch) {
    axios.post(`${API_URL}/auth/login`, { email, password })
    .then((response) => {
      cookie.save('github-email', response.data.token, { path: '/' });
      dispatch({ type: AUTH_USER });
      browserHistory.push('/');
    })
    .catch((error) => {
      errorHandler(dispatch, error.response, AUTH_ERROR);
    });
  };
}

export function logoutUser() {
  return function (dispatch) {
    dispatch({ type: UNAUTH_USER });
    cookie.remove('github-email');
  };
}