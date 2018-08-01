package ionflow

import (
	"bytes"
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

	// create a chain id
	chainID := crypto.Keccak256Hash([]byte("test argument")) // Ion argument

	// start compile and deploy ion
	contractChan := CompileAndDeployIon(ctx, blockchain, userAKey, chainID)

	// commit first block after sent transaction for deployment of patricia trie lib
	blockchain.Commit()
	// patriciaTrieContractInstance := <-contractChan
	<-contractChan

	// commit after transaction for deployment of Ion (with reference to Patricia Trie lib) as been sent
	blockchain.Commit()
	ionContractInstance := <-contractChan

	// call contract variable
	methodName := "chainId"
	out := new([32]byte)
	CallContract(ctx, blockchain, ionContractInstance.Contract, userAAddr, ionContractInstance.Address, methodName, out)

	if !bytes.Equal((*out)[:], chainID.Bytes()) {
		t.Fatal("ERROR chainID result from contract call, and sent to contract constructor differ")
	}
}
