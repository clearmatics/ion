package ionflow

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"log"
	"math/big"
	"os"
	"regexp"
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

// --------------------------------------------------------------------------------
// --------------------------------------------------------------------------------
// GENERIC UTIL FUNCTIONS
// --------------------------------------------------------------------------------
// --------------------------------------------------------------------------------

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
	client bind.ContractTransactor,
	from, to *common.Address,
	amount *big.Int,
	gasLimit uint64,
	payloadBytecode []byte,
) *types.Transaction {

	nonce, err := client.PendingNonceAt(ctx, *from) // uint64(0)
	if err != nil {
		log.Fatal("Error getting pending nonce ", err)
	}
	gasPrice, err := client.SuggestGasPrice(ctx) //new(big.Int)
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
	client bind.ContractTransactor,
	userKey *ecdsa.PrivateKey,
	binStr string,
	abiStr string,
	amount *big.Int,
	gasLimit uint64,
	constructorArgs ...interface{},
) *types.Transaction {
	payload := generateContractPayload(binStr, abiStr, constructorArgs...)
	userAddr := crypto.PubkeyToAddress(userKey.PublicKey)
	tx := newTx(ctx, client, &userAddr, nil, amount, gasLimit, payload)
	signedTx := signTx(tx, userKey)

	err := client.SendTransaction(ctx, signedTx)
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
	client bind.ContractTransactor,
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
	tx := newTx(ctx, client, &from, &to, amount, gasLimit, payload)
	signedTx := signTx(tx, userKey)

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal("ERROR sending contract deployment transaction")
	}
	return signedTx
}

// --------------------------------------------------------------------------------
// --------------------------------------------------------------------------------
// ION SPECIFIC FUNCTIONS
// --------------------------------------------------------------------------------
// --------------------------------------------------------------------------------

// CompileAndDeployIon specific compile and deploy ion contract
func CompileAndDeployIon(
	ctx context.Context,
	client bind.ContractTransactor,
	userKey *ecdsa.PrivateKey,
	chainID interface{},
) <-chan ContractInstance {
	// ---------------------------------------------
	// COMPILE ION AND DEPENDENCIES
	// ---------------------------------------------
	basePath := os.Getenv("GOPATH") + "/src/github.com/clearmatics/ion/contracts/"
	ionContractPath := basePath + "Ion.sol"

	contracts, err := compiler.CompileSolidity("", ionContractPath)
	if err != nil {
		log.Fatal("ERROR failed to compile Ion.sol:", err)
	}

	patriciaTrieContract := contracts[basePath+"libraries/PatriciaTrie.sol:PatriciaTrie"]
	patriciaTrieBinStr, patriciaTrieABIStr := getContractBytecodeAndABI(patriciaTrieContract)

	ionContract := contracts[basePath+"Ion.sol:Ion"]
	ionBinStr, ionABIStr := getContractBytecodeAndABI(ionContract)

	// ---------------------------------------------
	// DEPLOY PATRICIA LIB ADDRESS
	// ---------------------------------------------
	patriciaTrieSignedTx := compileAndDeployContract(
		ctx,
		client,
		userKey,
		patriciaTrieBinStr,
		patriciaTrieABIStr,
		nil,
		uint64(3000000),
	)

	resChan := make(chan ContractInstance)

	// Go-Routine that waits for PatriciaTrie Library and Ion Contract to be deployed
	// Ion depends on PatriciaTrie library
	go func() {
		defer close(resChan)
		deployBackend := client.(bind.DeployBackend)

		// wait for PatriciaTrie library to be deployed
		patriciaTrieAddr, err := bind.WaitDeployed(ctx, deployBackend, patriciaTrieSignedTx)
		if err != nil {
			log.Fatal("ERROR while waiting for contract deployment")
		}

		// ---------------------------------------------
		// DEPLOY ION CONTRACT WITH PATRICIA LIB ADDRESS
		// ---------------------------------------------
		// replace palceholder with Prticia Trie Lib address
		var re = regexp.MustCompile(`__.*__`)
		ionBinStrWithLibAddr := re.ReplaceAllString(ionBinStr, patriciaTrieAddr.Hex()[2:])
		ionSignedTx := compileAndDeployContract(
			ctx,
			client,
			userKey,
			ionBinStrWithLibAddr,
			ionABIStr,
			nil,
			uint64(3000000),
			chainID,
		)

		// only stop blocking the first result after the Ion contract as been deploy
		// this guarantees that it works well with the blockchain simulator Commit()
		resChan <- ContractInstance{patriciaTrieContract, patriciaTrieAddr}

		// wait for Ion to be deployed
		ionAddr, err := bind.WaitDeployed(ctx, deployBackend, ionSignedTx)
		if err != nil {
			log.Fatal("ERROR while waiting for contract deployment")
		}

		resChan <- ContractInstance{ionContract, ionAddr}
	}()

	return resChan
}

