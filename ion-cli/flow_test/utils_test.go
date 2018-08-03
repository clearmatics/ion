package ionflow

import (
	"bytes"
	"context"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

func TestRegisterChain(t *testing.T) {
	// check comments on TestCompileAndDeploy()
	ctx := context.Background()
	initialBalance := big.NewInt(1000000000)
	userAKey, _ := crypto.GenerateKey()
	userAAddr := crypto.PubkeyToAddress(userAKey.PublicKey)
	alloc := make(core.GenesisAlloc)
	alloc[userAAddr] = core.GenesisAccount{
		Balance: initialBalance,
	}
	blockchain := backends.NewSimulatedBackend(alloc)
	var chainID [32]byte
	copy(chainID[:], crypto.Keccak256Hash([]byte("DEPLOYEDCHAINID")).Bytes())
	contractChan := CompileAndDeployIon(ctx, blockchain, userAKey, chainID)
	blockchain.Commit()
	<-contractChan
	blockchain.Commit()
	ionContractInstance := <-contractChan

	// deploy validation contract
	contractChan = CompileAndDeployValidation(ctx, blockchain, userAKey, chainID)
	blockchain.Commit()
	validationContractInstance := <-contractChan

	var chainIDA, addressArray [32]byte
	copy(addressArray[:], validationContractInstance.Address.Bytes())
	copy(chainIDA[:], crypto.Keccak256Hash([]byte("TESTCHAINID")).Bytes())
	txRegisterChain := TransactionContract(
		ctx,
		blockchain,
		userAKey,
		ionContractInstance.Contract,
		ionContractInstance.Address,
		nil,
		uint64(3000000),
		"RegisterChain",
		chainIDA,
		addressArray,
	)
	blockchain.Commit()

	registerChainReceipt, err := bind.WaitMined(ctx, blockchain, txRegisterChain)
	if err != nil {
		log.Fatal("ERROR while waiting for contract deployment")
	}
	if registerChainReceipt.Status == 0 {
		t.Fatalf("ERROR transaction of RegisterChain failed!: %#v\n", registerChainReceipt)
	}

	methodName := "chains"
	var isChainRegistered bool
	CallContract(
		ctx,
		blockchain,
		ionContractInstance.Contract,
		userAAddr,
		ionContractInstance.Address,
		methodName,
		&isChainRegistered,
		chainIDA,
	)

	if !isChainRegistered {
		t.Log("ERROR expecting value of chains(validation.address) to be true, but it was ", isChainRegistered)
	}
}
