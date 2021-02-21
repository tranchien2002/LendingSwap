pragma solidity ^0.5.0;

contract BridgeRegistry {
    address public harmonyBridge;
    address payable public bridgeBank;
    address public oracle;
    address public valset;
    address public operator;

    event EthLogBridgeBankSet(address _bridgeBank);
    event EthLogOracleSet(address _oracle);
    event EthLogHarmonyBridgeSet(address _harmonyBridge);
    event EthLogValsetSet(address _valset);

    modifier onlyOperator() {
        require(msg.sender == operator, "Must be the operator.");
        _;
    }

    constructor() public {
        operator = msg.sender;
    }

    function getOperator() public view returns(address) {
        return operator;
    }

    function setHarmonyBridge(address _harmonyBridge) public onlyOperator {
        harmonyBridge = _harmonyBridge;
        emit EthLogHarmonyBridgeSet(_harmonyBridge);
    }

    function getHarmonyBridge() public view returns(address) {
        return harmonyBridge;
    }

    function setBridgeBank(address payable _bridgeBank) public onlyOperator {
        bridgeBank = _bridgeBank;
        emit EthLogBridgeBankSet(address(bridgeBank));
    }

    function getBridgeBank() public view returns(address payable) {
        return bridgeBank;
    }

    function setOracle(address _oracle) public onlyOperator {
        oracle = _oracle;
        emit EthLogOracleSet(oracle);
    }

    function getOracle() public view returns(address) {
        return oracle;
    }

     function setValset(address _valset) public onlyOperator {
        valset = _valset;
        emit EthLogValsetSet(_valset);
    }

    function getValset() public view returns(address) {
        return valset;
    }
}
