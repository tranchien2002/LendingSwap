pragma solidity ^0.5.0;

interface IBridgeRegistry {
    function getOperator() external view returns (address);

    function getEthereumBridge() external view returns (address);

    function getBridgeBank() external view returns (address payable);

    function getOracle() external view returns (address);

    function getValset() external view returns (address);
}
