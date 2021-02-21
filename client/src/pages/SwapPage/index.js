import React, { useEffect, useMemo, useState } from 'react';
import { useSelector } from 'react-redux';
import { Row, Col, Input, Select } from 'antd';
import { tokens } from 'constant/support';
import WalletsETH from 'components/WalletsETH';
import WalletsHmy from 'components/WalletsHmy';
import {
  allowance,
  balanceOf,
  getAddressToken,
  getHmyAddressToken,
  convertToken
} from 'utils/getTokenInfo';
import {
  swapToken_1_1,
  swapETHForONE,
  swapETHForWETH,
  swapETHForToken,
  swapTokenForToken,
  swapTokenForWETH,
  swapTokenForONE,
  approve_Eth
} from 'utils/transferEth';

import {
  swapToken_1_1_hmy,
  swapONEForETH,
  swapONEForWONE,
  swapONEForToken,
  swapTokenForToken_hmy,
  swapTokenForWONE,
  swapTokenForETH,
  approveHmy
} from 'utils/transferHmy';
import {
  RightOutlined,
  ArrowRightOutlined,
  ArrowDownOutlined,
  SwapOutlined
} from '@ant-design/icons';
import useInterval from 'utils/useInterval';

import './style.scss';

const { Option } = Select;

const listWallets = {
  ETH: <WalletsETH />,
  ONE: <WalletsHmy />
};

