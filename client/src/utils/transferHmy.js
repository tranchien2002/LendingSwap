import { Harmony } from '@harmony-js/core';
import { ChainID, ChainType, Unit } from '@harmony-js/utils';
import BridgeBankHmy from 'contracts/BridgeBankHmy.json';
import { getHmyContractAddress } from 'utils/getHmyContractAddress';

const ERC20 = require('contracts/IERC20.json');

const options = {
  gasLimit: 6721900,
  gasPrice: 1000000000
};

const hmy = new Harmony('https://api.s0.b.hmny.io', {
  chainType: ChainType.Harmony,
  chainId: ChainID.HmyTestnet
});

export const approveHmy = async (walletAddress, tokenAddress, chainId) => {
  const contractEthAddress = getHmyContractAddress(chainId);
  const erc20 = hmy.contracts.createContract(ERC20.abi, tokenAddress);
  try {
    erc20.wallet.defaultSigner = hmy.crypto.getAddress(walletAddress).checksum;
    erc20.wallet.signTransaction = async tx => {
      try {
        tx.from = hmy.crypto.getAddress(walletAddress).checksum;
        const signTx = await window.harmony.signTransaction(tx);
        return signTx;
      } catch (e) {
        console.log(e);
      }
    };
    let weiValue = new Unit(1000000).asOne().toWei();
    await erc20.methods.approve(contractEthAddress.bridgeBank, weiValue).send({ ...options });
  } catch (e) {
    console.log(e);
  }
};

export const swapToken_1_1_hmy = async (walletAddress, receiver, tokenAddress, amount, chainId) => {
  const contractAddress = getHmyContractAddress(chainId);
  const bridgeBank = hmy.contracts.createContract(BridgeBankHmy.abi, contractAddress.bridgeBank);
  try {
    bridgeBank.wallet.defaultSigner = hmy.crypto.getAddress(walletAddress).checksum;
    bridgeBank.wallet.signTransaction = async tx => {
      try {
        tx.from = hmy.crypto.getAddress(walletAddress).checksum;
        const signTx = await window.harmony.signTransaction(tx);
        return signTx;
      } catch (e) {
        console.log(e);
      }
    };
    let weiValue = new Unit(amount).asOne().toWei();
    await bridgeBank.methods
      .swapToken_1_1(receiver, tokenAddress, weiValue)
      .send({ ...options })
      .then(e => console.log(e));
  } catch (e) {
    console.log(e);
  }
};

export const swapONEForETH = async (walletAddress, receiver, amount, chainId) => {
  const contractAddress = getHmyContractAddress(chainId);
  const bridgeBank = hmy.contracts.createContract(BridgeBankHmy.abi, contractAddress.bridgeBank);
  try {
    bridgeBank.wallet.defaultSigner = hmy.crypto.getAddress(walletAddress).checksum;
    bridgeBank.wallet.signTransaction = async tx => {
      try {
        tx.from = hmy.crypto.getAddress(walletAddress).checksum;
        const signTx = await window.harmony.signTransaction(tx);
        return signTx;
      } catch (e) {
        console.log(e);
      }
    };
    let weiValue = new Unit(amount).asOne().toWei();
    await bridgeBank.methods
      .swapONEForETH(receiver, weiValue)
      .send({ ...options, value: weiValue });
  } catch (e) {
    console.log(e);
  }
};

export const swapONEForWONE = async (walletAddress, receiver, amount, chainId) => {
  const contractAddress = getHmyContractAddress(chainId);
  const bridgeBank = hmy.contracts.createContract(BridgeBankHmy.abi, contractAddress.bridgeBank);
  try {
    bridgeBank.wallet.defaultSigner = hmy.crypto.getAddress(walletAddress).checksum;
    bridgeBank.wallet.signTransaction = async tx => {
      try {
        tx.from = hmy.crypto.getAddress(walletAddress).checksum;
        const signTx = await window.harmony.signTransaction(tx);
        return signTx;
      } catch (e) {
        console.log(e);
      }
    };
    let weiValue = new Unit(amount).asOne().toWei();
    await bridgeBank.methods
      .swapONEForWrappedONE(receiver, weiValue)
      .send({ ...options, value: weiValue });
  } catch (e) {
    console.log(e);
  }
};

