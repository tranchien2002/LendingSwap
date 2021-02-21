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

	ctypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	tmLog "github.com/tendermint/tendermint/libs/log"
	"github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/contract"
	"github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/txs"
	"github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/types"
)

// TODO: Move relay functionality out of EthereumSub into a new Relayer parent struct

// EthereumSub is an Ethereum listener that can relay txs to Harmony
type EthereumSub struct {
	EthereumProvider       string
	HarmonyProvider        string
	EthereumBridgeRegistry common.Address
	HarmonyBridgeRegistry  common.Address
	ValidatorName          string
	EthPrivateKey          *ecdsa.PrivateKey
	HmyPrivatekey          *ecdsa.PrivateKey
	Logger                 tmLog.Logger
}

// NewEthereumSub initializes a new EthereumSub
func NewEthereumSub(inBuf io.Reader, validatorMoniker, ethereumProvider string, harmonyProvider string,
	ethereumBridgeRegistry common.Address, harmonyBridgeRegistry common.Address, ethPrivateKey *ecdsa.PrivateKey, hmyPrivateKey *ecdsa.PrivateKey, logger tmLog.Logger) (EthereumSub, error) {
	return EthereumSub{
		EthereumProvider:       ethereumProvider,
		HarmonyProvider:        harmonyProvider,
		EthereumBridgeRegistry: ethereumBridgeRegistry,
		HarmonyBridgeRegistry:  harmonyBridgeRegistry,
		ValidatorName:          "validator",
		EthPrivateKey:          ethPrivateKey,
		HmyPrivatekey:          hmyPrivateKey,
		Logger:                 logger,
	}, nil
}

// Start an Ethereum chain subscription
func (sub EthereumSub) Start() {
	client, err := EthSetupWebsocketClient(sub.EthereumProvider)
	if err != nil {
		sub.Logger.Error(err.Error())
		os.Exit(1)
	}
	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		sub.Logger.Error(err.Error())
		os.Exit(1)
	}
	sub.Logger.Info("Started Ethereum websocket with provider:", sub.EthereumProvider)

	// We will check logs for new events
	logs := make(chan ctypes.Log)

	// Start BridgeBank subscription, prepare contract ABI and EthLogLock event signature
	bridgeBankAddress, subBridgeBank := sub.EthStartContractEventSub(logs, client, txs.BridgeBank)
	bridgeBankContractABI := contract.EthLoadABI(txs.BridgeBank)
	eventLogLockSignature := bridgeBankContractABI.Events[types.EthLogLock.String()].ID.Hex()

	// Start harmonyBridge subscription, prepare contract ABI and EthLogNewUnlockClaim event signature
	_, subHarmonyBridge := sub.EthStartContractEventSub(logs, client, txs.HarmonyBridge)
	harmonyBridgeContractABI := contract.EthLoadABI(txs.HarmonyBridge)
	eventLogNewUnlockClaimSignature := harmonyBridgeContractABI.Events[types.EthLogNewUnlockClaim.String()].ID.Hex()

	for {
		select {
		// Handle any errors
		case err := <-subBridgeBank.Err():
			sub.Logger.Error("Ethereum - Sub bridgeBank error: ", err.Error())
			client, err = EthSetupWebsocketClient(sub.EthereumProvider)
			if err != nil {
				sub.Logger.Error(err.Error())
				os.Exit(1)
			}
			_, subBridgeBank = sub.EthStartContractEventSub(logs, client, txs.BridgeBank)
			_, subHarmonyBridge = sub.EthStartContractEventSub(logs, client, txs.HarmonyBridge)
		case err := <-subHarmonyBridge.Err():
			sub.Logger.Error("Ethereum - Sub harmonyBridge error:", err.Error())
			client, err = EthSetupWebsocketClient(sub.EthereumProvider)
			if err != nil {
				sub.Logger.Error(err.Error())
				os.Exit(1)
			}
			_, subBridgeBank = sub.EthStartContractEventSub(logs, client, txs.BridgeBank)
			_, subHarmonyBridge = sub.EthStartContractEventSub(logs, client, txs.HarmonyBridge)
		// vLog is raw event data
		case vLog := <-logs:
			sub.Logger.Info(fmt.Sprintf("Witnessed tx %s on block %d\n", vLog.TxHash.Hex(), vLog.BlockNumber))
			var err error
			switch vLog.Topics[0].Hex() {
			case eventLogLockSignature:
				err = sub.EthHandleLogLockEvent(clientChainID, bridgeBankAddress, bridgeBankContractABI,
					types.EthLogLock.String(), vLog)
			case eventLogNewUnlockClaimSignature:
				err = sub.EthHandleLogNewUnlockClaim(sub.EthereumBridgeRegistry, harmonyBridgeContractABI,
					types.EthLogNewUnlockClaim.String(), vLog)
			}
			// TODO: Check local events store for status, if retryable, attempt relay again
			if err != nil {
				sub.Logger.Error("Ethereum error: ", err.Error())
			}
		}
	}
}

