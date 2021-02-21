pragma solidity ^0.5.0;

import "../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../node_modules/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "./libraries/openzeppelin-upgradeability/VersionedInitializable.sol";

import "./Valset.sol";
import "./BridgeBank/BridgeBank.sol";
import "./BridgeRegistry.sol";
import "./Oracle.sol";

contract HarmonyBridge is VersionedInitializable {
    using SafeMath for uint256;


    uint256 public constant HARMONYBRIDGE_REVISION = 0x1;
    uint256 public unlockClaimCount;
    address public operator;
    Valset public valset;
    BridgeBank public bridgeBank;
    BridgeRegistry public bridgeRegistry;
    Oracle public oracle;

    mapping(uint256 => UnlockClaim) public unlockClaims;

    enum Status {Null, Pending, Success, Failed}

    struct UnlockClaim {
        address harmonySender;
        address payable ethereumReceiver;
        address originalValidator;
        address token;
        uint256 amount;
        Status status;
    }

    event EthLogNewUnlockClaim(
        uint256 _unlockID,
        address _harmonySender,
        address payable _ethereumReceiver,
        address _validatorAddress,
        address _tokenAddress,
        uint256 _amount
    );

    event EthLogUnlockCompleted(uint256 _unlockID);

    modifier isPending(uint256 _unlockID) {
        require(isUnlockClaimActive(_unlockID), "Unlock claim is not active");
        _;
    }

    modifier onlyOperator() {
        require(msg.sender == operator, "Must be the operator.");
        _;
    }

    modifier isActive() {
        require(
            bridgeRegistry.getOracle() != address(0) && bridgeRegistry.getBridgeBank() != address(0),
            "The Operator must set the oracle and bridge bank for bridge activation"
        );
        _;
    }

    function initialize(
        address _bridgeRegistry
    ) public initializer {
        bridgeRegistry = BridgeRegistry(_bridgeRegistry);
        operator = bridgeRegistry.getOperator();
        valset = Valset(bridgeRegistry.getValset());
        bridgeBank = BridgeBank(bridgeRegistry.getBridgeBank());
        oracle = Oracle(bridgeRegistry.getOracle());
        unlockClaimCount = 0;
    }

    function getRevision() internal pure returns (uint256) {
        return HARMONYBRIDGE_REVISION;
    }

    function newUnlockClaim(
        address _harmonySender,
        address payable _ethereumReceiver,
        address _token,
        uint256 _amount
    ) public isActive {
        require(_amount > 0, "Amount token must be greater than zero");
        require(
            valset.isActiveValidator(msg.sender),
            "Must be an active validator"
        );
        require(
            bridgeBank.checkUnlockable(_token, _amount),
            "Not enough locked assets to complete the proposed prophecy"
        );

        // Create the new UnlockClaim
        UnlockClaim memory unlockClaim = UnlockClaim(
            _harmonySender,
            _ethereumReceiver,
            msg.sender,
            _token,
            _amount,
            Status.Pending
        );

        // Increment count and add the new UnlockClaim to the mapping
        unlockClaimCount = unlockClaimCount.add(1);
        unlockClaims[unlockClaimCount] = unlockClaim;

        emit EthLogNewUnlockClaim(
            unlockClaimCount,
            _harmonySender,
            _ethereumReceiver,
            msg.sender,
            _token,
            _amount
        );
    }

    function completeUnlockClaim(uint256 _unlockID)
        public
        isPending(_unlockID)
    {
        require(
            msg.sender == address(oracle),
            "Only the Oracle may complete prophecies"
        );

        unlockClaims[_unlockID].status = Status.Success;

        unlockTokens(_unlockID);

        emit EthLogUnlockCompleted(_unlockID);
    }

    function unlockTokens(uint256 _unlockID) internal {
        UnlockClaim memory unlockClaim = unlockClaims[_unlockID];

        if (unlockClaim.token == bridgeBank.ETHAddress()){
            bridgeBank.unlockETH(
                unlockClaim.ethereumReceiver,
                unlockClaim.amount
            );
        } else {
            bridgeBank.unlockERC20(
                unlockClaim.ethereumReceiver,
                unlockClaim.token,
                unlockClaim.amount
            );
        }
    }

    function isUnlockClaimActive(uint256 _unlockID) public view returns (bool) {
        return unlockClaims[_unlockID].status == Status.Pending;
    }

    function isUnlockClaimValidatorActive(uint256 _unlockID)
        public
        view
        returns (bool)
    {
        return
            valset.isActiveValidator(unlockClaims[_unlockID].originalValidator);
    }
}
