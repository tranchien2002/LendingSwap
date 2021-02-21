package types

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Event enum containing supported chain events
type Event byte

const (
	// Unsupported is an invalid Harmony or Ethereum event
	Unsupported Event = iota
	// EthLogLock is for Ethereum event EthLogLock
	EthLogLock
	// EthLogNewUnlockClaim is an Ethereum event named 'EthLogNewUnlockClaim'
	EthLogNewUnlockClaim
	// HmyLogLock is for Harmony event HmyLogLock
	HmyLogLock
	// HmyLogNewUnlockClaim is for Harmony event HmyLogNewUnlock
	HmyLogNewUnlockClaim
)

// String returns the event type as a string
func (d Event) String() string {
	return [...]string{"unsupported", "EthLogLock", "EthLogNewUnlockClaim", "HmyLogLock", "HmyLogNewUnlockClaim"}[d]
}

// EthLogLockEvent struct is used by EthLogLock
type EthLogLockEvent struct {
	EthereumChainID     *big.Int
	BridgeBankAddress   common.Address
	ID                  [32]byte
	EthereumSender      common.Address
	HarmonyReceiver     common.Address
	EthereumToken       common.Address
	HarmonyToken        common.Address
	EthereumTokenAmount *big.Int
	HarmonyTokenAmount  *big.Int
	Nonce               *big.Int
}

// String implements fmt.Stringer
func (e EthLogLockEvent) String() string {
	return fmt.Sprintf("\nChain ID: %v\nBridge contract address: %v\nEthereum Token: %v\nHarmony Token: %v\nEthereum Sender: %v\nHarmony Recipient: %v\nEthereum Token Amount: %v\nHarmony Token Amount: %v\nNonce: %v\n",
		e.EthereumChainID, e.BridgeBankAddress.Hex(), e.EthereumToken.Hex(), e.HarmonyToken.Hex(), e.EthereumSender.Hex(),
		e.HarmonyReceiver.Hex(), e.EthereumTokenAmount, e.HarmonyTokenAmount, e.Nonce)
}

// EthLogNewUnlockClaimEvent struct which represents a EthLogNewUnlockClaim event
type EthLogNewUnlockClaimEvent struct {
	UnlockID         *big.Int
	HarmonySender    common.Address
	EthereumReceiver common.Address
	ValidatorAddress common.Address
	TokenAddress     common.Address
	Amount           *big.Int
}

// String implements fmt.Stringer
func (p EthLogNewUnlockClaimEvent) String() string {
	return fmt.Sprintf("\nUnlocl ID: %v\nHarmony Sender: %v\n"+
		"Ethereum Receiver: %v\nEthereum Token: %v\nAmount: %v\nValidator Address: %v\n\n",
		p.UnlockID, p.HarmonySender.Hex(), p.EthereumReceiver.Hex(),
		p.TokenAddress.Hex(), p.Amount, p.ValidatorAddress.Hex())
}

// HmyLogLockEvent struct is used by HmyLogLock
type HmyLogLockEvent struct {
	HarmonyChainID      *big.Int
	BridgeBankAddress   common.Address
	ID                  [32]byte
	HarmonySender       common.Address
	EthereumReceiver    common.Address
	HarmonyToken        common.Address
	EthereumToken       common.Address
	HarmonyTokenAmount  *big.Int
	EthereumTokenAmount *big.Int
	Nonce               *big.Int
}

// String implements fmt.Stringer
func (e HmyLogLockEvent) String() string {
	return fmt.Sprintf("\nChain ID: %v\nBridge contract address: %v\nHarmony Token: %v\nEthereum Token: %v\nHarmony Sender: %v\nEthereum Recipient: %v\nHarmony Token Amount: %v\nEthereum Token Amount: %v\nNonce: %v\n",
		e.HarmonyChainID, e.BridgeBankAddress.Hex(), e.HarmonyToken.Hex(), e.EthereumToken.Hex(), e.HarmonySender.Hex(),
		e.EthereumReceiver.Hex(), e.HarmonyTokenAmount, e.EthereumTokenAmount, e.Nonce)
}

// HmyLogNewUnlockClaimEvent struct which represents a HmyLogNewUnlockClaim event
type HmyLogNewUnlockClaimEvent struct {
	UnlockID         *big.Int
	EthereumSender   common.Address
	HarmonyReceiver  common.Address
	ValidatorAddress common.Address
	TokenAddress     common.Address
	Amount           *big.Int
}

// String implements fmt.Stringer
func (p HmyLogNewUnlockClaimEvent) String() string {
	return fmt.Sprintf("\nUnlocl ID: %v\nEthereum Sender: %v\n"+
		"Harmony Receiver: %v\nHarmony Token: %v\nAmount: %v\nValidator Address: %v\n\n",
		p.UnlockID, p.EthereumSender.Hex(), p.HarmonyReceiver.Hex(),
		p.TokenAddress.Hex(), p.Amount, p.ValidatorAddress.Hex())
}
