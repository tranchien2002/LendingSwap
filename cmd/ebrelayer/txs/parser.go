package txs

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/types"
)

const (
	nullAddress = "0x0000000000000000000000000000000000000000"
)

// EthUnlockClaimToSignedOracleClaim packages and signs a unlock claim's data, returning a new oracle claim
func EthUnlockClaimToSignedOracleClaim(event types.EthLogNewUnlockClaimEvent, key *ecdsa.PrivateKey) (EthOracleClaim, error) {
	oracleClaim := EthOracleClaim{}

	// Generate a hashed claim message which contains UnlockClaim's data
	fmt.Println("Generating unique message for UnlockClaim", event.UnlockID)
	message := EthGenerateClaimMessage(event)

	// Sign the message using the validator's private key
	fmt.Println("Signing message...")
	signature, err := SignClaim(PrefixMsg(message), key)
	if err != nil {
		return oracleClaim, err
	}
	fmt.Println("Signature generated:", hexutil.Encode(signature))

	oracleClaim.UnlockID = event.UnlockID
	var message32 [32]byte
	copy(message32[:], message)
	oracleClaim.Message = message32
	oracleClaim.Signature = signature
	return oracleClaim, nil
}

// HmyUnlockClaimToSignedOracleClaim packages and signs a unlock claim's data, returning a new oracle claim
func HmyUnlockClaimToSignedOracleClaim(event types.HmyLogNewUnlockClaimEvent, key *ecdsa.PrivateKey) (HmyOracleClaim, error) {
	oracleClaim := HmyOracleClaim{}

	// Generate a hashed claim message which contains UnlockClaim's data
	fmt.Println("Generating unique message for UnlockClaim", event.UnlockID)
	message := HmyGenerateClaimMessage(event)

	// Sign the message using the validator's private key
	fmt.Println("Signing message...")
	signature, err := SignClaim(PrefixMsg(message), key)
	if err != nil {
		return oracleClaim, err
	}
	fmt.Println("Signature generated:", hexutil.Encode(signature))

	oracleClaim.UnlockID = event.UnlockID
	var message32 [32]byte
	copy(message32[:], message)
	oracleClaim.Message = message32
	oracleClaim.Signature = signature
	return oracleClaim, nil
}

// isZeroAddress checks an Ethereum address and returns a bool which indicates if it is the null address
func isZeroAddress(address common.Address) bool {
	return address == common.HexToAddress(nullAddress)
}

// HarmonyEventToEthereumClaim parses and packages an Ethereum event struct with a validator address in an EthBridgeClaim msg
func HarmonyEventToEthereumClaim(event *types.HmyLogLockEvent) (EthUnlockClaim, error) {
	witnessClaim := EthUnlockClaim{}

	// chainID type casting (*big.Int -> int)
	chainID := event.HarmonyChainID

	// harmonySender type casting (address.common -> string)
	harmonySender := event.HarmonySender

	// ethereumReceiver type casting (address.common -> string)
	ethereumReceiver := event.EthereumReceiver

	// token type casting (address.common -> string)
	token := event.EthereumToken

	// amount is
	amount := event.EthereumTokenAmount

	// Package the information in a unique EthBridgeClaim
	witnessClaim.HarmonyChainID = chainID
	witnessClaim.Token = token
	witnessClaim.HarmonySender = harmonySender
	witnessClaim.EthereumReceiver = ethereumReceiver
	witnessClaim.Amount = amount

	return witnessClaim, nil
}

// EthereumEventToHarmonyClaim parses and packages an Ethereum event struct with a validator address in an EthBridgeClaim msg
func EthereumEventToHarmonyClaim(event *types.EthLogLockEvent) (HmyUnlockClaim, error) {
	witnessClaim := HmyUnlockClaim{}

	// chainID type casting (*big.Int -> int)
	chainID := event.EthereumChainID

	// ethereumSender type casting (address.common -> string)
	ethereumSender := event.EthereumSender

	// harmonyReceiver type casting (address.common -> string)
	harmonyReceiver := event.HarmonyReceiver

	// token type casting (address.common -> string)
	token := event.HarmonyToken

	// amount is
	amount := event.HarmonyTokenAmount

	// Package the information in a unique EthBridgeClaim
	witnessClaim.EthereumChainID = chainID
	witnessClaim.Token = token
	witnessClaim.EthereumSender = ethereumSender
	witnessClaim.HarmonyReceiver = harmonyReceiver
	witnessClaim.Amount = amount

	return witnessClaim, nil
}
