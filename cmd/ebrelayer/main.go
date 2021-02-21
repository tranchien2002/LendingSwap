package main

import (
	"bufio"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	tmLog "github.com/tendermint/tendermint/libs/log"

	"github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/contract"
	"github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/relayer"
	"github.com/trinhtan/horizon-hackathon/cmd/ebrelayer/txs"
)

func init() {
	// Construct Root Command
	rootCmd.AddCommand(
		initRelayerCmd(),
		generateBindingsCmd(),
	)
}

var rootCmd = &cobra.Command{
	Use:          "ebrelayer",
	Short:        "Streams live events from Ethereum and Cosmos and relays event information to the opposite chain",
	SilenceUsage: true,
}

//	initRelayerCmd
func initRelayerCmd() *cobra.Command {

	initRelayerCmd := &cobra.Command{
		Use:     "init [ethereumProvider] [Eth-bridgeRegistryContractAddress] [harmonyProvider] [Hmy-bridgeRegistryContract] [validatorMoniker]",
		Short:   "Validate credentials and initialize subscriptions to both chains",
		Args:    cobra.ExactArgs(5),
		Example: "ebrelayer ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028  wss://ws.s0.b.hmny.io 0x30753E4A8aad7F8597332E813735Def5dD395028 validator",
		RunE:    RunInitRelayerCmd,
	}

	return initRelayerCmd
}

//	generateBindingsCmd : Generates ABIs and bindings for Bridge smart contracts which facilitate contract interaction
func generateBindingsCmd() *cobra.Command {
	generateBindingsCmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generates Bridge smart contracts ABIs and bindings",
		Args:    cobra.ExactArgs(0),
		Example: "generate",
		RunE:    RunGenerateBindingsCmd,
	}

	return generateBindingsCmd
}

// RunInitRelayerCmd executes initRelayerCmd
func RunInitRelayerCmd(cmd *cobra.Command, args []string) error {
	// Load the validator's Ethereum private key from environment variables
	ethereumPrivateKey, err := txs.LoadEthereumPrivateKey()
	if err != nil {
		return errors.Errorf("invalid [ETHEREUM_PRIVATE_KEY] environment variable")
	}

	harmonyPrivateKey, err := txs.LoadHarmonyPrivateKey()
	if err != nil {
		return errors.Errorf("invalid [HARMONY_PRIVATE_KEY] environment variable")
	}

	if !relayer.IsWebsocketURL(args[0]) {
		return errors.Errorf("invalid [web3-provider]: %s", args[0])
	}
	ethereumProvider := args[0]

	if !common.IsHexAddress(args[1]) {
		return errors.Errorf("invalid [bridge-registry-contract-address]: %s", args[1])
	}
	ethereumBridgeRegistry := common.HexToAddress(args[1])

	if !relayer.IsWebsocketURL(args[2]) {
		return errors.Errorf("invalid [hmy-provider]: %s", args[2])
	}
	harmonyProvider := args[2]

	if !common.IsHexAddress(args[3]) {
		return errors.Errorf("invalid [bridge-registry-contract-address]: %s", args[3])
	}
	harmonyBridgeRegistry := common.HexToAddress(args[3])

	if len(strings.Trim(args[4], "")) == 0 {
		return errors.Errorf("invalid [validator-moniker]: %s", args[4])
	}

	validatorMoniker := args[4]

	// Universal logger
	logger := tmLog.NewTMLogger(tmLog.NewSyncWriter(os.Stdout))

	// Initialize new Ethereum event listener
	inBuf := bufio.NewReader(cmd.InOrStdin())

	ethereumSub, err := relayer.NewEthereumSub(inBuf, validatorMoniker, ethereumProvider, harmonyProvider,
		ethereumBridgeRegistry, harmonyBridgeRegistry, ethereumPrivateKey, harmonyPrivateKey, logger)
	if err != nil {
		return err
	}

	harmonySub, err := relayer.NewHarmonySub(inBuf, validatorMoniker, harmonyProvider, ethereumProvider,
		harmonyBridgeRegistry, ethereumBridgeRegistry, harmonyPrivateKey, ethereumPrivateKey, logger)
	if err != nil {
		return err
	}

	go harmonySub.Start()
	go ethereumSub.Start()

	// Exit signal enables graceful shutdown
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	return nil
}

// RunGenerateBindingsCmd : executes the generateBindingsCmd
func RunGenerateBindingsCmd(cmd *cobra.Command, args []string) error {
	ethereumContracts := contract.EthLoadBridgeContracts()

	// Compile contracts, generating contract bins and abis
	err := contract.EthCompileContracts(ethereumContracts)
	if err != nil {
		return err
	}

	// Generate contract bindings from bins and abis
	err = contract.EthGenerateBindings(ethereumContracts)
	if err != nil {
		return err
	}

	harmonyContracts := contract.HmyLoadBridgeContracts()

	// Compile contracts, generating contract bins and abis
	err = contract.HmyCompileContracts(harmonyContracts)
	if err != nil {
		return err
	}

	// Generate contract bindings from bins and abis
	err = contract.HmyGenerateBindings(harmonyContracts)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
