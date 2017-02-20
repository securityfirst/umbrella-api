import { combineReducers } from 'redux';
import categoryReducer from './categoryReducer'
import authReducer from './authReducer'
import { routerReducer } from 'react-router-redux'

export default combineReducers({
  routing: routerReducer,
  categories: categoryReducer,
  auth: authReducer,
});