// CompileAndDeployValidation method
func CompileAndDeployValidation(
	ctx context.Context,
	client bind.ContractTransactor,
	userKey *ecdsa.PrivateKey,
	chainID interface{},
) <-chan ContractInstance {
	// ---------------------------------------------
	// COMPILE VALIDATION AND DEPENDENCIES
	// ---------------------------------------------
	basePath := os.Getenv("GOPATH") + "/src/github.com/clearmatics/ion/contracts/"
	validationContractPath := basePath + "Validation.sol"

	contracts, err := compiler.CompileSolidity("", validationContractPath)
	if err != nil {
		log.Fatal("ERROR failed to compile Ion.sol:", err)
	}

	validationContract := contracts[basePath+"Validation.sol:Validation"]
	validationBinStr, validationABIStr := getContractBytecodeAndABI(validationContract)

	// ---------------------------------------------
	// DEPLOY VALIDATION CONTRACT
	// ---------------------------------------------
	validationSignedTx := compileAndDeployContract(
		ctx,
		client,
		userKey,
		validationBinStr,
		validationABIStr,
		nil,
		uint64(3000000),
		chainID,
	)

	resChan := make(chan ContractInstance)

	// Go-Routine that waits for PatriciaTrie Library and Ion Contract to be deployed
	// Ion depends on PatriciaTrie library
	go func() {
		defer close(resChan)
		deployBackend := client.(bind.DeployBackend)

		// wait for PatriciaTrie library to be deployed
		validationAddr, err := bind.WaitDeployed(ctx, deployBackend, validationSignedTx)
		if err != nil {
			log.Fatal("ERROR while waiting for contract deployment")
		}
		resChan <- ContractInstance{validationContract, validationAddr}
	}()

	return resChan
}

// CompileAndDeployTriggerVerifierAndConsumerFunction method
func CompileAndDeployTriggerVerifierAndConsumerFunction(
	ctx context.Context,
	client bind.ContractTransactor,
	userKey *ecdsa.PrivateKey,
	ionContractAddress common.Address,
) <-chan ContractInstance {
	// ---------------------------------------------
	// COMPILE VALIDATION AND DEPENDENCIES
	// ---------------------------------------------
	basePath := os.Getenv("GOPATH") + "/src/github.com/clearmatics/ion/contracts/"
	triggerEventVerifierContractPath := basePath + "TriggerEventVerifier.sol"
	consumerFunctionContractPath := basePath + "Function.sol"

	contracts, err := compiler.CompileSolidity("", consumerFunctionContractPath, triggerEventVerifierContractPath)
	if err != nil {
		log.Fatal("ERROR failed to compile Ion.sol:", err)
	}

	triggerEventVerifierContract := contracts[triggerEventVerifierContractPath+":TriggerEventVerifier"]
	triggerEventVerifierBinStr, triggerEventVerifierABIStr := getContractBytecodeAndABI(triggerEventVerifierContract)
	consumerFunctionContract := contracts[consumerFunctionContractPath+":Function"]
	consumerFunctionBinStr, consumerFunctionABIStr := getContractBytecodeAndABI(consumerFunctionContract)

	// ---------------------------------------------
	// DEPLOY TRIGGER EVENT CONTRACT
	// ---------------------------------------------
	triggerEventSignedTx := compileAndDeployContract(
		ctx,
		client,
		userKey,
		triggerEventVerifierBinStr,
		triggerEventVerifierABIStr,
		nil,
		uint64(3000000),
	)

	resChan := make(chan ContractInstance)

	// Go-Routine that waits for PatriciaTrie Library and Ion Contract to be deployed
	// Ion depends on PatriciaTrie library
	go func() {
		defer close(resChan)
		deployBackend := client.(bind.DeployBackend)

		// wait for trigger event contract to be deployed
		triggerEventAddr, err := bind.WaitDeployed(ctx, deployBackend, triggerEventSignedTx)
		if err != nil {
			log.Fatal("ERROR while waiting for contract deployment")
		}

		// ---------------------------------------------
		// DEPLOY CONSUMER FUNCTION CONTRACT
		// ---------------------------------------------
		consumerFunctionSignedTx := compileAndDeployContract(
			ctx,
			client,
			userKey,
			consumerFunctionBinStr,
			consumerFunctionABIStr,
			nil,
			uint64(3000000),
			ionContractAddress,
			triggerEventAddr,
		)

		resChan <- ContractInstance{triggerEventVerifierContract, triggerEventAddr}

		// wait for consumer function contract to be deployed
		consumerFunctionAddr, err := bind.WaitDeployed(ctx, deployBackend, consumerFunctionSignedTx)
		if err != nil {
			log.Fatal("ERROR while waiting for contract deployment")
		}

		resChan <- ContractInstance{consumerFunctionContract, consumerFunctionAddr}
	}()

	return resChan
}
