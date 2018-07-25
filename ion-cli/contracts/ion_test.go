package contract

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/clearmatics/ion/ion-cli/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

var CLIENT = "http://127.0.0.1:8501"
var KEY = `{"address":"2be5ab0e43b6dc2908d5321cf318f35b80d0c10d","crypto":{"cipher":"aes-128-ctr","ciphertext":"0b11aa865046778a1b16a9b8cb593df704e3fe09f153823d75442ad1aab66caa","cipherparams":{"iv":"4aa66b789ee2d98cf77272a72eeeaa50"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"b957fa7b7577240fd3791168bbe08903af4c8cc62c304f1df072dc2a59b1765e"},"mac":"197a06eb0449301d871400a6bdf6c136b6f7658ee41e3f2f7fd81ca11cd954a3"},"id":"a3cc1eae-3e36-4659-b759-6cf416216e72","version":3}`
var IONADDR = common.HexToAddress("0x6aa4444974f60bf3a0bf074d3c194f88ae4d4613")

var DEPLOYEDCHAINID = `ab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075`
var TESTCHAINID = `22b55e8a4f7c03e1689da845dd463b09299cb3a574e64c68eafc4e99077a7254`
var TESTBLOCK = `{"difficulty": "2","extraData": "0xd88301080b846765746888676f312e31302e32856c696e757800000000000000dd2ba07230e2186ee83ef77d88298c068205167718d48ba5b6ba1de552d0c6ce156011a58b49ed91855de154346968a7eeaaf20914022e58e4f6c0e1e02567ec00", "gasLimit": "5635559972940396", "gasUsed": "273138", "hash": "0x6f98a4b7bffb6c5b3dce3923be8a87eeef94ba22e3266cfcfd53407e70294fa4", "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", "miner":"0x0000000000000000000000000000000000000000","mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000","nonce":"0x0000000000000000","number": "5446","parentHash": "0xaa912ad61a8aa3e2d1144e4c76b746720e41682122a8b77eff890099a0ff6284","receiptsRoot": "0x1d000ef3f5ca9ebc62cc8aaa07e8fbd103583d1e3cbd28c13e62bc8eac5eb2f1","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size": "2027","stateRoot":"0xb347dd25d9a8a456448aed25e072c9db54f464be5e3ce1f505cc171cacf3a967","timestamp": "1531327572","totalDifficulty": "10893","transactions": [ "0x63eff998322fd9ec22bbe141ea74ab929197d2db65834e6f4db65743a214cea3","0xa581c3669e5c927e624949d378a5a9df949d4e7f15e1e974c754929408e4b8a5","0x51f1e414334270b7a338f4d81eb82a5560b406f992bf1b3a2371964425e7c0d8","0xc199cd22b3285ea30d798204c3c2fdb8cebfb4648589aa9687aecd9296705ff6","0x4da9368a70e4cfcee28f4c95d69d1256a7d649505f6971b0435bc90f963833f8","0x3cd690a88f4eff005e85f12492afe84344355e9913ea391e52cc0c39debc19e1","0x5dc2e7ea90a0b2630c8138d1357c78ec3d0f55ed23d2951f3c3754ccb9d47446","0xc7f92719dd9f10e8e49ce31a1d271a268269f5c6103629b65869f595109d0462","0x97ff99ad8a3ae45e933464d09b485b7e1adf2fae15ea88d4215cd676b9ca959e","0x343b25b3c1140eb6bf24dbb7ef8595d62178e9ed686fb5d7e6431840c1194314","0x15eb2874404febc7c5cf63bc8ee8100d3f66bf32b69c66805f2fd24732cee39d","0xdfa64978248b67cd5941fe29fc4297ea311aca517ad0e43d71ca59b760fa9ede","0x63f77993f0db424f3bfc202d6f2d3a4cc33979588ef156deff28987c352d44bc"],"transactionsRoot": "0xcb9ecdf5483a1435113250201f690124501cfb0c071b697fcfee88c9a368ef35","uncles": []}`
var TESTRLPENCODING = `f9025fa0aa912ad61a8aa3e2d1144e4c76b746720e41682122a8b77eff890099a0ff6284a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347940000000000000000000000000000000000000000a0b347dd25d9a8a456448aed25e072c9db54f464be5e3ce1f505cc171cacf3a967a0cb9ecdf5483a1435113250201f690124501cfb0c071b697fcfee88c9a368ef35a01d000ef3f5ca9ebc62cc8aaa07e8fbd103583d1e3cbd28c13e62bc8eac5eb2f1b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002821546871405830e4c2a6c83042af2845b463454b861d88301080b846765746888676f312e31302e32856c696e757800000000000000dd2ba07230e2186ee83ef77d88298c068205167718d48ba5b6ba1de552d0c6ce156011a58b49ed91855de154346968a7eeaaf20914022e58e4f6c0e1e02567ec00a00000000000000000000000000000000000000000000000000000000000000000880000000000000000`
var TEST_NODE_VALUE = `0xf86982093f85174876e80083015f909407340652d03d131cd5737aac4a88623682e7e4c40180820bf9a070d26860a32ef4d08d6d91afa73c067af3211dd692a372770927dc9cbddd7869a05aac135e61c984c356509fc27d41b9f0c9c1f23c76d99571491bb0d15936608a`
var TEST_PATH = `0x80`
var TEST_PARENT_NODES = `0xf8c3f851a0448f4ee6a987bf17a91096e25247c3d7d78dbd08afddb5cfd4186d6a9f36bbc080808080808080a0c47289442eb85e0ca1f12c5ac6168f15513036935879931655dadfad3586dcb78080808080808080f86e30b86bf86982093f85174876e80083015f909407340652d03d131cd5737aac4a88623682e7e4c40180820bf9a070d26860a32ef4d08d6d91afa73c067af3211dd692a372770927dc9cbddd7869a05aac135e61c984c356509fc27d41b9f0c9c1f23c76d99571491bb0d15936608a`

