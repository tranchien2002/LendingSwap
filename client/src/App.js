import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import './App.css';

// import { connectOneWallet } from './utils/connectOneWallet';
// import { connectMetamask } from './utils/connectMetamask';
// import { convertToken } from 'utils/getTokenInfo';
import { Layout } from 'antd';

import Head from 'components/Head';
import Footer from 'components/Footer';
import SwapPage from 'pages/SwapPage';

function App() {
  // useEffect(() => {
  //   // connectOneWallet();
  //   // connectMetamask();
  //   let test = async () => {
  //     let price = await convertToken('BAND', 'USD', 100);
  //     console.log('price', price);
  //   };
  //   test();
  // });

  return (
    <Router>
      <Layout className='App'>
        <Head />
        <Switch>
          <Route exact path='/' component={SwapPage} />
        </Switch>
        <Footer />
      </Layout>
    </Router>
  );
}

export default App;
