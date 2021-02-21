package txs

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// EthOracleClaim contains data required to make an EthOracleClaim
type EthOracleClaim struct {
	UnlockID  *big.Int
	Message   [32]byte
	Signature []byte
}

// HmyOracleClaim contains data required to make an HmyOracleClaim
type HmyOracleClaim struct {
	UnlockID  *big.Int
	Message   [32]byte
	Signature []byte
}

// EthUnlockClaim contains data required to make an Ethereum UnlockClaim
type EthUnlockClaim struct {
	HarmonyChainID   *big.Int
	HarmonySender    common.Address
	EthereumReceiver common.Address
	Token            common.Address
	Amount           *big.Int
}

// HmyUnlockClaim contains data required to make an Harmony UnlockClaim
type HmyUnlockClaim struct {
	EthereumChainID *big.Int
	EthereumSender  common.Address
	HarmonyReceiver common.Address
	Token           common.Address
	Amount          *big.Int
}