// NOTE: These tests are skipped if go test -short is called

// Ensure that Ion is deployed as expected
func Test_IonDeployement(t *testing.T) {
	// Setup simulated block chain
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1000000000)}
	blockchain := backends.NewSimulatedBackend(alloc)

	patriciaAddress, _, _, err := DeployPatriciaTrie(
		auth,
		blockchain,
	)
	if err != nil {
		log.Fatalf("Failed to deploy PatriciaTrie library: %v", err)
	}

	// Register id of another chain
	deployedChainID, _ := utils.StringToBytes32(DEPLOYEDCHAINID)

	_, _, ion, err := LinkDeployIon(
		auth,
		blockchain,
		deployedChainID,
		patriciaAddress,
		"__./contracts/libraries/PatriciaTrie.s__",
	)
	if err != nil {
		log.Fatalf("Failed to link and deploy Ion: %v", err)
	}

	// commit all pending transactions
	blockchain.Commit()

	chainID, err := ion.ChainId(nil)
	if err != nil {
		log.Fatalf("Failed to retrieve chainID: %v", err)
	}

	// Transform into string
	CHAINID := fmt.Sprintf("%x", chainID)
	assert.Equal(t, DEPLOYEDCHAINID, CHAINID)
}

// Ensure chains are registered correctly
func Test_RegisterChain(t *testing.T) {
	// Setup simulated block chain
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1000000000)}
	blockchain := backends.NewSimulatedBackend(alloc)

	patriciaAddress, _, _, err := DeployPatriciaTrie(
		auth,
		blockchain,
	)
	if err != nil {
		log.Fatalf("Failed to deploy PatriciaTrie library: %v", err)
	}

	// Register id of another chain
	deployedChainID, _ := utils.StringToBytes32(DEPLOYEDCHAINID)

	_, _, ion, err := LinkDeployIon(
		auth,
		blockchain,
		deployedChainID,
		patriciaAddress,
		"__./contracts/libraries/PatriciaTrie.s__",
	)
	if err != nil {
		log.Fatalf("Failed to link and deploy Ion: %v", err)
	}

	// commit all pending transactions
	blockchain.Commit()

	testChainID, _ := utils.StringToBytes32(TESTCHAINID)

	// Register an alternate chain
	_, err = ion.RegisterChain(auth, testChainID)
	if err != nil {
		log.Fatalf("Failed to register chain: %v", err)
	}

	// commit all pending transactions
	blockchain.Commit()

	// Find deployed chainId
	chain, err := ion.Chains(nil, big.NewInt(0))
	if err != nil {
		log.Fatalf("Failed to retrieve chainID: %v", err)
	}

	// Transform into string
	CHAIN := fmt.Sprintf("%x", chain)
	assert.Equal(t, TESTCHAINID, CHAIN)

}

