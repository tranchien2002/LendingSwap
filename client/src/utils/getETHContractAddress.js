const contractAddress = {
  42: {
    bridgeBank: '0x17F6A3398b47D7Bdf5cBfa88e36A6FaEEa71f1e6'
  }
};

export const getETHContractAddress = _chainId => {
  return contractAddress[_chainId];
};
