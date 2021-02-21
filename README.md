<h1 align="center">Welcome to Lending Swap 👋</h1>

<p>
  <img src="https://img.shields.io/badge/version-1.0-blue.svg?cacheSeconds=2592000" />
  
</p>

## Introduction

Lending Swap is an application used to swap atoms between tokens, native tokens from Ethereum to Binance Smart Chain and vice versa.

Between two bridges of Ethereum - BSC, there are 2 relayers to be able to mint corresponding tokens when swapping between the two platform.

There will be validators at each bridge to verify the atomic swaps, when **75%** consensus is reached, the transaction will be executed by the relayer.

In order to be able to transfer the equivalent value of tokens, for example **ETH -> ONE**, **KNC (ERC20) -> WETH (HRC20)**, we use price data taken from **Band Oracle**.

<div float: left;
  width: 33.33%;
  padding: 5px;>
<image src='./readme-images/Ethereum-icon.png'  width='10%'/>
<image src='./readme-images/aave.png'  width='10%'/>
<image src='./readme-images/band-logo.png'  width='10%'/>
<image src='./readme-images/binance.png'  width='10%'/>
</div>

## Usage

### **Ethereum -> BSC**

- ETH -> BNB
- ETH -> Wrapped ETH
- ERC20 -> BNB
- ERC20 -> BEP-20
- Wrapped BNB (ERC20) -> BNB

### **BSC -> Ethereum**

- BNB -> ETH
- BNB -> Wrapped BNB
- BEP-20 -> ERC20
- Wrapped BNB (BEP-20) -> ETH
- BEP-20 -> ETH

###

## Architecture

### Sequence Diagram

- Ethereum -> Binance smart chain

<image src='./readme-images/ethtobsc.png' />

- Binance smart chain -> Ethereum

<image src='./readme-images/bsctoeth.png' />

### Lending

On the ethereum side, when a user performs transactions to swap to Binance Smart Chain, the token amount will be locked and then staked in **AAVE** (has been upgraded to version 2). The yield obtained from the lending will be returned to the validators involved in verifying the swap transactions.

<image src='./readme-images/lending.jpg'>

### Swap

There are 2 BridgeBank contracts on 2 sides that are responsible for receiving transactions to swap tokens and emit out events for the other contract. Contract on that side will be based on the price data taken from the **Band Protocol** to be able to unlock the corresponding amount of tokens. Before being unlocked, the transaction will have to get **75% validator** approval

## Technical

### Frontend:

<p align="center">
<image src='./readme-images/react.png' width='40%'/>
<p>
### Oracle Protocol:
<p align="center">
<image src='./readme-images/band.png' width='30%' padding='20%'/>
</p>
### Lengding platform:
<p align="center">
<image src='./readme-images/aave_logo.jpg' width='30%'>
</p>
### Smart contract:

- Main contract in Binance Smart Chain:

  - BridgeBank.sol
  - BridgeRegistry.sol
  - EthereumBridge.sol
  - Valset.sol

- Main contract in Ethereum:

  - BridgeBank.sol
  - BridgeRegistry.sol
  - BscBridge.sol
  - Valset.sol

## In the next version :herb:

- Set up a consensus mechanism for validators to verify transactions when unlocking tokens. Validators will receive a reward from the interest received from lending pools when validating transactions

- Optimize to choose the most profitable lending platform to optimize rewards for validators