// Fail if chain is registered more than once
func Test_FailRegisterChain(t *testing.T) {
	// Setup simulated block chain
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1000000000)}
	blockchain := backends.NewSimulatedBackend(alloc)

	patriciaAddress, _, _, err := DeployPatriciaTrie(
		auth,
		blockchain,
	)
	if err != nil {
		log.Fatalf("Failed to deploy PatriciaTrie library: %v", err)
	}

	// Register id of another chain
	deployedChainID, _ := utils.StringToBytes32(DEPLOYEDCHAINID)

	_, _, ion, err := LinkDeployIon(
		auth,
		blockchain,
		deployedChainID,
		patriciaAddress,
		"__./contracts/libraries/PatriciaTrie.s__",
	)
	if err != nil {
		log.Fatalf("Failed to link and deploy Ion: %v", err)
	}

	// commit all pending transactions
	blockchain.Commit()

	testChainID, _ := utils.StringToBytes32(TESTCHAINID)

	// Register an alternate chain
	_, err = ion.RegisterChain(auth, testChainID)
	if err != nil {
		log.Fatalf("Failed to register chain: %v", err)
	}

	// commit all pending transactions
	blockchain.Commit()

	// Register the same chain again
	_, err = ion.RegisterChain(auth, testChainID)
	assert.NotEqual(t, nil, err)
}

// Ensure chains are registered correctly
func Test_SubmitBlock(t *testing.T) {
	// Setup simulated block chain
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1000000000000000000)}
	auth.GasPrice = big.NewInt(1)
	auth.GasLimit = uint64(1000000)
	blockchain := backends.NewSimulatedBackend(alloc)

	patriciaAddress, _, _, err := DeployPatriciaTrie(
		auth,
		blockchain,
	)
	if err != nil {
		log.Fatalf("Failed to deploy PatriciaTrie library: %v", err)
	}

	// Register id of another chain
	deployedChainID, _ := utils.StringToBytes32(DEPLOYEDCHAINID)

	_, _, ion, err := LinkDeployIon(
		auth,
		blockchain,
		deployedChainID,
		patriciaAddress,
		"__./contracts/libraries/PatriciaTrie.s__",
	)
	if err != nil {
		log.Fatalf("Failed to link and deploy Ion: %v", err)
	}

	// commit all pending transactions
	blockchain.Commit()

	testChainID, _ := utils.StringToBytes32(TESTCHAINID)

	// Register an alternate chain
	_, err = ion.RegisterChain(auth, testChainID)
	if err != nil {
		log.Fatalf("Failed to register chain: %v", err)
	}

	// commit all pending transactions
	blockchain.Commit()

	// Submit block
	var blockHeader utils.Header
	err = json.Unmarshal([]byte(TESTBLOCK), &blockHeader)
	if err != nil {
		log.Fatal("Unmarshal failed", err)
	}
	blockHash, _ := utils.StringToBytes32(blockHeader.Root)
	// blockParentHash, _ := utils.StringToBytes32(blockHeader.ParentHash)
	// blockTxHash, _ := utils.StringToBytes32(blockHeader.TxHash)
	// blockReceiptHash, _ := utils.StringToBytes32(blockHeader.ReceiptHash)

	rlpEncodingBytes, err := hex.DecodeString(TESTRLPENCODING)

	_, err = ion.SubmitBlock(auth, testChainID, blockHash, rlpEncodingBytes)
	if err != nil {
		log.Fatalf("Failed to submit block: %v", err)
	}

}
