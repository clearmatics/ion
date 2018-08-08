package contract

import (
	"context"
	"crypto/ecdsa"
	"log"
	"os"
	"regexp"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/compiler"
)

// CompileAndDeployIon specific compile and deploy ion contract
func CompileAndDeployIon(
	ctx context.Context,
	client bind.ContractBackend,
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