export const swapONEForToken = async (walletAddress, receiver, amount, destToken, chainId) => {
  const contractAddress = getHmyContractAddress(chainId);
  const bridgeBank = hmy.contracts.createContract(BridgeBankHmy.abi, contractAddress.bridgeBank);
  try {
    bridgeBank.wallet.defaultSigner = hmy.crypto.getAddress(walletAddress).checksum;
    bridgeBank.wallet.signTransaction = async tx => {
      try {
        tx.from = hmy.crypto.getAddress(walletAddress).checksum;
        const signTx = await window.harmony.signTransaction(tx);
        return signTx;
      } catch (e) {
        console.log(e);
      }
    };
    let weiValue = new Unit(amount).asOne().toWei();
    await bridgeBank.methods
      .swapONEForToken(receiver, weiValue, destToken)
      .send({ ...options, value: weiValue });
  } catch (e) {
    console.log(e);
  }
};

export const swapTokenForToken_hmy = async (
  walletAddress,
  receiver,
  hmyToken,
  amount,
  destToken,
  chainId
) => {
  const contractAddress = getHmyContractAddress(chainId);
  const bridgeBank = hmy.contracts.createContract(BridgeBankHmy.abi, contractAddress.bridgeBank);
  try {
    bridgeBank.wallet.defaultSigner = hmy.crypto.getAddress(walletAddress).checksum;
    bridgeBank.wallet.signTransaction = async tx => {
      try {
        tx.from = hmy.crypto.getAddress(walletAddress).checksum;
        const signTx = await window.harmony.signTransaction(tx);
        return signTx;
      } catch (e) {
        console.log(e);
      }
    };
    let weiValue = new Unit(amount).asOne().toWei();
    await bridgeBank.methods
      .swapTokenForToken(receiver, hmyToken, weiValue, destToken)
      .send({ ...options });
  } catch (e) {
    console.log(e);
  }
};

export const swapTokenForWONE = async (walletAddress, receiver, harmonyToken, amount, chainId) => {
  const contractAddress = getHmyContractAddress(chainId);
  const bridgeBank = hmy.contracts.createContract(BridgeBankHmy.abi, contractAddress.bridgeBank);

  try {
    bridgeBank.wallet.defaultSigner = hmy.crypto.getAddress(walletAddress).checksum;
    bridgeBank.wallet.signTransaction = async tx => {
      try {
        tx.from = hmy.crypto.getAddress(walletAddress).checksum;
        const signTx = await window.harmony.signTransaction(tx);
        return signTx;
      } catch (e) {
        console.log(e);
      }
    };
    let weiValue = new Unit(amount).asOne().toWei();
    await bridgeBank.methods
      .swapTokenForWONE(receiver, harmonyToken, weiValue)
      .send({ ...options });
  } catch (e) {
    console.log(e);
  }
};

export const swapTokenForETH = async (walletAddress, receiver, hmyToken, amount, chainId) => {
  const contractAddress = getHmyContractAddress(chainId);
  const bridgeBank = hmy.contracts.createContract(BridgeBankHmy.abi, contractAddress.bridgeBank);

  try {
    bridgeBank.wallet.defaultSigner = hmy.crypto.getAddress(walletAddress).checksum;
    bridgeBank.wallet.signTransaction = async tx => {
      try {
        tx.from = hmy.crypto.getAddress(walletAddress).checksum;
        const signTx = await window.harmony.signTransaction(tx);
        return signTx;
      } catch (e) {
        console.log(e);
      }
    };
    let weiValue = new Unit(amount).asOne().toWei();
    await bridgeBank.methods.swapTokenForETH(receiver, hmyToken, weiValue).send({ ...options });
  } catch (e) {
    console.log(e);
  }
};
