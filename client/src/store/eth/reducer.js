import * as actions from './action';
const initState = {
  address: null,
  web3: null,
  chainId: 42,
  bridgeBankInstance: null,
  isAuthorized: false
};

const reducer = (state = initState, action) => {
  switch (action.type) {
    case actions.SET_ADDRESS:
      return {
        ...state,
        address: action.address
      };
    case actions.SET_BRIDGE_BANK:
      return {
        ...state,
        bridgeBankInstance: action.bridgeBankInstance
      };
    case actions.SET_CHAINID:
      return {
        ...state,
        chainId: action.chainId
      };
    case actions.SET_WEB3:
      return {
        ...state,
        web3: action.web3
      };
    case actions.SET_IS_AUTHORIZED:
      return {
        ...state,
        isAuthorized: action.isAuthorized
      };
    default:
      return state;
  }
};

export { reducer as ethReducer };
