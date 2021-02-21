pragma solidity ^0.5.0;

contract BridgeRegistry {
    address public ethereumBridge;
    address payable public bridgeBank;
    address public oracle;
    address public valset;
    address public operator;

    event HmyLogBridgeBankSet(address _bridgeBank);
    event HmyLogOracleSet(address _oracle);
    event HmyLogHarmonyBridgeSet(address _ethereumBridge);
    event HmyLogValsetSet(address _valset);

    modifier onlyOperator() {
        require(msg.sender == operator, "Must be the operator.");
        _;
    }

    constructor() public {
        operator = msg.sender;
    }

    function getOperator() public view returns (address) {
        return operator;
    }

    function setEthereumBridge(address _ethereumBridge) public onlyOperator {
        ethereumBridge = _ethereumBridge;
        emit HmyLogHarmonyBridgeSet(_ethereumBridge);
    }

    function getEthereumBridge() public view returns (address) {
        return ethereumBridge;
    }

    function setBridgeBank(address payable _bridgeBank) public onlyOperator {
        bridgeBank = _bridgeBank;
        emit HmyLogBridgeBankSet(address(bridgeBank));
    }

    function getBridgeBank() public view returns (address payable) {
        return bridgeBank;
    }

    function setOracle(address _oracle) public onlyOperator {
        oracle = _oracle;
        emit HmyLogOracleSet(oracle);
    }

    function getOracle() public view returns (address) {
        return oracle;
    }

    function setValset(address _valset) public onlyOperator {
        valset = _valset;
        emit HmyLogValsetSet(_valset);
    }

    function getValset() public view returns (address) {
        return valset;
    }
}
