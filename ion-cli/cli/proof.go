package cli

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/maxrobot/go-ethereum/crypto"
)

// GenerateTxProof takes a transaction and block, returns the trie root, tx index and proof path
func GenerateTxProof(client *ethclient.Client, transaction string, blockNum string) (root common.Hash, txIdx []byte, leaf []byte, proof *ethdb.MemDatabase) {
	// Select a specific block
	blockNumber := new(big.Int)
	blockNumber.SetString(blockNum, 10)

	// Fetch header of block num
	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch block of block num
	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}

	// Select a transaction, index should be 49
	transx := block.Transaction(common.HexToHash(transaction))

	// Generate the trie
	trieObj := new(trie.Trie)
	for idx, tx := range block.Transactions() {
		rlpIdx, _ := rlp.EncodeToBytes(uint(idx))  // rlp encode index of transaction
		rlpTransaction, _ := rlp.EncodeToBytes(tx) // rlp encode transaction

		trieObj.Update(rlpIdx, rlpTransaction)

		// Get the information about the transaction I care about...
		if transx == tx {
			txIdx = rlpIdx
			leaf = rlpTransaction
		}

	}

	root = trieObj.Hash()

	// Generate a merkle proof for a key
	proof = ethdb.NewMemDatabase()
	trieObj.Prove(txIdx, 0, proof)
	if proof == nil {
		fmt.Printf("prover: nil proof")
	}

	return
}

// VerifyTxProof takes a transaction and block, returns the trie root, tx index and proof path
func VerifyTxProof(client *ethclient.Client, transaction string, blockNum string) {
	// Select a specific block
	blockNumber := new(big.Int)
	blockNumber.SetString(blockNum, 10)

	// Fetch header of block num
	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	// assert.Equal(t, expectedBlockHash, header.Hash())

	// Fetch block of block num
	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}
	// assert.Equal(t, expectedBlockHash, block.Hash())

	// Select a transaction, index should be 49
	transx := block.Transaction(common.HexToHash(transaction))
	var txIdx []byte
	var leaf []byte
	// fmt.Printf("\nTransaction:\n% 0x", transx.Hash().Bytes())

	// Generate the trie
	trieObj := new(trie.Trie)
	for idx, tx := range block.Transactions() {
		rlpIdx, _ := rlp.EncodeToBytes(uint(idx))  // rlp encode index of transaction
		rlpTransaction, _ := rlp.EncodeToBytes(tx) // rlp encode transaction

		trieObj.Update(rlpIdx, rlpTransaction)

		txRlpHash := crypto.Keccak256Hash(rlpTransaction)

		fmt.Printf("TxHash[%d]: \t% 0x\n\tHash(RLP(Tx)): \t% 0x\n",
			idx, tx.Hash().Bytes(), txRlpHash.Bytes())

		// Get the information about the transaction I care about...
		if transx == tx {
			txIdx = rlpIdx
			leaf = rlpTransaction
		}

	}
	root := trieObj.Hash()
	expectedRoot := block.TxHash()

	fmt.Printf("\nExpected Root:\t%x\nRecovered Root:\t%x\n", expectedRoot, root)

	// Generate a merkle proof for a key
	proof := ethdb.NewMemDatabase()
	trieObj.Prove(txIdx, 0, proof)
	if proof == nil {
		fmt.Printf("prover: nil proof")
	}

	// Verify the proof
	val, _, err := trie.VerifyProof(root, txIdx, proof)
	if err != nil {
		fmt.Printf("prover: failed to verify proof: %v\nraw proof: %x", err, proof)
	}
	if !bytes.Equal(val, leaf) {
		fmt.Printf("prover: verified value mismatch: have %x, want 'k'", val)
	}
	fmt.Printf("\nVerified Value:\t%x\nExpected Leaf:\t%x\n", val, leaf)
}
