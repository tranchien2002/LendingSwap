pragma solidity ^0.5.0;

interface IWETHGateway {
    function depositETH(address onBehalfOf, uint16 referralCode)
        external
        payable;

    function withdrawETH(uint256 amount, address to) external;
}
