const contractAddress = {
  0: {
    bridgeBank: 'one17jwyxhsn4uukjyd32nnszhf0au03lz30j4yvwj'
  }
};

export const getHmyContractAddress = _chainId => {
  return contractAddress[_chainId];
};
