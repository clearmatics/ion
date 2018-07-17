package contract

import (
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

// Initialisation array
var val1 = common.HexToAddress("0x8671e5e08d74f338ee1c462340842346d797afd3")
var val2 = common.HexToAddress("0x8671e5e08d74f338ee1c462340842346d797afd3")
var initValidators = []common.Address{val1, val2}

var GENESISHASH = "c3bac257bbd04893316a76d41b6ff70de5f65c9f24db128864a6322d8e0e2f28"

// Test validation contract gets deployed correctly
func Test_DeployValidation(t *testing.T) {

	// Setup simulated block chain
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1000000000)}
	blockchain := backends.NewSimulatedBackend(alloc)

	genesisHash := [32]byte{}
	copy(genesisHash[:], []byte(GENESISHASH))

	// Deploy contract
	address, _, _, err := DeployValidation(
		auth,
		blockchain,
		initValidators,
		genesisHash,
	)
	// commit all pending transactions
	blockchain.Commit()

	if err != nil {
		t.Fatalf("Failed to deploy the Validation contract: %v", err)
	}

	if len(address.Bytes()) == 0 {
		t.Error("Expected a valid deployment address. Received empty address byte array instead")
	}

}

// Test inbox contract gets deployed correctly
func TestGetValidators(t *testing.T) {

	//Setup simulated block chain
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1000000000)}
	blockchain := backends.NewSimulatedBackend(alloc)

	genesisHash := [32]byte{}
	copy(genesisHash[:], []byte(GENESISHASH))

	// Deploy contract
	_, _, contract, _ := DeployValidation(
		auth,
		blockchain,
		initValidators,
		genesisHash,
	)
	// commit all pending transactions
	blockchain.Commit()

	validators, _ := contract.GetValidators(&bind.CallOpts{})
	assert.Equal(t, validators[0], val1)
	assert.Equal(t, validators[1], val2)

}

// Test that the latest block returns the last submitted block hash
func Test_LatestBlock(t *testing.T) {

	//Setup simulated block chain
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1000000000)}
	blockchain := backends.NewSimulatedBackend(alloc)

	genesisHash := [32]byte{}
	copy(genesisHash[:], []byte(GENESISHASH))

	// Deploy contract
	_, _, contract, _ := DeployValidation(
		auth,
		blockchain,
		initValidators,
		genesisHash,
	)
	// commit all pending transactions
	blockchain.Commit()

	latestBlock, _ := contract.LatestBlock(&bind.CallOpts{})
	assert.Equal(t, latestBlock, genesisHash)
}

// Test that the Validators is updated upon deployment
func Test_Validators(t *testing.T) {

	//Setup simulated block chain
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1000000000)}
	blockchain := backends.NewSimulatedBackend(alloc)

	genesisHash := [32]byte{}
	copy(genesisHash[:], []byte(GENESISHASH))

	// Deploy contract
	_, _, contract, _ := DeployValidation(
		auth,
		blockchain,
		initValidators,
		genesisHash,
	)
	// commit all pending transactions
	blockchain.Commit()

	idx := big.NewInt(0)
	validator, err := contract.Validators(&bind.CallOpts{}, idx)
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v", err)
	}
	assert.Equal(t, initValidators[0], validator)
}
