package contract

// -------------------------------------------------------
//    Contract Contains functionality for loading the
//				 smart contract
// -------------------------------------------------------

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/txs"
)

// File paths to Ethereum smart contract ABIs
const (
	EthereumBridgeBankABI     = "/generated/ethereum/abi/BridgeBank/BridgeBank.abi"
	EthereumHarmonyBridgeABI  = "/generated/ethereum/abi/HarmonyBridge/HarmonyBridge.abi"
	EthereumBridgeRegistryABI = "/generated/ethereum/abi/BridgeRegistry/BridgeRegistry.abi"

	HarmonyBridgeBankABI     = "/generated/harmony/abi/BridgeBank/BridgeBank.abi"
	HarmonyEthereumBridgeABI = "/generated/harmony/abi/EthereumBridge/EthereumBridge.abi"
	HarmonyBridgeRegistryABI = "/generated/harmony/abi/BridgeRegistry/BridgeRegistry.abi"
)

// EthoadABI loads a smart contract as an abi.ABI
func EthLoadABI(contractType txs.ContractRegistry) abi.ABI {
	var (
		_, b, _, _ = runtime.Caller(0)
		dir        = filepath.Dir(b)
	)

	var filePath string
	switch contractType {
	case txs.HarmonyBridge:
		filePath = EthereumHarmonyBridgeABI
	case txs.BridgeBank:
		filePath = EthereumBridgeBankABI
	case txs.BridgeRegistry:
		filePath = EthereumBridgeRegistryABI
	}

	// Read the file containing the contract's ABI
	contractRaw, err := ioutil.ReadFile(dir + filePath)
	if err != nil {
		panic(err)
	}

	// Convert the raw abi into a usable format
	contractABI, err := abi.JSON(strings.NewReader(string(contractRaw)))
	if err != nil {
		panic(err)
	}
	return contractABI
}

// LoadHarmonyABI loads a smart contract as an abi.ABI
func HmyLoadABI(contractType txs.ContractRegistry) abi.ABI {
	var (
		_, b, _, _ = runtime.Caller(0)
		dir        = filepath.Dir(b)
	)

	var filePath string
	switch contractType {
	case txs.EthereumBridge:
		filePath = HarmonyEthereumBridgeABI
	case txs.BridgeBank:
		filePath = HarmonyBridgeBankABI
	case txs.BridgeRegistry:
		filePath = HarmonyBridgeRegistryABI
	}

	// Read the file containing the contract's ABI
	contractRaw, err := ioutil.ReadFile(dir + filePath)
	if err != nil {
		panic(err)
	}

	// Convert the raw abi into a usable format
	contractABI, err := abi.JSON(strings.NewReader(string(contractRaw)))
	if err != nil {
		panic(err)
	}
	return contractABI
}
