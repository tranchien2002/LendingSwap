package relayer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io"
	"math/big"
	"os"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	htypes "github.com/harmony-one/harmony/core/types"
	tmLog "github.com/tendermint/tendermint/libs/log"
	"github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/contract"
	"github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/txs"
	"github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/types"
	"github.com/trinhtan/horizon-hackathon/hmyclient"
)

// HarmonySub is an Harmony listener that can relay txs to Ethereum
type HarmonySub struct {
	HarmonyProvider        string
	EthereumProvider       string
	HarmonyBridgeRegistry  common.Address
	EthereumBridgeRegistry common.Address
	ValidatorName          string
	HmyPrivateKey          *ecdsa.PrivateKey
	EthPrivateKey          *ecdsa.PrivateKey
	Logger                 tmLog.Logger
}

// NewHarmonySub initializes a new HarmonySub
func NewHarmonySub(inBuf io.Reader, validatorMoniker,
	harmonyProvider string, ethereumProvider string, harmonyBridgeRegistry common.Address, ethereumBridgeRegistry common.Address, hmyPrivateKey *ecdsa.PrivateKey, ethPrivateKey *ecdsa.PrivateKey,
	logger tmLog.Logger) (HarmonySub, error) {

	return HarmonySub{
		HarmonyProvider:        harmonyProvider,
		EthereumProvider:       ethereumProvider,
		HarmonyBridgeRegistry:  harmonyBridgeRegistry,
		EthereumBridgeRegistry: ethereumBridgeRegistry,
		ValidatorName:          "validator",
		HmyPrivateKey:          hmyPrivateKey,
		EthPrivateKey:          ethPrivateKey,
		Logger:                 logger,
	}, nil
}

// Start an Harmony chain subscription
func (sub HarmonySub) Start() {
	client, err := HmySetupWebsocketClient(sub.HarmonyProvider)
	if err != nil {
		sub.Logger.Error(err.Error())
		os.Exit(1)
	}
	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		sub.Logger.Error(err.Error())
		os.Exit(1)
	}
	sub.Logger.Info("Started Harmony websocket with provider:", sub.HarmonyProvider)

	// We will check logs for new events
	logs := make(chan htypes.Log)

	// Start BridgeBank subscription, prepare contract ABI and HmyLockLog event signature
	bridgeBankAddress, subBridgeBank := sub.HmyStartContractEventSub(logs, client, txs.BridgeBank)
	bridgeBankContractABI := contract.HmyLoadABI(txs.BridgeBank)
	eventLogLockSignature := bridgeBankContractABI.Events[types.HmyLogLock.String()].ID.Hex()

	// Start ethereumBridge subscription, prepare contract ABI and HmyLogNewUnlockClaim event signature
	_, subEthereumBridge := sub.HmyStartContractEventSub(logs, client, txs.EthereumBridge)
	ethereumBridgeContractABI := contract.HmyLoadABI(txs.EthereumBridge)
	eventLogNewUnlockClaimSignature := ethereumBridgeContractABI.Events[types.HmyLogNewUnlockClaim.String()].ID.Hex()

	for {
		select {
		// Handle any errors
		case err := <-subBridgeBank.Err():
			sub.Logger.Error("Harmony - Sub bridgeBank error: ", err.Error())
			client, err = HmySetupWebsocketClient(sub.HarmonyProvider)
			if err != nil {
				sub.Logger.Error(err.Error())
				os.Exit(1)
			}
			_, subBridgeBank = sub.HmyStartContractEventSub(logs, client, txs.BridgeBank)
			_, subEthereumBridge = sub.HmyStartContractEventSub(logs, client, txs.EthereumBridge)
		case err := <-subEthereumBridge.Err():
			sub.Logger.Error("Harmony - Sub ethereumBridge error: ", err.Error())
			client, err = HmySetupWebsocketClient(sub.HarmonyProvider)
			if err != nil {
				sub.Logger.Error(err.Error())
				os.Exit(1)
			}
			_, subBridgeBank = sub.HmyStartContractEventSub(logs, client, txs.BridgeBank)
			_, subEthereumBridge = sub.HmyStartContractEventSub(logs, client, txs.EthereumBridge)
		// vLog is raw event data
		case vLog := <-logs:
			sub.Logger.Info(fmt.Sprintf("Witnessed tx %s on block %d\n", vLog.TxHash.Hex(), vLog.BlockNumber))
			var err error
			switch vLog.Topics[0].Hex() {
			case eventLogLockSignature:
				err = sub.HmyHandleLogLockEvent(clientChainID, bridgeBankAddress, bridgeBankContractABI, types.HmyLogLock.String(), vLog)
			case eventLogNewUnlockClaimSignature:
				err = sub.HmyHandleLogNewUnlockClaim(sub.HarmonyBridgeRegistry, ethereumBridgeContractABI,
					types.HmyLogNewUnlockClaim.String(), vLog)
			}
			// TODO: Check local events store for status, if retryable, attempt relay again
			if err != nil {
				sub.Logger.Error("Harmony error: ", err.Error())
			}
		}
	}

}

