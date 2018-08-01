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

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

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
	abiPatriciaTrie, err := abi.JSON(strings.NewReader(contractABIStr))
	if err != nil {
		log.Fatal("ERROR reading contract ABI ", err)
	}
	packedABI, err := abiPatriciaTrie.Pack("", constructorArgs...)
	if err != nil {
		log.Fatal("ERROR packing ABI ", err)
	}
	payloadBytecode := append(bytecode, packedABI...)
	return payloadBytecode
}

func generateContractDeployTx(
	ctx context.Context,
	client bind.ContractTransactor,
	from common.Address,
	amount *big.Int,
	gasLimit uint64,
	payloadBytecode []byte,
) *types.Transaction {

	nonce, err := client.PendingNonceAt(ctx, from) // uint64(0)
	if err != nil {
		log.Fatal("Error getting pending nonce ", err)
	}
	gasPrice, err := client.SuggestGasPrice(ctx) //new(big.Int)
	if err != nil {
		log.Fatal("Error suggesting gas price ", err)
	}

	// create contract transaction NewContractCreation is the same has NewTransaction with `to` == nil
	// tx := types.NewTransaction(nonce, nil, amount, gasLimit, gasPrice, payloadBytecode)
	tx := types.NewContractCreation(nonce, amount, gasLimit, gasPrice, payloadBytecode)
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
	tx := generateContractDeployTx(ctx, client, userAddr, amount, gasLimit, payload)
	signedTx := signTx(tx, userKey)

	err := client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal("ERROR sending contract deployment transaction")
	}
	return signedTx
}

// ContractInstance is just an util type to output contract and address
type ContractInstance struct {
	Contract *compiler.Contract
	Address  common.Address
}

// CompileAndDeployIon specific compile and deploy ion contract
func CompileAndDeployIon(ctx context.Context, client bind.ContractTransactor, userKey *ecdsa.PrivateKey) <-chan ContractInstance {
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

	// TODO: make this a go routine
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
		constructorArg1Ion := crypto.Keccak256Hash([]byte("test argument")) // Ion argument
		ionSignedTx := compileAndDeployContract(
			ctx,
			client,
			userKey,
			ionBinStrWithLibAddr,
			ionABIStr,
			nil,
			uint64(3000000),
			constructorArg1Ion,
		)

		// only stop blocking the first result after the Ion contract as been deploy
		// this garuantees that it works well with the blockchain simulator Commit()
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
