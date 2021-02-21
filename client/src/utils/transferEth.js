import Web3 from 'web3';
import BridgeBank from 'contracts/BridgeBank.json';
import { getETHContractAddress } from 'utils/getETHContractAddress';
import { ChainID, ChainType } from '@harmony-js/utils';
import { Harmony } from '@harmony-js/core';
import { message } from 'antd';
const ERC20 = require('contracts/IERC20.json');
const web3 = new Web3(window.ethereum);

const hmy = new Harmony('https://api.s0.b.hmny.io', {
  chainType: ChainType.Harmony,
  chainId: ChainID.HmyTestnet
});

export const approve_Eth = async (walletAddress, tokenAddress, chainId) => {
  const contractEthAddress = getETHContractAddress(chainId);
  const erc20 = new web3.eth.Contract(ERC20.abi, tokenAddress);
  const weiValue = web3.utils.toWei('1000000000', 'ether');
  await erc20.methods
    .approve(contractEthAddress.bridgeBank, weiValue)
    .send({ from: walletAddress });
};

export const swapToken_1_1 = async (walletAddress, receiver, tokenAddress, amount, chainId) => {
  const contractAddress = getETHContractAddress(chainId);
  const bridgeBank = new web3.eth.Contract(BridgeBank.abi, contractAddress.bridgeBank);
  let receiver_hmy = hmy.crypto.getAddress(receiver).checksum;
  const weiValue = web3.utils.toWei(amount, 'ether');
  await bridgeBank.methods
    .swapToken_1_1(receiver_hmy, tokenAddress, weiValue)
    .send({ from: walletAddress })
    .then(e => message.success(e.blockHash));
};

export const swapETHForONE = async (walletAddress, receiver, amount, chainId) => {
  const contractAddress = getETHContractAddress(chainId);
  const bridgeBank = new web3.eth.Contract(BridgeBank.abi, contractAddress.bridgeBank);
  let receiver_hmy = hmy.crypto.getAddress(receiver).checksum;
  const weiValue = web3.utils.toWei(amount, 'ether');
  await bridgeBank.methods
    .swapETHForONE(receiver_hmy, weiValue)
    .send({ value: weiValue, from: walletAddress })
    .then(e => message.success(e.blockHash));
};

export const swapETHForWETH = async (walletAddress, receiver, amount, chainId) => {
  const contractAddress = getETHContractAddress(chainId);
  const bridgeBank = new web3.eth.Contract(BridgeBank.abi, contractAddress.bridgeBank);
  let receiver_hmy = hmy.crypto.getAddress(receiver).checksum;
  const weiValue = web3.utils.toWei(amount, 'ether');
  await bridgeBank.methods
    .swapETHForWrappedETH(receiver_hmy, weiValue)
    .send({ value: weiValue, from: walletAddress })
    .then(e => message.success(e.blockHash));
};

export const swapETHForToken = async (walletAddress, receiver, amount, destToken, chainId) => {
  const contractAddress = getETHContractAddress(chainId);
  const bridgeBank = new web3.eth.Contract(BridgeBank.abi, contractAddress.bridgeBank);
  let receiver_hmy = hmy.crypto.getAddress(receiver).checksum;
  const weiValue = web3.utils.toWei(amount, 'ether');
  await bridgeBank.methods
    .swapETHForToken(receiver_hmy, weiValue, destToken)
    .send({ value: weiValue, from: walletAddress })
    .then(e => message.success(e.blockHash));
};
export const swapTokenForToken = async (
  walletAddress,
  receiver,
  ethToken,
  amount,
  destToken,
  chainId
) => {
  const contractAddress = getETHContractAddress(chainId);
  const bridgeBank = new web3.eth.Contract(BridgeBank.abi, contractAddress.bridgeBank);
  let receiver_hmy = hmy.crypto.getAddress(receiver).checksum;
  const weiValue = web3.utils.toWei(amount, 'ether');
  await bridgeBank.methods
    .swapTokenForToken(receiver_hmy, ethToken, weiValue, destToken)
    .send({ from: walletAddress })
    .then(e => message.success(e.blockHash));
};
export const swapTokenForWETH = async (walletAddress, receiver, ethToken, amount, chainId) => {
  const contractAddress = getETHContractAddress(chainId);
  const bridgeBank = new web3.eth.Contract(BridgeBank.abi, contractAddress.bridgeBank);
  let receiver_hmy = hmy.crypto.getAddress(receiver).checksum;
  const weiValue = web3.utils.toWei(amount, 'ether');
  console.log('wei', weiValue);
  await bridgeBank.methods
    .swapTokenForWrappedETH(receiver_hmy, ethToken, weiValue)
    .send({ from: walletAddress })
    .then(e => message.success(e.blockHash));
};
export const swapTokenForONE = async (walletAddress, receiver, ethToken, amount, chainId) => {
  const contractAddress = getETHContractAddress(chainId);
  const bridgeBank = new web3.eth.Contract(BridgeBank.abi, contractAddress.bridgeBank);
  let receiver_hmy = hmy.crypto.getAddress(receiver).checksum;
  const weiValue = web3.utils.toWei(amount, 'ether');
  await bridgeBank.methods
    .swapTokenForONE(receiver_hmy, ethToken, weiValue)
    .send({ from: walletAddress })
    .then(e => message.success(e.blockHash));
};
