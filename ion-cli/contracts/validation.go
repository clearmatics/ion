// Copyright (c) 2018 Clearmatics Technologies Ltd
package contract

import (
	"context"
	"crypto/ecdsa"
	"log"
	"os"

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

// Registers chain with Validation contract specified
func RegisterChain(
	ctx context.Context,
	backend bind.ContractBackend,
	userKey *ecdsa.PrivateKey,
	contract *compiler.Contract,
	toAddr common.Address,
	chainId common.Hash,
	ionAddr common.Address,
	validators []common.Address,
	registerHash common.Hash,
) (tx *types.Transaction) {
	tx = TransactionContract(
		ctx,
		backend,
		userKey,
		contract,
		toAddr,
		nil,
		uint64(3000000),
		"RegisterChain",
		chainId,
		ionAddr,
		validators,
		registerHash,
	)

	return
}

// Submits block to Validation contract specified
func SubmitBlock(
	ctx context.Context,
	backend bind.ContractBackend,
	userKey *ecdsa.PrivateKey,
	contract *compiler.Contract,
	toAddr common.Address,
	chainId common.Hash,
	unsignedBlockHeaderRLP []byte,
	signedBlockHeaderRLP []byte,
) (tx *types.Transaction) {
	tx = TransactionContract(
		ctx,
		backend,
		userKey,
		contract,
		toAddr,
		nil,
		uint64(3000000),
		"SubmitBlock",
		chainId,
		unsignedBlockHeaderRLP,
		signedBlockHeaderRLP,
	)
	return
}

// Queries validation contract to see is block is valid
func ValidBlock(
	ctx context.Context,
	backend bind.ContractBackend,
	contract *compiler.Contract,
	userAddr common.Address,
	toAddr common.Address,
	chainId common.Hash,
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
		chainId,
		blockHash,
	)
	return
}

// Queries validation contract to see is block is valid
func LatestValidBlock(
	ctx context.Context,
	backend bind.ContractBackend,
	contract *compiler.Contract,
	userAddr common.Address,
	toAddr common.Address,
	chainId common.Hash,
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
		chainId,
	)
	return
}
