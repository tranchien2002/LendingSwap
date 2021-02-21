import * as React from 'react';

import { Layout } from 'antd';

import './style.scss';

const { Header } = Layout;

export default function Head() {
  return (
    <Header className='header'>
      <div className='container'>
        <h1>
          {/* <img alt='logo' src='/logo/logo-stand-1.png' className='style-logo' /> */}
          {/* <img alt='logo' src='/logo/logo-stand-2.png' className='style-logo' /> */}
          <img alt='logo' src='/logo/logo-no-stand-1.gif' className='style-logo' />
          {/* <img alt='logo' src='/logo/logo-no-stand-2.gif' className='style-logo' /> */}
          {/* <img alt='logo' src='/logo/logo-no-stand-1.png' className='style-logo' /> */}
          {/* <img alt='logo' src='/logo/logo-no-stand-2.png' className='style-logo' /> */}
        </h1>
      </div>
    </Header>
  );
}
