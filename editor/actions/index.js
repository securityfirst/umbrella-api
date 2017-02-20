import axios from 'axios';
import { browserHistory } from 'react-router';
import cookie from 'react-cookie';
import { CATEGORIES,
         SEGMENTS,
         CHECK_ITEMS } from './types';

const API_URL = 'http://127.0.0.1:8080/v2';

export function errorHandler(dispatch, error, type) {
  let errorMessage = '';

  if(error.data.error) {
    errorMessage = error.data.error;
  } else if(error.data) {
    errorMessage = error.data;
  } else {
    errorMessage = error;
  }

  if(error.status === 401) {
    dispatch({
      type: type,
      payload: 'You are not authorized to do this. Please login and try again.'
    });
  } else {
    dispatch({
      type: type,
      payload: errorMessage
    });
  }
}

export function getRepos() {
  console.log("get repos");
  return function(dispatch) {
    axios.get(`${API_URL}/api/tree`)
    .then(response => {
      dispatch({
        type: CATEGORIES,
        payload: response.data
      });
      // console.log(response.data);
    })
  }
}