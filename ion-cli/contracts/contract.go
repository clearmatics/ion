// Copyright (c) 2018 Clearmatics Technologies Ltd
package contract

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"encoding/json"
	"log"
	"math/big"
	"os"
	"strings"
	"errors"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/clearmatics/autonity/common/compiler"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// ContractInstance is just an util type to output contract and address
type ContractInstance struct {
	Contract *compiler.Contract
	Abi *abi.ABI
	Path string
}

// GENERIC UTIL FUNCTIONS

func GetContractBytecodeAndABI(c *compiler.Contract) (string, string) {
	cABIBytes, err := json.Marshal(c.Info.AbiDefinition)
	if err != nil {
		log.Fatal("ERROR marshalling contract ABI:", err)
	}

	contractBinStr := c.Code[2:]
	contractABIStr := string(cABIBytes)
	return contractBinStr, contractABIStr
}

func generateContractPayload(contractBinStr string, contractABIStr string, constructorArgs ...interface{}) []byte {
	bytecode := common.Hex2Bytes(contractBinStr)
	abiContract, err := abi.JSON(strings.NewReader(contractABIStr))
	if err != nil {
		log.Fatal("ERROR reading contract ABI ", err)
	}
	packedABI, err := abiContract.Pack("", constructorArgs...)
	if err != nil {
		log.Fatal("ERROR packing ABI ", err)
	}
	payloadBytecode := append(bytecode, packedABI...)
	return payloadBytecode
}

func newTx(
	ctx context.Context,
	backend bind.ContractBackend,
	from, to *common.Address,
	amount *big.Int,
	gasLimit uint64,
	payloadBytecode []byte,
) *types.Transaction {

	nonce, err := backend.PendingNonceAt(ctx, *from) // uint64(0)
	if err != nil {
		log.Fatal("Error getting pending nonce ", err)
	}
	gasPrice, err := backend.SuggestGasPrice(ctx) //new(big.Int)
	if err != nil {
		log.Fatal("Error suggesting gas price ", err)
	}

	// create contract transaction NewContractCreation is the same has NewTransaction with `to` == nil
	// tx := types.NewTransaction(nonce, nil, amount, gasLimit, gasPrice, payloadBytecode)
	var tx *types.Transaction
	if to == nil {
		tx = types.NewContractCreation(nonce, amount, gasLimit, gasPrice, payloadBytecode)
	} else {
		tx = types.NewTransaction(nonce, *to, amount, gasLimit, gasPrice, payloadBytecode)
	}
	return tx
}

// method created just to easily sign a tranasaction
func signTx(tx *types.Transaction, userKey *ecdsa.PrivateKey) *types.Transaction {
	signer := types.HomesteadSigner{} // this functions makes it easier to change signer if needed
	signedTx, err := types.SignTx(tx, signer, userKey)
	if err != nil {
		log.Fatal("Error signing tx: ", err)
	}
	return signedTx
}

func DeployContract(
	ctx context.Context,
	backend bind.ContractBackend,
	userKey *ecdsa.PrivateKey,
	payload []byte,
	amount *big.Int,
	gasLimit uint64,
) (*types.Transaction, error){
	userAddr := crypto.PubkeyToAddress(userKey.PublicKey)
	tx := newTx(ctx, backend, &userAddr, nil, amount, gasLimit, payload)
	signedTx := signTx(tx, userKey)

	err := backend.SendTransaction(ctx, signedTx)
	if err != nil {
	    return nil, err
	}
	return signedTx, nil
}

func CompilePayload(
	binStr string,
	abiStr string,
	constructorArgs ...interface{},
) ([]byte){
	return generateContractPayload(binStr, abiStr, constructorArgs...)
}

