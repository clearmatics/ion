package contract

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"log"
	"math/big"
	"os"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// ContractInstance is just an util type to output contract and address
type ContractInstance struct {
	Contract *compiler.Contract
	Address  common.Address
}

// GENERIC UTIL FUNCTIONS

func getContractBytecodeAndABI(c *compiler.Contract) (string, string) {
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

func compileAndDeployContract(
	ctx context.Context,
	backend bind.ContractBackend,
	userKey *ecdsa.PrivateKey,
	binStr string,
	abiStr string,
	amount *big.Int,
	gasLimit uint64,
	constructorArgs ...interface{},
) *types.Transaction {
	payload := generateContractPayload(binStr, abiStr, constructorArgs...)
	userAddr := crypto.PubkeyToAddress(userKey.PublicKey)
	tx := newTx(ctx, backend, &userAddr, nil, amount, gasLimit, payload)
	signedTx := signTx(tx, userKey)

	err := backend.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal("ERROR sending contract deployment transaction")
	}
	return signedTx
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
) {
	abiStr, err := json.Marshal(contract.Info.AbiDefinition)
	if err != nil {
		log.Fatal("ERROR marshalling abi to string", err)
	}

	abiContract, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatal("ERROR reading contract ABI ", err)
	}

	input, err := abiContract.Pack(methodName, args...)
	if err != nil {
		log.Fatal("ERROR packing the method name for the contract call: ", err)
	}
	msg := ethereum.CallMsg{From: from, To: &to, Data: input}
	output, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		log.Fatal("ERROR calling the Ion Contract", err)
	}
	err = abiContract.Unpack(out, methodName, output)
	if err != nil {
		log.Fatal("ERROR upacking the call: ", err)
	}
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
) *types.Transaction {
	abiStr, err := json.Marshal(contract.Info.AbiDefinition)
	if err != nil {
		log.Fatal("ERROR marshalling abi to string", err)
	}

	abiContract, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatal("ERROR reading contract ABI ", err)
	}

	payload, err := abiContract.Pack(methodName, args...)
	if err != nil {
		log.Fatal("ERROR packing the method name for the contract call: ", err)
	}

	from := crypto.PubkeyToAddress(userKey.PublicKey)
	tx := newTx(ctx, backend, &from, &to, amount, gasLimit, payload)
	signedTx := signTx(tx, userKey)

	err = backend.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal("ERROR sending transaction", err)
	}
	return signedTx
}

func CompileContract(contract string) (compiledContract *compiler.Contract) {
	basePath := os.Getenv("GOPATH") + "/src/github.com/clearmatics/ion/contracts/"
	contractPath := basePath + contract + ".sol"

	contracts, err := compiler.CompileSolidity("", contractPath)
	if err != nil {
		log.Fatal("ERROR failed to compile contract:", err)
	}

	compiledContract = contracts[basePath+contract+".sol:"+contract]

	return
}
