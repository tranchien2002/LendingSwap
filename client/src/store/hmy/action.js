import { Harmony } from '@harmony-js/core';
import { ChainID, ChainType } from '@harmony-js/utils';
import WalletConnectProvider from '@walletconnect/web3-provider';

const provider = new WalletConnectProvider({
  rpc: {
    0: 'https://api.s0.b.hmny.io',
    1: 'https://api.s1.b.hmny.io'
    // ...
  },
  qrcodeModalOptions: {
    mobileLinks: ['trust']
  }
});

const hmy = new Harmony('https://api.s0.b.hmny.io', {
  chainType: ChainType.Harmony,
  chainId: ChainID.HmyTestnet
});

export const SET_HMY = '@@hmy/SET_HMY';
export const setHmy = hmy => async dispatch => {
  dispatch({ type: SET_HMY, hmy });
};

export const SET_ADDRESS = '@@hmy/SET_ADDDRESS';
export const setAddress = address => async dispatch => {
  dispatch({ type: SET_ADDRESS, address });
};

export const SET_IS_AUTHORIZED = '@@hmy/SET_IS_AUTHORIZED';
export const setIsAuthorized = isAuthorized => async dispatch => {
  dispatch({ type: SET_IS_AUTHORIZED, isAuthorized });
};

export const SET_BRIDGE_BANK = '@@hmy/SET_BRIDGE_BANK';
export const setBridgeBank = bridgeBankInstance => dispatch => {
  dispatch({ type: SET_BRIDGE_BANK, bridgeBankInstance });
};

export const SET_CHAINID = '@@hmy/SET_CHAINID';
export const setChainId = chainId => dispatch => {
  dispatch({ type: SET_CHAINID, chainId });
};
