import { CATEGORIES,
         SEGMENTS,
         CHECK_ITEMS } from '../actions/types';

const START_STATE = { error: '', message: '', categories:'', content: '', authenticated: false}

export default function (state = START_STATE, action) {
  console.log(action.type);
  switch(action.type) {
    case CATEGORIES:
      return { ...state, categories: action.payload };
    case SEGMENTS:
      return { ...state, error: action.payload };
    case CHECK_ITEMS:
      return { ...state, error: action.payload };
  }

  return state;
}