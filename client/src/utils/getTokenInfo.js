import dai from 'assets/icons/dai.png';
import knc from 'assets/icons/knc.png';
import link from 'assets/icons/link.png';
import eth from 'assets/icons/eth.png';
import one from 'assets/icons/one.png';
import Web3 from 'web3';
import { getETHContractAddress } from 'utils/getETHContractAddress';
import { getHmyContractAddress } from 'utils/getHmyContractAddress';
import { Harmony } from '@harmony-js/core';
import { ChainID, ChainType, fromWei, hexToNumber, Units } from '@harmony-js/utils';
const ERC20 = require('contracts/IERC20.json');
const { Client } = require('@bandprotocol/bandchain.js');
const options = {
  gasLimit: 6721900,
  gasPrice: 1000000000
};
const hmy = new Harmony('https://api.s0.b.hmny.io', {
  chainType: ChainType.Harmony,
  chainId: ChainID.HmyTestnet
});

const tokenInfo = {
  42: [
    {
      symbol: 'DAI',
      ethAddress: '0xff795577d9ac8bd7d90ee22b6c1703490b6512fd',
      hmyAddress: '0xff795577d9ac8bd7d90ee22b6c1703490b6512fd',
      icon: dai
    },
    {
      symbol: 'KNC',
      ethAddress: '0x3F80c39c0b96A0945f9F0E9f55d8A8891c5671A8',
      hmyAddress: 'one1uz230xxr88yrhs709fc865yunlf72prrefcm29',
      icon: knc
    },
    {
      symbol: 'LINK',
      ethAddress: '0xAD5ce863aE3E4E9394Ab43d4ba0D80f419F61789',
      hmyAddress: '0xff795577d9ac8bd7d90ee22b6c1703490b6512fd',
      icon: link
    },
    {
      symbol: 'ETH',
      ethAddress: '0x0000000000000000000000000000000000000001',
      hmyAddress: 'one16vxu2p4v7qf65nkeclll974ckvv0rjcyv3lc8a',
      icon: eth
    },
    {
      symbol: 'ONE',
      ethAddress: '0x503bE5F89B1d0342880D983724bD2e3E9a827904',
      hmyAddress: '0x0000000000000000000000000000000000000001',
      icon: one
    }
  ]
};

export const getSymbolETH = (_chainId, _address) => {
  let symbol = '';
  let listToken = tokenInfo[_chainId] ? tokenInfo[_chainId] : [];
  listToken.forEach(token => {
    if (token.ethAddress === _address) {
      symbol = token.symbol;
    }
  });
  return symbol;
};

export const getAddressToken = (_chainId, _symbol) => {
  let address = '';
  let listToken = tokenInfo[_chainId] ? tokenInfo[_chainId] : [];
  listToken.forEach(token => {
    if (token.symbol === _symbol) {
      address = token.ethAddress;
    }
  });
  return address;
};

export const getHmyAddressToken = (_chainId, _symbol) => {
  let address = '';
  let listToken = tokenInfo[_chainId] ? tokenInfo[_chainId] : [];
  listToken.forEach(token => {
    if (token.symbol === _symbol) {
      address = token.hmyAddress;
    }
  });
  return address;
};

export const getIconETH = (_chainId, _address) => {
  let icon;
  let listToken = tokenInfo[_chainId] ? tokenInfo[_chainId] : [];
  listToken.forEach(token => {
    if (token.ethAddress === _address) {
      icon = token.icon;
    }
  });
  return icon;
};

export const convertToken = async (src, target, amount) => {
  const endpoint = 'https://api-gm-lb.bandchain.org';
  const bandchain = new Client(endpoint);
  const res = await bandchain.getReferenceData([src + '/' + target]);
  let rate;
  if (res.length > 0) {
    rate = res[0].rate;
  }
  return rate * amount;
};

// chain: 0: ETH 1: Harmony
export const balanceOf = async (tokenAddress, walletAddress, roadSwap) => {
  let web3 = new Web3(window.ethereum);
  let balance;
  if (!walletAddress) {
    return 0;
  }
  if (roadSwap) {
    if (tokenAddress === '0x0000000000000000000000000000000000000001') {
      balance = await hmy.blockchain.getBalance({ address: walletAddress });
      balance = parseInt(balance.result);
    } else {
      const erc20 = hmy.contracts.createContract(ERC20.abi, tokenAddress);
      balance = await erc20.methods.balanceOf(walletAddress).call(options);
      balance = parseInt(balance);
    }
  } else {
    if (tokenAddress === '0x0000000000000000000000000000000000000001') {
      balance = web3.eth.getBalance(walletAddress);
    } else {
      const erc20 = new web3.eth.Contract(ERC20.abi, tokenAddress);
      balance = await erc20.methods.balanceOf(walletAddress).call();
    }
  }
  return balance;
};

// chain: 0: ETH 1: Harmony
export const allowance = async (tokenAddress, walletAddress, chainId, roadSwap) => {
  let web3 = new Web3(window.ethereum);
  let erc20;
  let tokenAllowance = 0;
  if (!walletAddress) {
    return 0;
  }
  if (roadSwap) {
    if (tokenAddress !== '0x0000000000000000000000000000000000000001') {
      erc20 = hmy.contracts.createContract(ERC20.abi, tokenAddress);
      let contractAddress = getHmyContractAddress(chainId);
      tokenAllowance = await erc20.methods
        .allowance(walletAddress, contractAddress.bridgeBank)
        .call(options);
      tokenAllowance = parseInt(tokenAllowance);
    }
  } else {
    if (tokenAddress !== '0x0000000000000000000000000000000000000001') {
      erc20 = new web3.eth.Contract(ERC20.abi, tokenAddress);
      let contractAddress = getETHContractAddress(chainId);
      tokenAllowance = await erc20.methods
        .allowance(walletAddress, contractAddress.bridgeBank)
        .call();
    }
  }
  return tokenAllowance;
};
