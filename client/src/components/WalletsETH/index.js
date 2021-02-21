import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { wallets } from 'constant/support';
import { connectWalletConnect } from 'store/eth/action';
import { connectMetamask, signOut } from 'utils/connectMetamask';
import { Row, Col, Button } from 'antd';
import { PoweroffOutlined } from '@ant-design/icons';

import './style.scss';
import '../stylecommon.scss';

function WalletsETH() {
  const dispatch = useDispatch();
  const address = useSelector(state => state.eth.address);
  const [sizeMobile, setSizeMobile] = useState(false);

  useEffect(() => {
    if (
      /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
    ) {
      setSizeMobile(true);
    } else {
      setSizeMobile(false);
    }
  }, [sizeMobile]);

  useEffect(() => {
    // get session ETH in localStorage
    const sessionETH = localStorage.getItem('harmony_eth_session');
    const sessionObjETH = JSON.parse(sessionETH);
    if (sessionObjETH && sessionObjETH.address) {
      connectMetamask();
    }
  });

  function connectWalletMobile() {
    dispatch(connectWalletConnect());
  }

  return (
    <Row justify='center' className='list-wallets wallets-eth'>
      {!address ? (
        sizeMobile ? (
          <Col
            onClick={() => {
              connectWalletMobile();
            }}
          >
            <div className='info-wallet'>
              <img alt='icon' className='icon-wallet' src={wallets['WalletConnect'].icon} />
              <h4>{wallets['WalletConnect'].name}</h4>
            </div>
          </Col>
        ) : (
          <Col
            onClick={() => {
              connectMetamask();
            }}
          >
            <div className='info-wallet'>
              <img alt='icon' className='icon-wallet' src={wallets['MetaMask'].icon} />
              <h4>{wallets['MetaMask'].name}</h4>
            </div>
          </Col>
        )
      ) : (
        <div>
          <div className='display-address-connected'>
            <p className='address-short'>
              {address && address.length > 0
                ? address.substring(0, 6) +
                  '...' +
                  address.substring(address.length - 5, address.length - 1)
                : address}
            </p>
            <Button
              type='primary'
              shape='circle'
              icon={<PoweroffOutlined />}
              onClick={() => {
                signOut();
              }}
            />
          </div>
        </div>
      )}
    </Row>
  );
}

export default WalletsETH;
