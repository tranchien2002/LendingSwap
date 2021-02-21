package txs

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	harmonybridge "github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/contract/generated/ethereum/bindings/harmonybridge"
	oracle "github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/contract/generated/ethereum/bindings/oracle"
	"github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/types"
)

const (
	// GasLimit the gas limit in Gwei used for transactions sent with TransactOpts
	GasLimit = uint64(3000000)
)

// RelayUnlockClaimToEthereum relays the provided UnlockClaim to HarmonyBridge contract on the Ethereum network
func RelayUnlockClaimToEthereum(ethereumProvider string, ethereumBridgeRegistry common.Address, event types.Event,
	claim EthUnlockClaim, privateKey *ecdsa.PrivateKey) error {
	// Initialize client service, validator's tx auth, and target contract address
	client, auth, target := EthInitRelayConfig(ethereumProvider, ethereumBridgeRegistry, event, privateKey)

	// Initialize HarmonyBridge instance
	fmt.Println("\nFetching HarmonyBridge contract...")
	harmonyBridgeInstance, err := harmonybridge.NewHarmonyBridge(target, client)
	if err != nil {
		log.Fatal(err)
	}

	// Send transaction
	fmt.Println("Sending new UnlockClaim to HarmonyBridge...")
	tx, err := harmonyBridgeInstance.NewUnlockClaim(auth,
		claim.HarmonySender, claim.EthereumReceiver, claim.Token, claim.Amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("NewUnlockClaim tx hash:", tx.Hash().Hex())

	return nil
}

// RelayOracleClaimToEthereum relays the provided OracleClaim to Oracle contract on the Ethereum network
func RelayOracleClaimToEthereum(provider string, contractAddress common.Address, event types.Event,
	claim EthOracleClaim, privateKey *ecdsa.PrivateKey) error {
	// Initialize client service, validator's tx auth, and target contract address
	client, auth, target := EthInitRelayConfig(provider, contractAddress, event, privateKey)

	// Initialize Oracle instance
	fmt.Println("\nFetching Oracle contract...")
	oracleInstance, err := oracle.NewOracle(target, client)
	if err != nil {
		log.Fatal(err)
	}

	// Send transaction
	fmt.Println("Sending new OracleClaim to Oracle...")
	tx, err := oracleInstance.NewOracleClaim(auth, claim.UnlockID, claim.Message, claim.Signature)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("NewOracleClaim tx hash:", tx.Hash().Hex())
	return nil
}

// EthInitRelayConfig set up Ethereum client, validator's transaction auth, and the target contract's address
func EthInitRelayConfig(provider string, registry common.Address, event types.Event, privateKey *ecdsa.PrivateKey,
) (*ethclient.Client, *bind.TransactOpts, common.Address) {
	// Start Ethereum client
	client, err := ethclient.Dial(provider)
	if err != nil {
		log.Fatal(err)
	}

	// Load the validator's address
	sender, err := LoadSender(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), sender)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Set up TransactOpts auth's tx signature authorization
	transactOptsAuth := bind.NewKeyedTransactor(privateKey)
	transactOptsAuth.Nonce = big.NewInt(int64(nonce))
	transactOptsAuth.Value = big.NewInt(0) // in wei
	transactOptsAuth.GasLimit = GasLimit
	transactOptsAuth.GasPrice = gasPrice

	var targetContract ContractRegistry
	switch event {
	// New Claim
	case types.HmyLogLock:
		targetContract = HarmonyBridge
	// OracleClaims are sent to the Oracle contract
	case types.EthLogNewUnlockClaim:
		targetContract = Oracle
	default:
		panic("invalid target contract address")
	}

	// Get the specific contract's address
	target, err := EthGetAddressFromBridgeRegistry(privateKey, client, registry, targetContract)
	if err != nil {
		log.Fatal(err)
	}
	return client, transactOptsAuth, target
}
