import * as React from 'react';

import { Row, Col } from 'antd';

import './style.scss';

export default function Footer() {
  return (
    <div className='footer'>
      <div className='container'>
        <Row justify='start'>
          <Col xs={{ order: 2, span: 24 }} md={{ order: 1, span: 12 }} className='box-left'>
            <div className='footer-copyright'>
              <h3>Copyright © 2020 by PhoneFarm Team</h3>
            </div>
          </Col>
          <Col xs={{ order: 1, span: 24 }} md={{ order: 2, span: 12 }} className='box-right'>
            <Row className='icons icons-link' justify='end'>
              <Col xs={{ span: 2 }} md={{ span: 1 }} className='icon'>
                <a target='_blank' rel='noreferrer' href='https://phonefarm.finance/'>
                  <img src='https://phonefarm.finance/favicon.ico' alt='' width='36px' />
                </a>
              </Col>
              <Col xs={{ span: 2 }} md={{ span: 1 }} className='icon'>
                <a target='_blank' rel='noreferrer' href='https://t.me/phonefarm_official'>
                  <img
                    src='https://phonefarm.finance/static/media/telegram.a451f456.svg'
                    alt=''
                    width='36px'
                  />
                </a>
              </Col>
              <Col xs={{ span: 2 }} md={{ span: 1 }} className='icon'>
                <a target='_blank' rel='noreferrer' href='https://discord.com/invite/aBApkPx'>
                  <img
                    src='https://phonefarm.finance/static/media/discord.f10eba06.svg'
                    alt=''
                    width='36px'
                  />
                </a>
              </Col>
              <Col xs={{ span: 2 }} md={{ span: 1 }} className='icon'>
                <a target='_blank' rel='noreferrer' href='https://twitter.com/PhonefarmF'>
                  <img
                    src='https://phonefarm.finance/static/media/twitter.dd90206d.svg'
                    alt=''
                    width='36px'
                  />
                </a>
              </Col>
              <Col xs={{ span: 2 }} md={{ span: 1 }} className='icon'>
                <a target='_blank' rel='noreferrer' href='https://phonefarm-finance.medium.com'>
                  <img
                    src='https://phonefarm.finance/static/media/medium.965e9335.svg'
                    alt=''
                    width='36px'
                  />
                </a>
              </Col>
              <Col xs={{ span: 2 }} md={{ span: 1 }} className='icon'>
                <a target='_blank' rel='noreferrer' href='https://github.com/PhoneFarm-Project'>
                  <img
                    src='https://phonefarm.finance/static/media/github.1f33e75d.svg'
                    alt=''
                    width='36px'
                  />
                </a>
              </Col>
            </Row>
          </Col>
        </Row>
      </div>
    </div>
  );
}