// CallContract without changing the state
func CallContract(
	ctx context.Context,
	client bind.ContractCaller,
	contract *compiler.Contract,
	from, to common.Address,
	methodName string,
	out interface{},
	args ...interface{},
) (res interface{}, err error) {
	abiStr, err := json.Marshal(contract.Info.AbiDefinition)
	if err != nil {
		return nil, err
	}

	abiContract, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		return nil, err
	}

	input, err := abiContract.Pack(methodName, args...)
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{From: from, To: &to, Data: input}
	output, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, err
	}
	err = abiContract.Unpack(out, methodName, output)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// TransactionContract execute function in contract
func TransactionContract(
	ctx context.Context,
	backend bind.ContractBackend,
	userKey *ecdsa.PrivateKey,
	contract *compiler.Contract,
	to common.Address,
	amount *big.Int,
	gasLimit uint64,
	methodName string,
	args ...interface{},
) (*types.Transaction, error) {

    fmt.Print("Marshalling ABI\n")
	abiStr, err := json.Marshal(contract.Info.AbiDefinition)
	if err != nil {
	    errStr := fmt.Sprintf("ERROR marshalling abi to string: %s\n", err)
	    return nil, errors.New(errStr)
		log.Fatal()
	}


    fmt.Print("JSONify ABI\n")
	abiContract, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
	    errStr := fmt.Sprintf("ERROR reading contract ABI: %s\n", err)
	    return nil, errors.New(errStr)
	}


    fmt.Print("Packing Args to ABI\n")
	payload, err := abiContract.Pack(methodName, args...)
	if err != nil {
	    errStr := fmt.Sprintf("ERROR packing the method name for the contract call: %s\n", err)
	    return nil, errors.New(errStr)
	}


    fmt.Print("Retrieving public key\n")
	from := crypto.PubkeyToAddress(userKey.PublicKey)

    fmt.Print("Creating transaction\n")
	tx := newTx(ctx, backend, &from, &to, amount, gasLimit, payload)

    fmt.Print("Signing transaction\n")
	signedTx := signTx(tx, userKey)

    fmt.Print("SENDING TRANSACTION\n")

	err = backend.SendTransaction(ctx, signedTx)
	if err != nil {
	    errStr := fmt.Sprintf("ERROR sending transaction: %s\n", err)
	    return nil, errors.New(errStr)
	}
	return signedTx, nil
}

func CompileContract(contract string) (compiledContract *compiler.Contract, err error) {
	basePath := os.Getenv("GOPATH") + "/src/github.com/clearmatics/ion/contracts/"
	contractPath := basePath + contract + ".sol"

	contracts, err := compiler.CompileSolidity("", []string{}, contractPath)
	if err != nil {
	    return nil, err
	}

	compiledContract = contracts[basePath+contract+".sol:"+contract]

	return compiledContract, nil
}

func CompileContractAt(contractPath string) (compiledContract *compiler.Contract, err error) {
	path := strings.Split(contractPath, "/")
	contractName := path[len(path)-1]
	contractFolder := path[len(path)-2]

	i := strings.Index(contractPath, contractFolder)
	remapping := fmt.Sprintf("../=%s", contractPath[:i])

    contract, err := compiler.CompileSolidity("", []string{remapping}, contractPath)
	if err != nil {
	    return nil, err
	}

    compiledContract = contract[contractPath+":"+strings.Replace(contractName, ".sol", "", -1)]

	return compiledContract, nil
}

func CompileContractWithLibraries(contractPath string, libraries map[string]common.Address) (compiledContract *compiler.Contract, err error) {
	path := strings.Split(contractPath, "/")
	contractName := path[len(path)-1]
	contractFolder := path[len(path)-2]

    args := []string{}

    for name := range libraries {
        address := libraries[name]

        libraryArg := name + ":" + address.String()
        args = append(args, fmt.Sprintf("--libraries=%s", libraryArg))
    }

	i := strings.Index(contractPath, contractFolder)
	args = append(args, fmt.Sprintf("../=%s ", contractPath[:i]))

    contract, err := compiler.CompileSolidity("", args, contractPath)
	if err != nil {
	    return nil, err
	}

    compiledContract = contract[contractPath+":"+strings.Replace(contractName, ".sol", "", -1)]

	return compiledContract, nil
}
