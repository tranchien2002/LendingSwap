pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20Detailed.sol";
import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";

contract TestERC20 is ERC20, ERC20Detailed, ERC20Mintable {
  constructor() ERC20Detailed("Test ERC20", "TestERC20", 18) ERC20Mintable() public {
  }
}