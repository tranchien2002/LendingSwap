package contract

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const (
	SolcCmdText   = "[SOLC_CMD]"
	DirectoryText = "[DIRECTORY]"
	ContractText  = "[CONTRACT]"
)

var (
	// EthBaseABIBINGenCmd is the base command for contract compilation to ABI and BIN
	EthBaseABIBINGenCmd = strings.Join([]string{"solc ",
		fmt.Sprintf("--%s ./ethereum-contracts/contracts/%s%s.sol ", SolcCmdText, DirectoryText, ContractText),
		fmt.Sprintf("-o ./cmd/ebrelayer/contract/generated/ethereum/%s/%s ", SolcCmdText, ContractText),
		"--overwrite ",
		"--allow-paths *,"},
		"")
	// EthBaseBindingGenCmd is the base command for contract binding generation
	EthBaseBindingGenCmd = strings.Join([]string{"abigen ",
		fmt.Sprintf("--bin ./cmd/ebrelayer/contract/generated/ethereum/bin/%s/%s.bin ", ContractText, ContractText),
		fmt.Sprintf("--abi ./cmd/ebrelayer/contract/generated/ethereum/abi/%s/%s.abi ", ContractText, ContractText),
		fmt.Sprintf("--pkg %s ", ContractText),
		fmt.Sprintf("--type %s ", ContractText),
		fmt.Sprintf("--out ./cmd/ebrelayer/contract/generated/ethereum/bindings/%s/%s.go", ContractText, ContractText)},
		"")
	// HmyBaseABIBINGenCmd is the base command for contract compilation to ABI and BIN
	HmyBaseABIBINGenCmd = strings.Join([]string{"solc ",
		fmt.Sprintf("--%s ./harmony-contracts/contracts/%s%s.sol ", SolcCmdText, DirectoryText, ContractText),
		fmt.Sprintf("-o ./cmd/ebrelayer/contract/generated/harmony/%s/%s ", SolcCmdText, ContractText),
		"--overwrite ",
		"--allow-paths *,"},
		"")
	// HmyBaseBindingGenCmd is the base command for contract binding generation
	HmyBaseBindingGenCmd = strings.Join([]string{"abigen ",
		fmt.Sprintf("--bin ./cmd/ebrelayer/contract/generated/harmony/bin/%s/%s.bin ", ContractText, ContractText),
		fmt.Sprintf("--abi ./cmd/ebrelayer/contract/generated/harmony/abi/%s/%s.abi ", ContractText, ContractText),
		fmt.Sprintf("--pkg %s ", ContractText),
		fmt.Sprintf("--type %s ", ContractText),
		fmt.Sprintf("--out ./cmd/ebrelayer/contract/generated/harmony/bindings/%s/%s.go", ContractText, ContractText)},
		"")
)

// EthCompileContracts compiles contracts to BIN and ABI files
func EthCompileContracts(contracts BridgeContracts) error {
	for _, contract := range contracts {
		// Construct generic BIN/ABI generation cmd with contract's directory path and name
		baseDirectory := ""
		if contract.String() == BridgeBank.String() {
			baseDirectory = contract.String() + "/"
		}
		dirABIBINGenCmd := strings.Replace(EthBaseABIBINGenCmd, DirectoryText, baseDirectory, -1)
		contractABIBINGenCmd := strings.Replace(dirABIBINGenCmd, ContractText, contract.String(), -1)

		// Segment BIN and ABI generation cmds
		contractBINGenCmd := strings.Replace(contractABIBINGenCmd, SolcCmdText, "bin", -1)
		err := execCmd(contractBINGenCmd)
		if err != nil {
			return err
		}

		contractABIGenCmd := strings.Replace(contractABIBINGenCmd, SolcCmdText, "abi", -1)
		err = execCmd(contractABIGenCmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// HmyCompileContracts compiles contracts to BIN and ABI files
func HmyCompileContracts(contracts BridgeContracts) error {
	for _, contract := range contracts {
		// Construct generic BIN/ABI generation cmd with contract's directory path and name
		baseDirectory := ""
		if contract.String() == BridgeBank.String() {
			baseDirectory = contract.String() + "/"
		}
		dirABIBINGenCmd := strings.Replace(HmyBaseABIBINGenCmd, DirectoryText, baseDirectory, -1)
		contractABIBINGenCmd := strings.Replace(dirABIBINGenCmd, ContractText, contract.String(), -1)

		// Segment BIN and ABI generation cmds
		contractBINGenCmd := strings.Replace(contractABIBINGenCmd, SolcCmdText, "bin", -1)
		err := execCmd(contractBINGenCmd)
		if err != nil {
			return err
		}

		contractABIGenCmd := strings.Replace(contractABIBINGenCmd, SolcCmdText, "abi", -1)
		err = execCmd(contractABIGenCmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// EthGenerateBindings generates bindings for each ethereum contract
func EthGenerateBindings(contracts BridgeContracts) error {
	for _, contract := range contracts {
		genBindingCmd := strings.Replace(EthBaseBindingGenCmd, ContractText, contract.String(), -1)
		err := execCmd(genBindingCmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// HmyGenerateBindings generates bindings for each ethereum contract
func HmyGenerateBindings(contracts BridgeContracts) error {
	for _, contract := range contracts {
		genBindingCmd := strings.Replace(HmyBaseBindingGenCmd, ContractText, contract.String(), -1)
		err := execCmd(genBindingCmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// execCmd executes a bash cmd
func execCmd(cmd string) error {

	mainCmd := exec.Command("sh", "-c", cmd)

	var out bytes.Buffer
	var stderr bytes.Buffer

	mainCmd.Stdout = &out
	mainCmd.Stderr = &stderr

	err := mainCmd.Run()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	fmt.Println("Result: Successfully!")

	return nil
}
