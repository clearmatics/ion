// Copyright (c) 2018 Clearmatics Technologies Ltd
package contract

import (
	"context"
	"crypto/ecdsa"
	"log"
	"os"
	"strings"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/core/types"
)

// CompileAndDeployValidation method
func CompileAndDeployValidation(
	ctx context.Context,
	client bind.ContractBackend,
	userKey *ecdsa.PrivateKey,
	chainID interface{},
	ionContractAddress common.Address,
) <-chan ContractInstance {
	// ---------------------------------------------
	// COMPILE VALIDATION AND DEPENDENCIES
	// ---------------------------------------------
	basePath := os.Getenv("GOPATH") + "/src/github.com/clearmatics/ion/contracts/"
	validationContractPath := basePath + "Validation.sol"

	contracts, err := compiler.CompileSolidity("", validationContractPath)
	if err != nil {
		log.Fatal("ERROR failed to compile Validation.sol:", err)
	}

	validationContract := contracts[basePath+"Validation.sol:Validation"]
	validationBinStr, validationABIStr := GetContractBytecodeAndABI(validationContract)

	// ---------------------------------------------
	// DEPLOY VALIDATION CONTRACT
	// ---------------------------------------------
	validationSignedTx := CompileAndDeployContract(
		ctx,
		client,
		userKey,
		validationBinStr,
		validationABIStr,
		nil,
		uint64(3000000),
		chainID,
		ionContractAddress,
	)

	resChan := make(chan ContractInstance)

	// Go-Routine that waits for PatriciaTrie Library and Ion Contract to be deployed
	// Ion depends on PatriciaTrie library
	go func() {
		defer close(resChan)
		deployBackend := client.(bind.DeployBackend)

		// wait for PatriciaTrie library to be deployed
		_, err := bind.WaitDeployed(ctx, deployBackend, validationSignedTx)
		if err != nil {
			log.Fatal("ERROR while waiting for contract deployment")
		}
		abistruct, err := abi.JSON(strings.NewReader(validationABIStr))
        if err != nil {
		    log.Fatal("ERROR failed to compile Validation.sol:", err)
        }
		resChan <- ContractInstance{validationContract, &abistruct}
	}()

	return resChan
}

// RegisterChain with Validation contract specified
func RegisterChain(
	ctx context.Context,
	backend bind.ContractBackend,
	userKey *ecdsa.PrivateKey,
	contract *compiler.Contract,
	toAddr common.Address,
	chainID common.Hash,
	validators []common.Address,
	registerHash common.Hash,
) (tx *types.Transaction) {
	tx, err := TransactionContract(
		ctx,
		backend,
		userKey,
		contract,
		toAddr,
		nil,
		uint64(3000000),
		"RegisterChain",
		chainID,
		validators,
		registerHash,
	)

	if err != nil {
	    fmt.Println(err)
	}

	return
}

// SubmitBlock Submits block to Validation contract specified
func SubmitBlock(
	ctx context.Context,
	backend bind.ContractBackend,
	userKey *ecdsa.PrivateKey,
	contract *compiler.Contract,
	toAddr common.Address,
	chainID common.Hash,
	unsignedBlockHeaderRLP []byte,
	signedBlockHeaderRLP []byte,
) (tx *types.Transaction) {
	tx, err := TransactionContract(
		ctx,
		backend,
		userKey,
		contract,
		toAddr,
		nil,
		uint64(3000000),
		"SubmitBlock",
		chainID,
		unsignedBlockHeaderRLP,
		signedBlockHeaderRLP,
	)

	if err != nil {
	    fmt.Println(err)
	}

	return
}

// ValidBlock Queries validation contract to see is block is valid
func ValidBlock(
	ctx context.Context,
	backend bind.ContractBackend,
	contract *compiler.Contract,
	userAddr common.Address,
	toAddr common.Address,
	chainID common.Hash,
	blockHash common.Hash,
) (isBlockValid bool) {
	methodName := "m_blockhashes"
	CallContract(
		ctx,
		backend,
		contract,
		userAddr,
		toAddr,
		methodName,
		&isBlockValid,
		chainID,
		blockHash,
	)
	return
}

// LatestValidBlock Queries validation contract to see is block is valid
func LatestValidBlock(
	ctx context.Context,
	backend bind.ContractBackend,
	contract *compiler.Contract,
	userAddr common.Address,
	toAddr common.Address,
	chainID common.Hash,
) (latestBlock common.Hash) {
	methodName := "m_latestblock"
	CallContract(
		ctx,
		backend,
		contract,
		userAddr,
		toAddr,
		methodName,
		&latestBlock,
		chainID,
	)
	return
}
