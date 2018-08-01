package ionflow

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestCompileAndDeployIon(t *testing.T) {
	// ---------------------------------------------
	// START BLOCKCHAIN SIMULATOR
	// ---------------------------------------------

	ctx := context.Background()
	initialBalance := big.NewInt(1000000000)

	userAKey, _ := crypto.GenerateKey()
	userAAddr := crypto.PubkeyToAddress(userAKey.PublicKey)

	// start simulated blockchain
	alloc := make(core.GenesisAlloc)
	alloc[userAAddr] = core.GenesisAccount{
		Balance: initialBalance,
	}
	blockchain := backends.NewSimulatedBackend(alloc)

	chainID := crypto.Keccak256Hash([]byte("test argument")) // Ion argument
	contractChan := CompileAndDeployIon(ctx, blockchain, userAKey, chainID)
	blockchain.Commit()
	patriciaTrieContractInstance := <-contractChan
	t.Log(patriciaTrieContractInstance)

	blockchain.Commit()
	ionContractInstance := <-contractChan
	t.Log(ionContractInstance)
}