// HmyStartContractEventSub : starts an event subscription on the specified Harmony contract
func (sub HarmonySub) HmyStartContractEventSub(logs chan htypes.Log, client *hmyclient.Client,
	contractName txs.ContractRegistry) (common.Address, ethereum.Subscription) {
	// Get the contract address for this subscription
	subContractAddress, err := txs.HmyGetAddressFromBridgeRegistry(sub.HmyPrivateKey, client, sub.HarmonyBridgeRegistry, contractName)
	if err != nil {
		sub.Logger.Error(err.Error())
	}

	// We need the address in []bytes for the query
	subQuery := ethereum.FilterQuery{
		Addresses: []common.Address{subContractAddress},
	}

	// Start the contract subscription
	contractSub, err := client.SubscribeFilterLogs(context.Background(), subQuery, logs)
	if err != nil {
		sub.Logger.Error(err.Error())
	}
	sub.Logger.Info(fmt.Sprintf("Harmony - Subscribed to %v contract at address: %s", contractName, subContractAddress.Hex()))
	return subContractAddress, contractSub
}

// HmyHandleLogLockEvent unpacks a HmyLogLockEvent, and relays a tx to Ethereum
func (sub HarmonySub) HmyHandleLogLockEvent(clientChainID *big.Int, bridgeBankAddress common.Address,
	contractABI abi.ABI, eventName string, cLog htypes.Log) error {
	// Parse the event's attributes via contract ABI
	event := types.HmyLogLockEvent{}
	err := contractABI.Unpack(&event, eventName, cLog.Data)
	if err != nil {
		sub.Logger.Error("error unpacking: %v", err)
	}
	event.BridgeBankAddress = bridgeBankAddress
	event.HarmonyChainID = clientChainID

	sub.Logger.Info(event.String())

	types.HmyNewEventWrite(cLog.TxHash.Hex(), event)

	unlockClaim, err := txs.HarmonyEventToEthereumClaim(&event)
	if err != nil {
		return err
	}

	return txs.RelayUnlockClaimToEthereum(sub.EthereumProvider, sub.EthereumBridgeRegistry, types.HmyLogLock, unlockClaim, sub.EthPrivateKey)
}

// HmyHandleLogNewUnlockClaim unpacks a HmyLogNewUnlockClaim event, builds a new OracleClaim, and relays it to Harmony
func (sub HarmonySub) HmyHandleLogNewUnlockClaim(contractAddress common.Address, contractABI abi.ABI,
	eventName string, hLog htypes.Log) error {
	// Parse the event's attributes via contract ABI
	event := types.HmyLogNewUnlockClaimEvent{}
	err := contractABI.Unpack(&event, eventName, hLog.Data)
	if err != nil {
		sub.Logger.Error("error unpacking: %v", err)
	}
	sub.Logger.Info(event.String())

	oracleClaim, err := txs.HmyUnlockClaimToSignedOracleClaim(event, sub.HmyPrivateKey)
	if err != nil {
		return err
	}
	return txs.RelayOracleClaimToHarmony(sub.HarmonyProvider, contractAddress, types.HmyLogNewUnlockClaim,
		oracleClaim, sub.HmyPrivateKey)
}
