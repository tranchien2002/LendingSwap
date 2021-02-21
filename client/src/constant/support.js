import eth from 'assets/LogoTokens/eth.png';
import knc from 'assets/LogoTokens/knc.png';
import one from 'assets/LogoTokens/one.png';

import metamask from 'assets/LogoWallets/metamask.png';
import walletconnect from 'assets/LogoWallets/walletconnect.png';
import mathwallet from 'assets/LogoWallets/mathwallet.png';

export const tokens = [
  { symbol: 'ETH', icon: eth },
  { symbol: 'KNC', icon: knc },
  { symbol: 'ONE', icon: one }
];

export const wallets = {
  WalletConnect: { name: 'WalletConnect', icon: walletconnect },
  MetaMask: { name: 'MetaMask', icon: metamask },
  OneWallet: { name: 'One Wallet', icon: mathwallet }
};