// EthStartContractEventSub : starts an event subscription on the specified Ethereum contract
func (sub EthereumSub) EthStartContractEventSub(logs chan ctypes.Log, client *ethclient.Client,
	contractName txs.ContractRegistry) (common.Address, ethereum.Subscription) {
	// Get the contract address for this subscription
	subContractAddress, err := txs.EthGetAddressFromBridgeRegistry(sub.EthPrivateKey, client, sub.EthereumBridgeRegistry, contractName)
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
	sub.Logger.Info(fmt.Sprintf("Ethereum - Subscribed to %v contract at address: %s", contractName, subContractAddress.Hex()))
	return subContractAddress, contractSub
}

// EthHandleLogLockEvent unpacks an EthLogLockEvent, and relays a tx to Harmony
func (sub EthereumSub) EthHandleLogLockEvent(clientChainID *big.Int, contractAddress common.Address,
	contractABI abi.ABI, eventName string, cLog ctypes.Log) error {
	// Parse the event's attributes via contract ABI
	fmt.Println(cLog, "\n")
	event := types.EthLogLockEvent{}
	err := contractABI.Unpack(&event, eventName, cLog.Data)
	if err != nil {
		sub.Logger.Error("error unpacking: %v", err)
	}
	event.BridgeBankAddress = contractAddress
	event.EthereumChainID = clientChainID
	sub.Logger.Info(event.String())

	// Add the event to the record
	types.EthNewEventWrite(cLog.TxHash.Hex(), event)

	unlockClaim, err := txs.EthereumEventToHarmonyClaim(&event)
	if err != nil {
		return err
	}

	return txs.RelayUnlockClaimToHarmony(sub.HarmonyProvider, sub.HarmonyBridgeRegistry, types.EthLogLock, unlockClaim, sub.HmyPrivatekey)
}

// EthHandleLogNewUnlockClaim unpacks a EthLogNewUnlockClaim event, builds a new OracleClaim, and relays it to Ethereum
func (sub EthereumSub) EthHandleLogNewUnlockClaim(contractAddress common.Address, contractABI abi.ABI,
	eventName string, cLog ctypes.Log) error {
	// Parse the event's attributes via contract ABI
	event := types.EthLogNewUnlockClaimEvent{}
	err := contractABI.Unpack(&event, eventName, cLog.Data)
	if err != nil {
		sub.Logger.Error("error unpacking: %v", err)
	}
	sub.Logger.Info(event.String())

	oracleClaim, err := txs.EthUnlockClaimToSignedOracleClaim(event, sub.EthPrivateKey)
	if err != nil {
		return err
	}
	return txs.RelayOracleClaimToEthereum(sub.EthereumProvider, contractAddress, types.EthLogNewUnlockClaim,
		oracleClaim, sub.EthPrivateKey)
}
