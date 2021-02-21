pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20Detailed.sol";
import "../../node_modules/openzeppelin-solidity/contracts/utils/ReentrancyGuard.sol";
import "../libraries/openzeppelin-upgradeability/VersionedInitializable.sol";
import "../interfaces/aave-protocol-v2/ILendingPool.sol";
import "../interfaces/aave-protocol-v2/IWETHGateway.sol";
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
    address public WETH;
    address public ETHAddress = address(
        0x1111111111111111111111111111111111111111
    );
    address public ONEAddress = address(
        0x2222222222222222222222222222222222222222
    );

    struct TokenData {
        uint256 lockedFund;
        address harmonyMappedToken;
        bool isActive;
    }

    mapping (address => TokenData) public tokensData;

    BridgeRegistry public bridgeRegistry;
    Oracle public oracle;
    HarmonyBridge public harmonyBridge;
    BandOracleInterface public bandOracleInterface;
    ILendingPool public lendingPool;
    IWETHGateway public wethGateway;

    event EthUpdateOracle(address _newOracle);
    event EthUpdateHarmonyBridge(address _newHarmonyBridge);
    event EthUpdateFee(uint256 _feeNumerator, uint256 _feeDenominator);
    event EthWithdrawETH(address _receiver, uint256 _amount);
    event EthWithdrawERC20(address _token, address _receiver, uint256 _amount);
    event EthLogLock(
        address _ethereumSender,
        address _harmonyReceiver,
        address _ethereumToken,
        address _harmonyToken,
        uint256 _ethereumTokenAmount,
        uint256 _harmonyTokenAmount,
        uint256 _nonce
    );
    event EthLogUnlock(address _to, address _token, uint256 _value);

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

    modifier onlyHarmonyBridge() {
        require(
            msg.sender == address(harmonyBridge),
            "Access restricted to the harmony bridge"
        );
        _;
    }

    modifier tokenMustBeActive(address _ethereumToken) {
        require(tokensData[_ethereumToken].isActive, "Token is not active");
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

    modifier tokenMustNotBeETH(address _token) {
        require(_token != ETHAddress, "Token must not be ETH");
        _;
    }

    function initialize(
        address _bridgeRegistry,
        address _bandOracleAddress,
        address _lendingPool,
        address _wethGateway,
        address _weth,
        address _harmonyETH
    ) payable public initializer {
        bridgeRegistry = BridgeRegistry(_bridgeRegistry);
        operator = bridgeRegistry.getOperator();
        oracle = Oracle(bridgeRegistry.getOracle());
        harmonyBridge = HarmonyBridge(bridgeRegistry.getHarmonyBridge());

        bandOracleInterface = BandOracleInterface(_bandOracleAddress);
        lendingPool = ILendingPool(_lendingPool);
        wethGateway = IWETHGateway(_wethGateway);
        WETH = _weth;
        lockNonce = 0;

        addDataTokenInternal(tokensData[ETHAddress], msg.value, _harmonyETH);
        DataTypes.ReserveData memory reserve = lendingPool.getReserveData(WETH);
        address aToken = reserve.aTokenAddress;
        IERC20(aToken).approve(address(wethGateway), uint256(-1));
    }

    function getRevision() internal pure returns (uint256) {
        return BRIDGEBANK_REVISION;
    }

    function() external payable onlyOperator {}

    function addToken(address _ethereumToken, uint256 _ethereumTokenAmount, address _harmonyToken) public onlyOperator {
        require(
            _ethereumToken != address(0) && _harmonyToken != address(0),
            "Token address must be valid"
        );

        require(tokensData[_ethereumToken].harmonyMappedToken == address(0), "Token already added!");

        IERC20(_ethereumToken).transferFrom(msg.sender, address(this), _ethereumTokenAmount);

        addDataTokenInternal(tokensData[_ethereumToken], _ethereumTokenAmount, _harmonyToken);

        IERC20(_ethereumToken).approve(
            address(lendingPool),
            uint256(-1)
        );
    }

    function addDataTokenInternal(TokenData storage data, uint256 _ethereumTokenAmount , address _harmonyToken)
        internal
    {
        data.lockedFund = _ethereumTokenAmount;
        data.harmonyMappedToken = _harmonyToken;
        data.isActive = true;
    }

    function deactivateToken(address _ethereumToken) public onlyOperator tokenAddressMustBeValid(_ethereumToken) {
        require(tokensData[_ethereumToken].isActive == true, "Token already deactivated!");
        tokensData[_ethereumToken].isActive = false;
    }

    function activateToken(address _ethereumToken) public onlyOperator tokenAddressMustBeValid(_ethereumToken) {
        require(tokensData[_ethereumToken].isActive == false, "Token already activated!");
        tokensData[_ethereumToken].isActive = true;
    }

    function isActiveToken(address _ethereumToken) public view returns (bool) {
        return tokensData[_ethereumToken].isActive;
    }

    function getTokenMappedAddress(address _ethereumToken)
        public
        view
        returns (address)
    {
        return tokensData[_ethereumToken].harmonyMappedToken;
    }

    function getLockedFund(address _ethereumToken) public view returns (uint256) {
        return tokensData[_ethereumToken].lockedFund;
    }

    /*
     * @dev: Locks received Ethereum/ERC20 funds.
     *
     * @param _recipient: bytes representation of destination address.
     * @param _token: token address in origin chain (0x0 if ethereum)
     * @param _amount: value of deposit
     */
    function swapToken_1_1(
        address _harmonyReceiver,
        address _ethereumToken,
        uint256 _ethereumTokenAmount
    ) public
        nonReentrant
        availableNonce
        tokenMustNotBeETH(_ethereumToken)
        tokenMustBeActive(_ethereumToken)
        amountMustGreaterThanZero(_ethereumTokenAmount)
        receiverMustBeValid(_harmonyReceiver)
    {
        uint256 fee = calculateFee(_ethereumTokenAmount);

        require(
            IERC20(_ethereumToken).transferFrom(
                msg.sender,
                address(this),
                _ethereumTokenAmount + fee
            ),
            "Contract token allowances insufficient to complete this lock request"
        );

        DataTypes.ReserveData memory reserve = lendingPool.getReserveData(_ethereumToken);
        if(reserve.aTokenAddress != address(0)){
            lendingPool.deposit(
                _ethereumToken,
                _ethereumTokenAmount + fee,
                address(this),
                0
            );
        }

        address _harmonyToken = tokensData[_ethereumToken].harmonyMappedToken;

        updateOnLock(msg.sender, _harmonyReceiver, _ethereumToken, _harmonyToken, _ethereumTokenAmount, _ethereumTokenAmount);
    }

    function swapETHForONE(address _harmonyReceiver, uint256 _amountETH)
        public
        payable
        nonReentrant
        availableNonce
        tokenMustBeActive(ETHAddress)
        amountMustGreaterThanZero(_amountETH)
        receiverMustBeValid(_harmonyReceiver)
    {
        uint256 fee = calculateFee(_amountETH);

        require(
            msg.value == _amountETH + fee,
            "The transactions value must be equal the specified amount (in wei)"
        );

        (bool success, ) = address(wethGateway).call.value(msg.value)(
            abi.encodeWithSignature(
                "depositETH(address,uint16)",
                address(this),
                0
            )
        );

        require(success, "Aave LendingPool: deposit ETH failed!");

        BandOracleInterface.ReferenceData memory data = bandOracleInterface
            .getReferenceData("ETH", "ONE");

        uint256 amountONE = _amountETH.mul(data.rate).div(1e18);

        updateOnLock(msg.sender, _harmonyReceiver, ETHAddress, ONEAddress, _amountETH, amountONE);
    }

    function swapETHForWrappedETH(address _harmonyReceiver, uint256 _amountETH)
        public
        payable
        nonReentrant
        availableNonce
        tokenMustBeActive(ETHAddress)
        amountMustGreaterThanZero(_amountETH)
        receiverMustBeValid(_harmonyReceiver)
    {
        uint256 fee = calculateFee(_amountETH);

        require(
            msg.value == _amountETH + fee,
            "The transactions value must be equal the specified amount (in wei)"
        );

        (bool success, ) = address(wethGateway).call.value(msg.value)(
            abi.encodeWithSignature(
                "depositETH(address,uint16)",
                address(this),
                0
            )
        );

        require(success, "Aave LendingPool: deposit ETH failed!");

        address _harmonyWETH = tokensData[ETHAddress].harmonyMappedToken;

        updateOnLock(msg.sender, _harmonyReceiver, ETHAddress, _harmonyWETH, _amountETH, _amountETH);
    }

    function swapETHForToken(
        address _harmonyReceiver,
        uint256 _amountETH,
        address _destToken
    ) public payable
        nonReentrant
        availableNonce
        tokenMustBeActive(ETHAddress)
        tokenMustBeActive(_destToken)
        amountMustGreaterThanZero(_amountETH)
        receiverMustBeValid(_harmonyReceiver)
    {
        uint256 fee = calculateFee(_amountETH);

        require(
            msg.value == _amountETH + fee,
            "The transactions value must be equal the specified amount (in wei)"
        );

        (bool success, ) = address(wethGateway).call.value(msg.value)(
            abi.encodeWithSignature(
                "depositETH(address,uint16)",
                address(this),
                0
            )
        );

        require(success, "Aave LendingPool: deposit ETH failed!");

        BandOracleInterface.ReferenceData memory data = bandOracleInterface
            .getReferenceData("ETH", ERC20Detailed(_destToken).symbol());

        uint256 harmonyTokenAmount = _amountETH.mul(data.rate).div(1e18);

        address harmonyToken = tokensData[_destToken].harmonyMappedToken;

        updateOnLock(msg.sender, _harmonyReceiver, ETHAddress, harmonyToken, _amountETH, harmonyTokenAmount);
    }

    function swapTokenForToken(
        address _harmonyReceiver,
        address _ethereumToken,
        uint256 _ethereumTokenAmount,
        address _destToken
    ) public
        availableNonce
        nonReentrant
        tokenMustBeActive(_ethereumToken)
        tokenMustBeActive(_destToken)
        amountMustGreaterThanZero(_ethereumTokenAmount)
        receiverMustBeValid(_harmonyReceiver)
    {
        uint256 fee = calculateFee(_ethereumTokenAmount);

        require(
            IERC20(_ethereumToken).transferFrom(
                msg.sender,
                address(this),
                _ethereumTokenAmount + fee
            ),
            "Contract token allowances insufficient to complete this lock request"
        );

        DataTypes.ReserveData memory reserve = lendingPool.getReserveData(_ethereumToken);
        if(reserve.aTokenAddress != address(0)){
            lendingPool.deposit(
                _ethereumToken,
                _ethereumTokenAmount + fee,
                address(this),
                0
            );
        }

        BandOracleInterface.ReferenceData memory data = bandOracleInterface
            .getReferenceData(ERC20Detailed(_ethereumToken).symbol(), ERC20Detailed(_destToken).symbol());
        uint256 harmonyTokenAmount = _ethereumTokenAmount.mul(data.rate).div(1e18);

        address harmonyToken = tokensData[_destToken].harmonyMappedToken;

        updateOnLock(msg.sender, _harmonyReceiver, _ethereumToken, harmonyToken, _ethereumTokenAmount, harmonyTokenAmount);
    }

    function swapTokenForWrappedETH(
        address _harmonyReceiver,
        address _ethereumToken,
        uint256 _ethereumTokenAmount
    ) public
        availableNonce
        nonReentrant
        tokenMustBeActive(_ethereumToken)
        amountMustGreaterThanZero(_ethereumTokenAmount)
        receiverMustBeValid(_harmonyReceiver)
    {
        uint256 fee = calculateFee(_ethereumTokenAmount);

        require(
            IERC20(_ethereumToken).transferFrom(
                msg.sender,
                address(this),
                _ethereumTokenAmount + fee
            ),
            "Contract token allowances insufficient to complete this lock request"
        );

        DataTypes.ReserveData memory reserve = lendingPool.getReserveData(_ethereumToken);
        if(reserve.aTokenAddress != address(0)){
            lendingPool.deposit(
                _ethereumToken,
                _ethereumTokenAmount + fee,
                address(this),
                0
            );
        }

        BandOracleInterface.ReferenceData memory data = bandOracleInterface
            .getReferenceData(ERC20Detailed(_ethereumToken).symbol(), "ETH");
        uint256 amountWETH = _ethereumTokenAmount.mul(data.rate).div(1e18);

        address harmonyWETH = tokensData[ETHAddress].harmonyMappedToken;

        updateOnLock(msg.sender, _harmonyReceiver, _ethereumToken, harmonyWETH, _ethereumTokenAmount, amountWETH);
    }

    function swapTokenForONE(
        address _harmonyReceiver,
        address _ethereumToken,
        uint256 _ethereumTokenAmount
    ) public
        availableNonce
        nonReentrant
        tokenMustBeActive(_ethereumToken)
        amountMustGreaterThanZero(_ethereumTokenAmount)
        receiverMustBeValid(_harmonyReceiver)
    {
        uint256 fee = calculateFee(_ethereumTokenAmount);

        require(
            IERC20(_ethereumToken).transferFrom(
                msg.sender,
                address(this),
                _ethereumTokenAmount + fee
            ),
            "Contract token allowances insufficient to complete this lock request"
        );

        DataTypes.ReserveData memory reserve = lendingPool.getReserveData(_ethereumToken);
        if(reserve.aTokenAddress != address(0)){
            lendingPool.deposit(
                _ethereumToken,
                _ethereumTokenAmount + fee,
                address(this),
                0
            );
        }

        BandOracleInterface.ReferenceData memory data = bandOracleInterface
            .getReferenceData(ERC20Detailed(_ethereumToken).symbol(), "ONE");
        uint256 amountONE = _ethereumTokenAmount.mul(data.rate).div(1e18);

        updateOnLock(msg.sender, _harmonyReceiver, _ethereumToken, ONEAddress, _ethereumTokenAmount, amountONE);
    }

    function unlockERC20(
        address payable _ethereumReceiver,
        address _ethereumToken,
        uint256 _ethereumTokenAmount
    ) public
        nonReentrant
        onlyHarmonyBridge
        amountMustGreaterThanZero(_ethereumTokenAmount)
        receiverMustBeValid(_ethereumReceiver)
    {
        require(tokensData[_ethereumToken].harmonyMappedToken != address(0), "Invalid token address");

        require(_ethereumTokenAmount <= getTotalERC20Balance(_ethereumToken),
            "Exceeded amount of Token allowed to withdraw"
        );

        uint256 selfBalance = IERC20(_ethereumToken).balanceOf(address(this));
        if (_ethereumTokenAmount <= selfBalance) {
            IERC20(_ethereumToken).transfer(_ethereumReceiver, _ethereumTokenAmount);
        } else {
            lendingPool.withdraw(_ethereumToken, _ethereumTokenAmount - selfBalance, _ethereumReceiver);
            IERC20(_ethereumToken).transfer(_ethereumReceiver, selfBalance);
        }

        updateOnUnlock(_ethereumToken, _ethereumReceiver, _ethereumTokenAmount);
    }

    function unlockETH(address payable _ethereumReceiver, uint256 _amountETH)
        public
        onlyHarmonyBridge
        nonReentrant
        amountMustGreaterThanZero(_amountETH)
        receiverMustBeValid(_ethereumReceiver)
    {
        require(_amountETH <= getTotalETHBalance(),
            "Exceeded amount of ETH allowed to withdraw"
        );

        if (_amountETH <= address(this).balance) {
            _ethereumReceiver.transfer(_amountETH);
        } else {
            wethGateway.withdrawETH(_amountETH - address(this).balance, _ethereumReceiver);
            _ethereumReceiver.transfer(address(this).balance);
        }

        updateOnUnlock(ETHAddress, _ethereumReceiver, _amountETH);
    }

    function updateOracle(address _oracleAddress) public onlyOperator {
        oracle = Oracle(_oracleAddress);
        emit EthUpdateOracle(_oracleAddress);
    }

    function updateHmyBridge(address _harmonyBridge) public onlyOperator {
        harmonyBridge = HarmonyBridge(_harmonyBridge);
        emit EthUpdateHarmonyBridge(_harmonyBridge);
    }

    function updateFee(uint256 _feeNumerator, uint256 _feeDenominator)
        public
        onlyOperator
    {
        feeNumerator = _feeNumerator;
        feeDenominator = _feeDenominator;
        emit EthUpdateFee(_feeNumerator, _feeDenominator);
    }

    function withdrawETH(address payable _ethereumReceiver, uint256 _amountETH) public onlyOperator nonReentrant {
        require(_amountETH <= getTotalETHBalance() - tokensData[ETHAddress].lockedFund,
            "Exceeded amount of ETH allowed to withdraw"
        );

        if (_amountETH <= address(this).balance) {
            _ethereumReceiver.transfer(_amountETH);
        } else {
            wethGateway.withdrawETH(_amountETH - address(this).balance, _ethereumReceiver);
            _ethereumReceiver.transfer(address(this).balance);
        }
        emit EthWithdrawETH(_ethereumReceiver, _amountETH);
    }

    function withdrawERC20(address _ethereumToken, address _ethereumReceiver, uint256 _ethereumTokenAmount) public onlyOperator nonReentrant {
        require(_ethereumTokenAmount <= getTotalERC20Balance(_ethereumToken) - tokensData[_ethereumToken].lockedFund,
            "Exceeded amount of Token allowed to withdraw"
        );

        uint256 selfBalance = IERC20(_ethereumToken).balanceOf(address(this));
        if (_ethereumTokenAmount <= selfBalance) {
            IERC20(_ethereumToken).transfer(_ethereumReceiver, _ethereumTokenAmount);
        } else {
            lendingPool.withdraw(_ethereumToken, _ethereumTokenAmount - selfBalance, _ethereumReceiver);
            IERC20(_ethereumToken).transfer(_ethereumReceiver, selfBalance);
        }
        emit EthWithdrawERC20(_ethereumToken, _ethereumReceiver, _ethereumTokenAmount);
    }

    function getTotalETHBalance() public view returns (uint256) {
        DataTypes.ReserveData memory reserve = lendingPool.getReserveData(WETH);
        address aToken = reserve.aTokenAddress;
        if(aToken != address(0)){
            return address(this).balance + IERC20(aToken).balanceOf(address(this));
        } else {
            return address(this).balance;
        }
    }

    function getTotalERC20Balance(address _ethereumToken) public view returns (uint256) {
        DataTypes.ReserveData memory reserve = lendingPool.getReserveData(_ethereumToken);
        address aToken = reserve.aTokenAddress;
        if(aToken != address(0)){
            return IERC20(_ethereumToken).balanceOf(address(this)) + IERC20(aToken).balanceOf(address(this));
        } else {
            return IERC20(_ethereumToken).balanceOf(address(this));
        }
    }

    function calculateFee(uint256 _amountToken) internal view returns (uint256) {
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
        address _ethereumSender,
        address _harmonyReceiver,
        address _ethereumToken,
        address _harmonyToken,
        uint256 _ethereumTokenAmount,
        uint256 _harmonyTokenAmount
    ) internal {
        require(_harmonyTokenAmount > 0, "Amount token must be greater than zero");

        lockNonce = lockNonce.add(1);

        tokensData[_ethereumToken].lockedFund = tokensData[_ethereumToken].lockedFund.add(_ethereumTokenAmount);

        emit EthLogLock(
            _ethereumSender,
            _harmonyReceiver,
            _ethereumToken,
            _harmonyToken,
            _ethereumTokenAmount,
            _harmonyTokenAmount,
            lockNonce
        );
    }

    function updateOnUnlock(address _ethereumToken, address _ethereumReceiver, uint256 _ethereumTokenAmount) internal {
        if(tokensData[_ethereumToken].lockedFund >= _ethereumTokenAmount){
            tokensData[_ethereumToken].lockedFund = tokensData[_ethereumToken].lockedFund.sub(_ethereumTokenAmount);
        } else {
            tokensData[_ethereumToken].lockedFund = 0;
        }

        emit EthLogUnlock(_ethereumReceiver, _ethereumToken, _ethereumTokenAmount);
    }

    function checkUnlockable(address _token, uint256 _amount) public view returns (bool) {
        if (_token == ETHAddress) {
            return getTotalETHBalance() >= _amount;
        } else {
            return getTotalERC20Balance(_token) >= _amount;
        }
    }
}
