import store from 'store';
import { setAddress, setHmy, setIsAuthorized } from 'store/hmy/action.js';

export const connectOneWallet = async () => {
  let isMathWallet = window.harmony && window.harmony.isMathWallet;
  if (isMathWallet) {
    let mathwallet = window.harmony;
    mathwallet.getAccount().then(async account => {
      store.dispatch(setAddress(account.address));
      // store.dispatch(setHmy(onewallet));
      store.dispatch(setIsAuthorized(true));
      localStorage.setItem(
        'harmony_hmy_session',
        JSON.stringify({
          address: account.address
        })
      );
      console.log('Wallet harmony connected');
    });
  } else {
    alert('Please connect Onewallet extension!');
  }
};

export const signOut = async () => {
  const { isAuthorized } = store.getState().hmy;
  if (isAuthorized) {
    let mathwallet = window.harmony;
    return mathwallet
      .forgetIdentity()
      .then(() => {
        store.dispatch(setIsAuthorized(false));
        store.dispatch(setAddress(''));
        store.dispatch(setHmy(null));
        localStorage.setItem(
          'harmony_hmy_session',
          JSON.stringify({
            address: ''
          })
        );
        console.log('Wallet harmony is sign out');
        return Promise.resolve();
      })
      .catch(err => {
        console.error(err.message);
      });
  }
};
