import * as actions from './action';
const initState = {
  address: null,
  hmy: null,
  isAuthorized: false,
  bridgeBankInstance: null,
  chainId: 0
};

const reducer = (state = initState, action) => {
  switch (action.type) {
    case actions.SET_ADDRESS:
      return {
        ...state,
        address: action.address
      };
    case actions.SET_HMY:
      return {
        ...state,
        hmy: action.hmy
      };
    case actions.SET_IS_AUTHORIZED:
      return {
        ...state,
        isAuthorized: action.isAuthorized
      };
    case actions.SET_CHAINID:
      return {
        ...state,
        chainId: action.chainId
      };
    case actions.SET_BRIDGE_BANK:
      return {
        ...state,
        bridgeBankInstance: action.bridgeBankInstance
      };
    default:
      return state;
  }
};

export { reducer as hmyReducer };