function SwapPage() {
  const addressETH = useSelector(state => state.eth.address);
  const addressHmy = useSelector(state => state.hmy.address);
  const ethChainId = useSelector(state => state.eth.chainId);
  const hmyChainId = useSelector(state => state.hmy.chainId);
  const listAddressDest = useMemo(() => {
    return { 1: addressETH, 0: addressHmy };
  }, [addressETH, addressHmy]);

  const [indexRoadSwap, setIndexRoadSwap] = useState(0);
  const [disableBtnSwap, setDisableBtnSwap] = useState(false);
  const [amountSource, setAmountSource] = useState(0);
  const [amountDest, setAmoutDest] = useState(0);
  const [tokenSource, setTokenSource] = useState('ONE');
  const [tokenDest, setTokenDest] = useState('ETH');
  const [walletSource, setWalletSource] = useState('ETH');
  const [walletDest, setWalletDest] = useState('ONE');
  const [addressDest, setAddressDest] = useState();
  const [toAddress, setToAddress] = useState();
  const [approvedToken, setApprovedToken] = useState(true);
  const [balanceSource, setBalanceSource] = useState(0);
  useInterval(async () => {
    if (tokenSource !== tokenDest) {
      let amountDest = await convertToken(tokenSource, tokenDest, amountSource);
      setAmoutDest(amountDest);
    } else {
      setAmoutDest(amountSource);
    }
    // checkBeforeSwap();
  }, 3000);

  useEffect(() => {
    setAddressDest(listAddressDest[indexRoadSwap]);
    setToAddress(listAddressDest[indexRoadSwap]);
    calcBalance(tokenSource);
    checkBeforeSwap();
  }, [listAddressDest, indexRoadSwap, addressDest, tokenDest, tokenSource, amountSource]);

  function reverseDirectionToken() {
    setTokenSource(tokenDest);
    setTokenDest(tokenSource);
  }

  async function calcBalance(tokenSymbol) {
    let tokenAddress = getAddressToken(ethChainId, tokenSymbol);
    let balanceToken;
    if (tokenSource !== tokenDest) {
      let amountDest = await convertToken(tokenSource, tokenDest, amountSource);
      setAmoutDest(amountDest);
    } else {
      setAmoutDest(amountSource);
    }
    if (indexRoadSwap) {
      // Harmony -> ETH
      tokenAddress = getHmyAddressToken(ethChainId, tokenSymbol);
      balanceToken = await balanceOf(tokenAddress, addressHmy, indexRoadSwap);
      setBalanceSource(balanceToken);
    } else {
      // ETH -> Harmony
      tokenAddress = getAddressToken(ethChainId, tokenSymbol);
      balanceToken = await balanceOf(tokenAddress, addressETH, indexRoadSwap);
      setBalanceSource(balanceToken);
    }
  }

  async function onChangeTokenSource(value) {
    setTokenSource(value);
  }

  function onChangeTokenDest(value) {
    setTokenDest(value);
  }

  const onChangeAddressTo = e => {
    const { value } = e.target;
    setToAddress(value);
  };

  const onBlurAddressTo = e => {
    const { value } = e.target;
    setToAddress(value);
  };

  const onChangeFormatNumber = async e => {
    const { value } = e.target;
    const reg = /^-?\d*(\.\d*)?$/;
    if ((!isNaN(value) && reg.test(value)) || value === '') {
      setAmountSource(value);
    }
  };

  const chooseRoad = (from, to, index) => {
    setIndexRoadSwap(index);
    setWalletSource(from);
    setWalletDest(to);
    setAddressDest(listAddressDest[index]);
    setToAddress(listAddressDest[index]);
  };

  const onBlurFormatNumber = e => {
    const { value } = e.target;
    let valueTemp = value;
    if (value.charAt(value.length - 1) === '.') {
      valueTemp = value.slice(0, -1);
    }
    setAmountSource(valueTemp.replace(/0*(\d+)/, '$1'));
  };

  function setMyAdress() {
    setToAddress(addressDest);
  }

  async function approveToken() {
    let tokenAddress;
    if (indexRoadSwap) {
      tokenAddress = getHmyAddressToken(ethChainId, tokenSource);
      await approveHmy(addressHmy, tokenAddress, hmyChainId);
    } else {
      tokenAddress = getAddressToken(ethChainId, tokenSource);
      await approve_Eth(addressETH, tokenAddress, ethChainId);
    }
    await checkBeforeSwap();
  }

  async function checkBeforeSwap() {
    let tokenAddress;
    let balanceToken;
    let allowanceToken;
    console.log('checbefore', indexRoadSwap);
    if (indexRoadSwap) {
      // Harmony -> ETH
      tokenAddress = getHmyAddressToken(ethChainId, tokenSource);
      balanceToken = await balanceOf(tokenAddress, addressHmy, indexRoadSwap);
      allowanceToken = await allowance(tokenAddress, addressHmy, hmyChainId, indexRoadSwap);
      if (
        parseFloat(allowanceToken) < parseFloat(amountSource) * 10 ** 18 &&
        tokenSource !== 'ONE'
      ) {
        setApprovedToken(false);
      } else if (balanceToken < amountSource * 10 ** 18) {
        setDisableBtnSwap(true);
      } else {
        setApprovedToken(true);
        setDisableBtnSwap(false);
      }
    } else {
      // ETH -> Harmony
      tokenAddress = getAddressToken(ethChainId, tokenSource);
      balanceToken = await balanceOf(tokenAddress, addressETH, indexRoadSwap);
      allowanceToken = await allowance(tokenAddress, addressETH, ethChainId, indexRoadSwap);
      if (
        parseFloat(allowanceToken) < parseFloat(amountSource) * 10 ** 18 &&
        tokenSource !== 'ETH'
      ) {
        setApprovedToken(false);
      } else if (balanceToken < amountSource * 10 ** 18) {
        setDisableBtnSwap(true);
      } else {
        setApprovedToken(true);
        setDisableBtnSwap(false);
      }
    }
  }

  const swap = async () => {
    console.log('token', tokenSource, tokenDest);
    if (indexRoadSwap) {
      // Harmony -> ETH
      if (tokenSource === 'ONE') {
        if (tokenDest === 'ONE') {
          await swapONEForWONE(addressHmy, addressETH, amountSource, hmyChainId);
        } else if (tokenDest === 'ETH') {
          await swapONEForETH(addressHmy, addressETH, amountSource, hmyChainId);
        } else {
          let addressDest = getHmyAddressToken(ethChainId, tokenDest);
          await swapONEForToken(addressHmy, addressETH, amountSource, addressDest, hmyChainId);
        }
      } else if (tokenSource === tokenDest) {
        let addressSource = getHmyAddressToken(ethChainId, tokenSource);
        await swapToken_1_1_hmy(addressHmy, addressETH, addressSource, amountSource, hmyChainId);
      } else {
        if (tokenDest === 'ETH') {
          let addressSource = getHmyAddressToken(ethChainId, tokenSource);
          await swapTokenForETH(addressHmy, addressETH, addressSource, amountSource, hmyChainId);
        } else if (tokenDest === 'ONE') {
          let addressSource = getHmyAddressToken(ethChainId, tokenSource);
          await swapTokenForWONE(addressHmy, addressETH, addressSource, amountSource, hmyChainId);
        } else {
          let addressSource = getHmyAddressToken(ethChainId, tokenSource);
          let addressDest = getHmyAddressToken(ethChainId, tokenDest);
          await swapTokenForToken_hmy(
            addressHmy,
            addressETH,
            addressSource,
            amountSource,
            addressDest,
            hmyChainId
          );
        }
      }
    } else {
      // ETH -> Harmony
      if (tokenSource === 'ETH') {
        if (tokenDest === 'ETH') {
          await swapETHForWETH(addressETH, addressHmy, amountSource, ethChainId);
        } else if (tokenDest === 'ONE') {
          await swapETHForONE(addressETH, addressHmy, amountSource, ethChainId);
        } else {
          let addressDest = getAddressToken(ethChainId, tokenDest);
          await swapETHForToken(addressETH, addressHmy, amountSource, addressDest, ethChainId);
        }
      } else if (tokenSource === tokenDest) {
        let addressSource = getAddressToken(ethChainId, tokenSource);
        await swapToken_1_1(addressETH, addressHmy, addressSource, amountSource, ethChainId);
      } else {
        if (tokenDest === 'ONE') {
          let addressSource = getAddressToken(ethChainId, tokenSource);
          await swapTokenForONE(addressETH, addressHmy, addressSource, amountSource, ethChainId);
        } else if (tokenDest === 'ETH') {
          let addressSource = getAddressToken(ethChainId, tokenSource);
          await swapTokenForWETH(addressETH, addressHmy, addressSource, amountSource, ethChainId);
        } else {
          let addressSource = getAddressToken(ethChainId, tokenSource);
          let addressDest = getAddressToken(ethChainId, tokenDest);
          await swapTokenForToken(
            addressETH,
            addressHmy,
            addressSource,
            amountSource,
            addressDest,
            ethChainId
          );
        }
      }
    }
  };

  return (
    <div className='swap-page'>
      <Row className='container' justify='space-between'>
        <Col xs={{ order: 2, span: 24 }} md={{ order: 1, span: 10 }}>
          <div className='content-swap'>
            <div className='input-and-select-token'>
              <div className='token-source-dest'>
                <div className='label-input-token'>
                  <div className='sc-hSdWYo euiRCS css-1rhdhic'>From</div>
                </div>
                <div className='box-input-token'>
                  <Input
                    size='large'
                    placeholder='0.0'
                    className='input-token'
                    onChange={e => onChangeFormatNumber(e)}
                    onBlur={e => onBlurFormatNumber(e)}
                    value={amountSource}
                  />
                  <Select
                    value={tokenSource}
                    style={{ width: 150 }}
                    showSearch
                    placeholder='Select a token'
                    optionFilterProp='children'
                    onChange={onChangeTokenSource}
                    filterOption={(input, option) =>
                      option.children[1].toLowerCase().indexOf(input.toLowerCase()) >= 0
                    }
                    className='button-select-token'
                  >
                    {tokens.map((token, i) => {
                      return (
                        <Option value={token.symbol} key={i}>
                          <img
                            alt='icon-token'
                            src={token.icon}
                            className='img-icon-token-select-option'
                          />
                          {token.symbol}
                        </Option>
                      );
                    })}
                  </Select>
                </div>
              </div>
              <p className='balance'>Balance: {parseFloat(balanceSource / 1e18).toFixed(2)}</p>
              <div className='symbol-arrow-down'>
                <ArrowDownOutlined
                  className='icon'
                  onClick={() => {
                    reverseDirectionToken();
                  }}
                />
              </div>
              <div className='token-source-dest'>
                <div className='label-input-token'>
                  <div className='sc-hSdWYo euiRCS css-1rhdhic'>To</div>
                </div>
                <div className='box-input-token'>
                  <Input
                    size='large'
                    placeholder='0.0'
                    className='input-token'
                    value={amountDest}
                    disabled
                  />
                  <Select
                    value={tokenDest}
                    style={{ width: 150 }}
                    showSearch
                    placeholder='Select a token'
                    optionFilterProp='children'
                    onChange={onChangeTokenDest}
                    filterOption={(input, option) =>
                      option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
                    }
                    className='button-select-token'
                  >
                    {tokens.map((token, i) => {
                      return (
                        <Option value={token.symbol} key={i}>
                          <img
                            alt='icon-token'
                            src={token.icon}
                            className='img-icon-token-select-option'
                          />
                          {token.symbol}
                        </Option>
                      );
                    })}
                  </Select>
                </div>
              </div>
              <div>
                <div className='label-input-token'>
                  <div className='sc-hSdWYo euiRCS css-1rhdhic'>To Address</div>
                </div>
                <div className='box-input-token'>
                  <Input
                    size='large'
                    placeholder='address...'
                    className='input-token'
                    value={toAddress}
                    onChange={e => onChangeAddressTo(e)}
                    onBlur={e => onBlurAddressTo(e)}
                  />
                </div>
                {addressDest ? (
                  <div className='button-use-my-address' onClick={() => setMyAdress()}>
                    Use my address
                  </div>
                ) : null}
              </div>
              {approvedToken ? (
                <div className='button-swap'>
                  <button
                    disabled={disableBtnSwap}
                    id='swap-button'
                    className='swap-button'
                    onClick={() => swap()}
                  >
                    <div className='css-10ob8xa'>
                      <SwapOutlined /> Swap Anyway
                    </div>
                  </button>
                </div>
              ) : (
                <div className='button-swap'>
                  <button
                    disabled={disableBtnSwap}
                    id='swap-button'
                    className='swap-button'
                    onClick={() => approveToken()}
                  >
                    <div className='css-10ob8xa'>
                      <SwapOutlined /> Approve
                    </div>
                  </button>
                </div>
              )}
            </div>
          </div>
        </Col>
        <Col xs={{ order: 1, span: 24 }} md={{ order: 2, span: 12 }}>
          <div className='area-chain-to-chain'>
            <Row justify='space-between'>
              <Col
                span={11}
                className={
                  'eth-to-harmony button-chain ' +
                  (indexRoadSwap === 0 ? 'enable-road' : 'disable-road')
                }
                onClick={() => {
                  chooseRoad('ETH', 'ONE', 0);
                }}
              >
                <Row>
                  <Col span={10}>
                    <div className='token-source'>
                      <h3>
                        <img alt='icon' src='/eth.svg' className='icon-token' />
                        ETH
                      </h3>
                    </div>
                  </Col>
                  <Col span={4}>
                    <div className='symbol-arrow'>
                      <RightOutlined />
                    </div>
                  </Col>
                  <Col span={10}>
                    <div className='token-dest'>
                      <h3>
                        <img alt='icon' src='/one.svg' className='icon-token' /> ONE
                      </h3>
                    </div>
                  </Col>
                </Row>
              </Col>
              <Col
                span={11}
                className={
                  'harmony-to-eth button-chain ' +
                  (indexRoadSwap === 1 ? 'enable-road' : 'disable-road')
                }
                onClick={() => {
                  chooseRoad('ONE', 'ETH', 1);
                }}
              >
                <Row>
                  <Col span={10}>
                    <div className='token-source'>
                      <h3>
                        <img alt='icon' src='/one.svg' className='icon-token' /> ONE
                      </h3>
                    </div>
                  </Col>
                  <Col span={4}>
                    <div className='symbol-arrow'>
                      <RightOutlined />
                    </div>
                  </Col>
                  <Col span={10}>
                    <div className='token-dest'>
                      <h3>
                        <img alt='icon' src='/eth.svg' className='icon-token' />
                        ETH
                      </h3>
                    </div>
                  </Col>
                </Row>
              </Col>
              <Col span={24} className='wallets-area'>
                <Row>
                  <Col span={11}>
                    <div className='wallet-source'>
                      <div className='text-source-dest'>
                        <h2>Wallet Source</h2>
                      </div>
                      {listWallets[walletSource]}
                    </div>
                  </Col>
                  <Col span={2} className='symbol-arrow-wallet'>
                    <ArrowRightOutlined />
                  </Col>
                  <Col span={11}>
                    <div className='wallet-dest'>
                      <div className='text-source-dest'>
                        <h2>Wallet Dest</h2>
                      </div>
                      {listWallets[walletDest]}
                    </div>
                  </Col>
                </Row>
              </Col>
            </Row>
          </div>
        </Col>
      </Row>
    </div>
  );
}

export default SwapPage;
