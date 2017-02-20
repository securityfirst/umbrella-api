import {
         AUTH_USER,
         UNAUTH_USER,
         ACTIVE_CAT,
         TREE,
         FETCH_USER,
         SEGMENTS,
         CHECK_ITEMS } from '../actions/types';

const START_STATE = { error: '', message: '', activeCategory:'', tree:'', user:'', content: ''}

export default function (state = START_STATE, action) {
  console.log(action.type);
  switch(action.type) {
    case TREE:
      return { ...state, tree:action.payload };
    case ACTIVE_CAT:
      return { ...state, activeCategory: action.payload };
    case FETCH_USER:
      return { ...state, user: action.payload };
    case UNAUTH_USER:
      return { ...state, user: null, categories: null, tree: null };
    case SEGMENTS:
      return { ...state, error: action.payload };
    case CHECK_ITEMS:
      return { ...state, error: action.payload };
  }

  return state;
}