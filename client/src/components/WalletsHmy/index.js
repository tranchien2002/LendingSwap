import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { wallets } from 'constant/support';
import { setHmy } from 'store/hmy/action';
import { connectOneWallet, signOut } from 'utils/connectOneWallet';
import { Row, Col, Button } from 'antd';
import { PoweroffOutlined } from '@ant-design/icons';

function WalletsHmy() {
  const dispatch = useDispatch();
  const address = useSelector(state => state.hmy.address);
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
    // get session Hmy in localStorage
    setTimeout(() => {
      const sessionHmy = localStorage.getItem('harmony_hmy_session');
      const sessionObjHmy = JSON.parse(sessionHmy);
      if (sessionObjHmy && sessionObjHmy.address) {
        connectOneWallet();
      }
    }, 500);
  });

  function connectWalletConnect() {
    dispatch(setHmy());
  }
  return (
    <Row justify='center' className='list-wallets wallets-eth'>
      {!address ? (
        sizeMobile ? (
          <Col
            onClick={() => {
              connectWalletConnect();
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
              connectOneWallet();
            }}
          >
            <div className='info-wallet'>
              <img alt='icon' className='icon-wallet' src={wallets['OneWallet'].icon} />
              <h4>{wallets['OneWallet'].name}</h4>
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

export default WalletsHmy;
