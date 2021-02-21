package txs

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/harmony-one/harmony/accounts/abi/bind"
	ethereumbridge "github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/contract/generated/harmony/bindings/ethereumbridge"
	oracle "github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/contract/generated/harmony/bindings/oracle"
	"github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/types"
	"github.com/trinhtan/horizon-hackathon/hmyclient"
)

// RelayUnlockClaimToHarmony relays the provided UnlockClaim to EthereumBridge contract on the Ethereum network
func RelayUnlockClaimToHarmony(harmonyProvider string, ethereumBridgeRegistry common.Address, event types.Event,
	claim HmyUnlockClaim, privateKey *ecdsa.PrivateKey) error {
	// Initialize client service, validator's tx auth, and target contract address
	client, auth, target := HmyInitRelayConfig(harmonyProvider, ethereumBridgeRegistry, event, privateKey)

	// Initialize EthereumBridge instance
	fmt.Println("\nFetching EthereumBridge contract...")
	ethereumBridgeInstance, err := ethereumbridge.NewEthereumBridge(target, client)
	if err != nil {
		log.Fatal(err)
	}

	// Send transaction
	fmt.Println("Sending new UnlockClaim to EthereumBridge...")
	tx, err := ethereumBridgeInstance.NewUnlockClaim(auth,
		claim.EthereumSender, claim.HarmonyReceiver, claim.Token, claim.Amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("NewUnlockClaim tx hash:", tx.Hash().Hex())
	return nil
}

// RelayOracleClaimToHarmony relays the provided OracleClaim to Oracle contract on the Ethereum network
func RelayOracleClaimToHarmony(provider string, contractAddress common.Address, event types.Event,
	claim HmyOracleClaim, privateKey *ecdsa.PrivateKey) error {
	// Initialize client service, validator's tx auth, and target contract address
	client, auth, target := HmyInitRelayConfig(provider, contractAddress, event, privateKey)

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

// HmyInitRelayConfig set up Ethereum client, validator's transaction auth, and the target contract's address
func HmyInitRelayConfig(provider string, registry common.Address, event types.Event, privateKey *ecdsa.PrivateKey,
) (*hmyclient.Client, *bind.TransactOpts, common.Address) {
	// Start Ethereum client
	client, err := hmyclient.Dial(provider)
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
	transactOptsAuth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(2))
	if err != nil {
		log.Fatal(err)
	}
	transactOptsAuth.Nonce = big.NewInt(int64(nonce))
	transactOptsAuth.Value = big.NewInt(0) // in wei
	transactOptsAuth.GasLimit = GasLimit
	transactOptsAuth.GasPrice = gasPrice

	var targetContract ContractRegistry
	switch event {
	// New Claim
	case types.EthLogLock:
		targetContract = EthereumBridge
	// OracleClaims are sent to the Oracle contract
	case types.HmyLogNewUnlockClaim:
		targetContract = Oracle
	default:
		panic("invalid target contract address")
	}

	// Get the specific contract's address
	target, err := HmyGetAddressFromBridgeRegistry(privateKey, client, registry, targetContract)
	if err != nil {
		log.Fatal(err)
	}
	return client, transactOptsAuth, target
}
