// Copyright (c) 2018 Clearmatics Technologies Ltd
package contract

import (
	"bytes"
	"context"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func Test_CompileAndDeployIon(t *testing.T) {
	// ---------------------------------------------
	// START BLOCKCHAIN SIMULATOR
	// ---------------------------------------------

	ctx := context.Background()
	initialBalance := big.NewInt(1000000000)

	userAKey, _ := crypto.GenerateKey()
	userAddr := crypto.PubkeyToAddress(userAKey.PublicKey)

	// start simulated blockchain
	alloc := make(core.GenesisAlloc)
	alloc[userAddr] = core.GenesisAccount{
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
	CallContract(ctx, blockchain, ionContractInstance.Contract, userAddr, ionContractInstance.Address, methodName, out)

	if !bytes.Equal((*out)[:], chainID.Bytes()) {
		t.Fatal("ERROR chainID result from contract call, and sent to contract constructor differ")
	}
}

func Test_RegisterChain(t *testing.T) {
	// ---------------------------------------------
	// HARD CODED DATA
	// ---------------------------------------------
	testValidators := []common.Address{
		common.HexToAddress("0x42eb768f2244c8811c63729a21a3569731535f06"),
		common.HexToAddress("0x6635f83421bf059cd8111f180f0727128685bae4"),
		common.HexToAddress("0x7ffc57839b00206d1ad20c69a1981b489f772031"),
		common.HexToAddress("0xb279182d99e65703f0076e4812653aab85fca0f0"),
		common.HexToAddress("0xd6ae8250b8348c94847280928c79fb3b63ca453e"),
		common.HexToAddress("0xda35dee8eddeaa556e4c26268463e26fb91ff74f"),
		common.HexToAddress("0xfc18cbc391de84dbd87db83b20935d3e89f5dd91"),
	}

	// check comments on TestCompileAndDeploy()
	ctx := context.Background()
	initialBalance := big.NewInt(1000000000)
	userAKey, _ := crypto.GenerateKey()
	userAddr := crypto.PubkeyToAddress(userAKey.PublicKey)
	alloc := make(core.GenesisAlloc)
	alloc[userAddr] = core.GenesisAccount{
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
	var ionContractAddr [20]byte
	copy(ionContractAddr[:], ionContractInstance.Address.Bytes())
	contractChan = CompileAndDeployValidation(ctx, blockchain, userAKey, chainID, ionContractAddr)
	blockchain.Commit()
	validationContractInstance := <-contractChan

	var chainIDA [32]byte
	var validationAddress [20]byte
	copy(validationAddress[:], validationContractInstance.Address.Bytes())
	copy(chainIDA[:], crypto.Keccak256Hash([]byte("TESTCHAINID")).Bytes())
	deployedChainID := common.HexToHash("0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075")
	txRegisterChain := RegisterChain(
		ctx,
		blockchain,
		userAKey,
		validationContractInstance.Contract,
		validationContractInstance.Address,
		chainIDA,
		testValidators,
		deployedChainID,
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
		validationContractInstance.Contract,
		userAddr,
		validationContractInstance.Address,
		methodName,
		&isChainRegistered,
		chainID,
	)

	if !isChainRegistered {
		t.Log("ERROR expecting value of chains(validation.address) to be true, but it was ", isChainRegistered)
	}
}
