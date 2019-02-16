// Copyright (c) 2018 Clearmatics Technologies Ltd
package utils

import (
	"log"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/rlp"
	"github.com/clearmatics/autonity/trie"
)

// ReceiptTrie generate trie for receipts
// TODO: the argument should be of type interface so that this is a generic function
func ReceiptTrie(receipts []*types.Receipt) *trie.Trie {
	var receiptRLPidxArr, receiptRLPArr [][]byte
	for idx, receipt := range receipts {
		idxRLP, err := rlp.EncodeToBytes(uint(idx))
		if err != nil {
			log.Fatalf("ReceiptTrie RLP error: %v", err)
		}
		txRLP, err := rlp.EncodeToBytes(receipt)
		if err != nil {
			log.Fatalf("ReceiptTrie RLP error: %v", err)
		}

		receiptRLPidxArr = append(receiptRLPidxArr, idxRLP)
		receiptRLPArr = append(receiptRLPArr, txRLP)
	}

	trieObj := generateTrie(receiptRLPidxArr, receiptRLPArr)

	return trieObj
}

// TxTrie generated Trie out of transaction array
// TODO: the argument should be of type interface so that this is a generic function
func TxTrie(transactions []*types.Transaction) *trie.Trie {
	var txRLPIdxArr, txRLPArr [][]byte
	for idx, tx := range transactions {
		idxRLP, err := rlp.EncodeToBytes(uint(idx))
		if err != nil {
			log.Fatalf("TxTrie RLP error: %v", err)
		}
		txRLP, err := rlp.EncodeToBytes(tx)
		if err != nil {
			log.Fatalf("TxTrie RLP error: %v", err)
		}

		txRLPIdxArr = append(txRLPIdxArr, idxRLP)
		txRLPArr = append(txRLPArr, txRLP)
	}

	trieObj := generateTrie(txRLPIdxArr, txRLPArr)

	return trieObj
}

func generateTrie(paths [][]byte, values [][]byte) *trie.Trie {
	if len(paths) != len(values) {
		log.Fatal("Paths array and Values array have different lengths when generating Trie")
	}

	trieDB := trie.NewDatabase(ethdb.NewMemDatabase())
	trieObj, _ := trie.New(common.Hash{}, trieDB) // empty trie

	for idx := range paths {
		p := paths[idx]
		v := values[idx]

		trieObj.Update(p, v) // update trie with the rlp encode index and the rlp encoded transaction
	}

	_, err := trieObj.Commit(nil) // commit to database (which in this case is stored in memory)
	if err != nil {
		log.Fatalf("commit error: %v", err)
	}

	return trieObj
}

// Proof creates an array of the proof pathj ordered
func Proof(trie *trie.Trie, path []byte) []byte {
	proof := generateProof(trie, path)
	proofRLP, err := rlp.EncodeToBytes(proof)
	if err != nil {
		log.Fatal("ERROR encoding proof: ", err)
	}
	return proofRLP
}

func generateProof(trie *trie.Trie, path []byte) []interface{} {
	proof := ethdb.NewMemDatabase()
	err := trie.Prove(path, 0, proof)
	if err != nil {
		log.Fatal("ERROR failed to create proof")
	}

	var proofArr []interface{}
	for nodeIt := trie.NodeIterator(nil); nodeIt.Next(true); {
		if val, err := proof.Get(nodeIt.Hash().Bytes()); val != nil && err == nil {
			var decodedVal interface{}
			err = rlp.DecodeBytes(val, &decodedVal)
			if err != nil {
				log.Fatalf("ERROR(%s) failed decoding RLP: 0x%0x\n", err, val)
			}
			proofArr = append(proofArr, decodedVal)
		}
	}
	return proofArr
}
