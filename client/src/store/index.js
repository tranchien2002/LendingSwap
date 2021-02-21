import { hmyReducer } from './hmy/reducer';
import { ethReducer } from './eth/reducer';
import thunk from 'redux-thunk';
import { createStore, applyMiddleware, combineReducers } from 'redux';
import { composeWithDevTools } from 'redux-devtools-extension';

const rootReducer = combineReducers({
  eth: ethReducer,
  hmy: hmyReducer
});
const store = createStore(rootReducer, composeWithDevTools(applyMiddleware(thunk)));
export default store;
