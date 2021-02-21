pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20Detailed.sol";
import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";

contract WrappedETH is ERC20, ERC20Detailed, ERC20Mintable {
  constructor() ERC20Detailed("Wrapped ETH", "ETH", 18) ERC20Mintable() public {
  }
}