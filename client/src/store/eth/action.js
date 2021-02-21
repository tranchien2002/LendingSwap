import Web3 from 'web3';
import WalletConnectProvider from '@walletconnect/web3-provider';
import BridgeBank from 'contracts/BridgeBank.json';
import { getETHContractAddress } from 'utils/getETHContractAddress';

var contractAddress;

//  Create WalletConnect Provider
const provider = new WalletConnectProvider({
  infuraId: '27e484dcd9e3efcfd25a83a78777cdf1',
  qrcodeModalOptions: {
    mobileLinks: ['metamask']
  }
});

export const connectWalletConnect = () => async dispatch => {
  await provider.enable();
  const web3 = new Web3(provider);
  dispatch({ type: SET_WEB3, web3 });
};

export const SET_ADDRESS = '@@eth/SET_ADDDRESS';
export const setAddress = address => async dispatch => {
  dispatch({ type: SET_ADDRESS, address });
};

export const SET_IS_AUTHORIZED = '@@eth/SET_IS_AUTHORIZED';
export const setIsAuthorized = isAuthorized => async dispatch => {
  dispatch({ type: SET_IS_AUTHORIZED, isAuthorized });
};

export const SET_CHAINID = '@@eth/SET_CHAINID';
export const setChainId = chainId => dispatch => {
  dispatch({ type: SET_CHAINID, chainId });
};

export const SET_WEB3 = '@@eth/SET_WEB3';
export const setWeb3 = web3 => async (dispatch, getState) => {
  dispatch({ type: SET_WEB3, web3 });
  if (web3) {
    let state = getState();
    let { chainId } = state.eth;
    contractAddress = getETHContractAddress(chainId);

    const bridgeBankInstance = new web3.eth.Contract(BridgeBank.abi, contractAddress.bridgeBank);
    dispatch(setBridgeBank(bridgeBankInstance));
  }
};

export const SET_BRIDGE_BANK = '@@eth/SET_BRIDGE_BANK';
export const setBridgeBank = bridgeBankInstance => dispatch => {
  dispatch({ type: SET_BRIDGE_BANK, bridgeBankInstance });
};
