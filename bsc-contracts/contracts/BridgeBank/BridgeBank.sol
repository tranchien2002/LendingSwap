pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20Detailed.sol";
import "../../node_modules/openzeppelin-solidity/contracts/utils/ReentrancyGuard.sol";
import "../libraries/openzeppelin-upgradeability/VersionedInitializable.sol";
import "../interfaces/band-oracle/BandOracleInterface.sol";
import "../Oracle.sol";
import "../BridgeRegistry.sol";

contract BridgeBank is ReentrancyGuard, VersionedInitializable {
    using SafeMath for uint256;

    uint256 public constant BRIDGEBANK_REVISION = 0x1;
    uint256 public feeNumerator;
    uint256 public feeDenominator;
    uint256 public SAFE_NUMBER = 1e12;
    uint256 public lockNonce;

    address public operator;
    address public ETHAddress =
        address(0x1111111111111111111111111111111111111111);
    address public ONEAddress =
        address(0x2222222222222222222222222222222222222222);

    struct TokenData {
        uint256 lockedFund;
        address ethereumMappedToken;
        bool isActive;
    }

    mapping(address => TokenData) public tokensData;

    BridgeRegistry public bridgeRegistry;
    Oracle public oracle;
    EthereumBridge public ethereumBridge;
    BandOracleInterface public bandOracleInterface;

    event HmyUpdateOracle(address _newOracle);
    event HmyUpdateEthereumBridge(address _newEthereumBridge);
    event HmyUpdateFee(uint256 _feeNumerator, uint256 _feeDenominator);
    event HmyWithdrawONE(address _receiver, uint256 _amount);
    event HmyWithdrawERC20(address _token, address _receiver, uint256 _amount);
    event HmyLogLock(
        address _harmonySender,
        address _ethereumReceiver,
        address _harmonyToken,
        address _ethereumToken,
        uint256 _harmonyTokenAmount,
        uint256 _ethereumTokenAmount,
        uint256 _nonce
    );
    event HmyLogUnlock(address _to, address _token, uint256 _value);

    modifier availableNonce() {
        require(lockNonce + 1 > lockNonce, "No available nonces.");
        _;
    }

    modifier onlyOperator() {
        require(msg.sender == operator, "Must be BridgeBank operator.");
        _;
    }

    modifier onlyOracle() {
        require(
            msg.sender == address(oracle),
            "Access restricted to the oracle"
        );
        _;
    }

    modifier onlyEthereumBridge() {
        require(
            msg.sender == address(ethereumBridge),
            "Access restricted to the harmony bridge"
        );
        _;
    }

    modifier tokenMustBeActive(address _harmonyToken) {
        require(tokensData[_harmonyToken].isActive, "Token is not active");
        _;
    }

    modifier amountMustGreaterThanZero(uint256 _amount) {
        require(_amount > 0, "Amount token must be greater than zero");
        _;
    }

    modifier receiverMustBeValid(address _receiver) {
        require(_receiver != address(0), "Receiver must be valid");
        _;
    }

    modifier tokenAddressMustBeValid(address _token) {
        require(_token != address(0), "Token address must be valid");
        _;
    }

    modifier tokenMustNotBeONE(address _token) {
        require(_token != ONEAddress, "Token must not be ONE");
        _;
    }

    function initialize(
        address _bridgeRegistry,
        address _bandOracleAddress,
        address _ethereumONE
    ) public payable initializer {
        bridgeRegistry = BridgeRegistry(_bridgeRegistry);
        operator = bridgeRegistry.getOperator();
        oracle = Oracle(bridgeRegistry.getOracle());
        ethereumBridge = EthereumBridge(bridgeRegistry.getEthereumBridge());

        bandOracleInterface = BandOracleInterface(_bandOracleAddress);
        lockNonce = 0;

        addDataTokenInternal(tokensData[ONEAddress], msg.value, _ethereumONE);
    }

    function getRevision() internal pure returns (uint256) {
        return BRIDGEBANK_REVISION;
    }

    function() external payable onlyOperator {}

    function addToken(
        address _harmonyToken,
        uint256 _harmonyTokenAmount,
        address _ethereumToken
    ) public onlyOperator {
        require(
            _harmonyToken != address(0) && _ethereumToken != address(0),
            "Token address must be valid"
        );

        require(
            tokensData[_harmonyToken].ethereumMappedToken == address(0),
            "Token already added!"
        );

        IERC20(_harmonyToken).transferFrom(
            msg.sender,
            address(this),
            _harmonyTokenAmount
        );

        addDataTokenInternal(
            tokensData[_harmonyToken],
            _harmonyTokenAmount,
            _ethereumToken
        );
    }

    function addDataTokenInternal(
        TokenData storage data,
        uint256 _harmonyTokenAmount,
        address _ethereumToken
    ) internal {
        data.lockedFund = _harmonyTokenAmount;
        data.ethereumMappedToken = _ethereumToken;
        data.isActive = true;
    }

    function deactivateToken(address _harmonyToken)
        public
        onlyOperator
        tokenAddressMustBeValid(_harmonyToken)
    {
        require(
            tokensData[_harmonyToken].isActive == true,
            "Token already deactivated!"
        );
        tokensData[_harmonyToken].isActive = false;
    }

    function activateToken(address _harmonyToken)
        public
        onlyOperator
        tokenAddressMustBeValid(_harmonyToken)
    {
        require(
            tokensData[_harmonyToken].isActive == false,
            "Token already activated!"
        );
        tokensData[_harmonyToken].isActive = true;
    }

    function isActiveToken(address _harmonyToken) public view returns (bool) {
        return tokensData[_harmonyToken].isActive;
    }

    function getTokenMappedAddress(address _harmonyToken)
        public
        view
        returns (address)
    {
        return tokensData[_harmonyToken].ethereumMappedToken;
    }

    function getLockedFund(address _harmonyToken)
        public
        view
        returns (uint256)
    {
        return tokensData[_harmonyToken].lockedFund;
    }

    /*
     * @dev: Locks received Ethereum/ERC20 funds.
     *
     * @param _recipient: bytes representation of destination address.
     * @param _token: token address in origin chain (0x0 if ethereum)
     * @param _amount: value of deposit
     */
    function swapToken_1_1(
        address _ethereumReceiver,
        address _harmonyToken,
        uint256 _harmonyTokenAmount
    )
        public
        nonReentrant
        availableNonce
        tokenMustNotBeONE(_harmonyToken)
        tokenMustBeActive(_harmonyToken)
        amountMustGreaterThanZero(_harmonyTokenAmount)
        receiverMustBeValid(_ethereumReceiver)
    {
        uint256 fee = calculateFee(_harmonyTokenAmount);

        require(
            IERC20(_harmonyToken).transferFrom(
                msg.sender,
                address(this),
                _harmonyTokenAmount + fee
            ),
            "Contract token allowances insufficient to complete this lock request"
        );

        address _ethereumToken = tokensData[_harmonyToken].ethereumMappedToken;

        updateOnLock(
            msg.sender,
            _ethereumReceiver,
            _harmonyToken,
            _ethereumToken,
            _harmonyTokenAmount,
            _harmonyTokenAmount
        );
    }

    function swapONEForETH(address _ethereumReceiver, uint256 _amountONE)
        public
        payable
        nonReentrant
        availableNonce
        amountMustGreaterThanZero(_amountONE)
        receiverMustBeValid(_ethereumReceiver)
    {
        uint256 fee = calculateFee(_amountONE);

        require(
            msg.value == _amountONE + fee,
            "The transactions value must be equal the specified amount (in wei)"
        );

        BandOracleInterface.ReferenceData memory data =
            bandOracleInterface.getReferenceData("ONE", "ETH");

        uint256 amountETH = _amountONE.mul(data.rate).div(1e18);

        updateOnLock(
            msg.sender,
            _ethereumReceiver,
            ONEAddress,
            ETHAddress,
            _amountONE,
            amountETH
        );
    }

    function swapONEForWrappedONE(address _ethereumReceiver, uint256 _amountONE)
        public
        payable
        nonReentrant
        availableNonce
        tokenMustBeActive(ONEAddress)
        amountMustGreaterThanZero(_amountONE)
        receiverMustBeValid(_ethereumReceiver)
    {
        uint256 fee = calculateFee(_amountONE);

        require(
            msg.value == _amountONE + fee,
            "The transactions value must be equal the specified amount (in wei)"
        );

        address _ethereumWONE = tokensData[ONEAddress].ethereumMappedToken;

        updateOnLock(
            msg.sender,
            _ethereumReceiver,
            ONEAddress,
            _ethereumWONE,
            _amountONE,
            _amountONE
        );
    }

    function swapONEForToken(
        address _ethereumReceiver,
        uint256 _amountONE,
        address _destToken
    )
        public
        payable
        nonReentrant
        availableNonce
        tokenMustBeActive(ONEAddress)
        tokenMustBeActive(_destToken)
        amountMustGreaterThanZero(_amountONE)
        receiverMustBeValid(_ethereumReceiver)
    {
        uint256 fee = calculateFee(_amountONE);

        require(
            msg.value == _amountONE + fee,
            "The transactions value must be equal the specified amount (in wei)"
        );

        BandOracleInterface.ReferenceData memory data =
            bandOracleInterface.getReferenceData(
                "ONE",
                ERC20Detailed(_destToken).symbol()
            );

        uint256 ethereumTokenAmount = _amountONE.mul(data.rate).div(1e18);

        address ethereumToken = tokensData[_destToken].ethereumMappedToken;

        updateOnLock(
            msg.sender,
            _ethereumReceiver,
            ONEAddress,
            ethereumToken,
            _amountONE,
            ethereumTokenAmount
        );
    }

    function swapTokenForToken(
        address _ethereumReceiver,
        address _harmonyToken,
        uint256 _harmonyTokenAmount,
        address _destToken
    )
        public
        availableNonce
        nonReentrant
        tokenMustBeActive(_harmonyToken)
        tokenMustBeActive(_destToken)
        amountMustGreaterThanZero(_harmonyTokenAmount)
        receiverMustBeValid(_ethereumReceiver)
    {
        uint256 fee = calculateFee(_harmonyTokenAmount);

        require(
            IERC20(_harmonyToken).transferFrom(
                msg.sender,
                address(this),
                _harmonyTokenAmount + fee
            ),
            "Contract token allowances insufficient to complete this lock request"
        );

        BandOracleInterface.ReferenceData memory data =
            bandOracleInterface.getReferenceData(
                ERC20Detailed(_harmonyToken).symbol(),
                ERC20Detailed(_destToken).symbol()
            );
        uint256 ethereumTokenAmount =
            _harmonyTokenAmount.mul(data.rate).div(1e18);

        address ethereumToken = tokensData[_destToken].ethereumMappedToken;

        updateOnLock(
            msg.sender,
            _ethereumReceiver,
            _harmonyToken,
            ethereumToken,
            _harmonyTokenAmount,
            ethereumTokenAmount
        );
    }

    function swapTokenForWrappedONE(
        address _ethereumReceiver,
        address _harmonyToken,
        uint256 _harmonyTokenAmount
    )
        public
        availableNonce
        nonReentrant
        tokenMustBeActive(_harmonyToken)
        amountMustGreaterThanZero(_harmonyTokenAmount)
        receiverMustBeValid(_ethereumReceiver)
    {
        uint256 fee = calculateFee(_harmonyTokenAmount);

        require(
            IERC20(_harmonyToken).transferFrom(
                msg.sender,
                address(this),
                _harmonyTokenAmount + fee
            ),
            "Contract token allowances insufficient to complete this lock request"
        );

        BandOracleInterface.ReferenceData memory data =
            bandOracleInterface.getReferenceData(
                ERC20Detailed(_harmonyToken).symbol(),
                "ONE"
            );
        uint256 amountWONE = _harmonyTokenAmount.mul(data.rate).div(1e18);

        address ethereumWONE = tokensData[ONEAddress].ethereumMappedToken;

        updateOnLock(
            msg.sender,
            _ethereumReceiver,
            _harmonyToken,
            ethereumWONE,
            _harmonyTokenAmount,
            amountWONE
        );
    }

    function swapTokenForETH(
        address _ethereumReceiver,
        address _harmonyToken,
        uint256 _harmonyTokenAmount
    )
        public
        availableNonce
        nonReentrant
        tokenMustBeActive(_harmonyToken)
        amountMustGreaterThanZero(_harmonyTokenAmount)
        receiverMustBeValid(_ethereumReceiver)
    {
        uint256 fee = calculateFee(_harmonyTokenAmount);

        require(
            IERC20(_harmonyToken).transferFrom(
                msg.sender,
                address(this),
                _harmonyTokenAmount + fee
            ),
            "Contract token allowances insufficient to complete this lock request"
        );

        BandOracleInterface.ReferenceData memory data =
            bandOracleInterface.getReferenceData(
                ERC20Detailed(_harmonyToken).symbol(),
                "ETH"
            );
        uint256 amountETH = _harmonyTokenAmount.mul(data.rate).div(1e18);

        updateOnLock(
            msg.sender,
            _ethereumReceiver,
            _harmonyToken,
            ETHAddress,
            _harmonyTokenAmount,
            amountETH
        );
    }

    function unlockERC20(
        address payable _harmonyReceiver,
        address _harmonyToken,
        uint256 _harmonyTokenAmount
    )
        public
        nonReentrant
        onlyEthereumBridge
        amountMustGreaterThanZero(_harmonyTokenAmount)
        receiverMustBeValid(_harmonyReceiver)
    {
        require(
            tokensData[_harmonyToken].ethereumMappedToken != address(0),
            "Invalid token address"
        );

        require(
            _harmonyTokenAmount <= getTotalERC20Balance(_harmonyToken),
            "Exceeded amount of Token allowed to withdraw"
        );

        IERC20(_harmonyToken).transfer(_harmonyReceiver, _harmonyTokenAmount);

        updateOnUnlock(_harmonyToken, _harmonyReceiver, _harmonyTokenAmount);
    }

    function unlockONE(address payable _harmonyReceiver, uint256 _amountONE)
        public
        onlyEthereumBridge
        nonReentrant
        amountMustGreaterThanZero(_amountONE)
        receiverMustBeValid(_harmonyReceiver)
    {
        require(
            _amountONE <= getTotalONEBalance(),
            "Exceeded amount of ETH allowed to withdraw"
        );
        _harmonyReceiver.transfer(_amountONE);

        updateOnUnlock(ONEAddress, _harmonyReceiver, _amountONE);
    }

    function updateOracle(address _oracleAddress) public onlyOperator {
        oracle = Oracle(_oracleAddress);
        emit HmyUpdateOracle(_oracleAddress);
    }

    function updateHmyBridge(address _ethereumBridge) public onlyOperator {
        ethereumBridge = EthereumBridge(_ethereumBridge);
        emit HmyUpdateEthereumBridge(_ethereumBridge);
    }

    function updateFee(uint256 _feeNumerator, uint256 _feeDenominator)
        public
        onlyOperator
    {
        feeNumerator = _feeNumerator;
        feeDenominator = _feeDenominator;
        emit HmyUpdateFee(_feeNumerator, _feeDenominator);
    }

    function withdrawONE(address payable _harmonyReceiver, uint256 _amountONE)
        public
        onlyOperator
        nonReentrant
    {
        require(
            _amountONE <=
                getTotalONEBalance() - tokensData[ONEAddress].lockedFund,
            "Exceeded amount of ONE allowed to withdraw"
        );

        _harmonyReceiver.transfer(_amountONE);
        emit HmyWithdrawONE(_harmonyReceiver, _amountONE);
    }

    function withdrawERC20(
        address _harmonyToken,
        address _harmonyReceiver,
        uint256 _harmonyTokenAmount
    ) public onlyOperator nonReentrant {
        require(
            _harmonyTokenAmount <=
                getTotalERC20Balance(_harmonyToken) -
                    tokensData[_harmonyToken].lockedFund,
            "Exceeded amount of Token allowed to withdraw"
        );

        IERC20(_harmonyToken).transfer(_harmonyReceiver, _harmonyTokenAmount);
        emit HmyWithdrawERC20(
            _harmonyToken,
            _harmonyReceiver,
            _harmonyTokenAmount
        );
    }

    function getTotalONEBalance() public view returns (uint256) {
        return address(this).balance;
    }

    function getTotalERC20Balance(address _harmonyToken)
        public
        view
        returns (uint256)
    {
        return IERC20(_harmonyToken).balanceOf(address(this));
    }

    function calculateFee(uint256 _amountToken)
        internal
        view
        returns (uint256)
    {
        uint256 fee;

        if (feeNumerator != 0 && feeDenominator != 0) {
            fee = _amountToken
                .mul(feeNumerator)
                .mul(SAFE_NUMBER)
                .div(feeDenominator)
                .div(SAFE_NUMBER);
        }

        return fee;
    }

    function updateOnLock(
        address _harmonySender,
        address _ethereumReceiver,
        address _harmonyToken,
        address _ethereumToken,
        uint256 _harmonyTokenAmount,
        uint256 _ethereumTokenAmount
    ) internal {
        require(_ethereumTokenAmount > 0, "Amount token must be greater than zero");

        lockNonce = lockNonce.add(1);

        tokensData[_ethereumToken].lockedFund = tokensData[_ethereumToken]
            .lockedFund
            .add(_ethereumTokenAmount);

        emit HmyLogLock(
            _harmonySender,
            _ethereumReceiver,
            _harmonyToken,
            _ethereumToken,
            _harmonyTokenAmount,
            _ethereumTokenAmount,
            lockNonce
        );
    }

    function updateOnUnlock(
        address _harmonyToken,
        address _harmonyReceiver,
        uint256 _harmonyTokenAmount
    ) internal {
        if(tokensData[_harmonyToken].lockedFund >= _harmonyTokenAmount){
            tokensData[_harmonyToken].lockedFund = tokensData[_harmonyToken]
                .lockedFund
                .sub(_harmonyTokenAmount);
        } else {
            tokensData[_harmonyToken].lockedFund = 0;
        }
        emit HmyLogUnlock(_harmonyReceiver, _harmonyToken, _harmonyTokenAmount);
    }

    function checkUnlockable(address _token, uint256 _amount) public view returns (bool) {
        if (_token == ONEAddress) {
            return getTotalONEBalance() >= _amount;
        } else {
            return getTotalERC20Balance(_token) >= _amount;
        }
    }
}
