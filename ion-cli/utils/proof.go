// Copyright (c) 2018 Clearmatics Technologies Ltd
package utils

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

func GenerateProof(ctx context.Context, client *rpc.Client, txHash common.Hash) ([]byte, error) {
	blockNumberStr, tx, err := BlockNumberByTransactionHash(ctx, client, txHash)
	if err != nil {
		fmt.Printf("Error: couldn't find block by tx hash: %s\n", err)
	}

	// Convert returned blocknumber
	var blockNumber big.Int
	blockNumber.SetString((*blockNumberStr)[2:], 16)

	clientETH := ethclient.NewClient(client)
	eventTxBlockNumber := blockNumber
	block, err := clientETH.BlockByNumber(ctx, &eventTxBlockNumber)
	if err != nil {
		fmt.Printf("Error: retrieving block: %s\n", err)
	}

	var idx byte
	txs := block.Transactions()
	txTrie := TxTrie(txs)
	blockReceipts := GetBlockTxReceipts(clientETH, block)
	receiptTrie := ReceiptTrie(blockReceipts)

	// Calculate transaction index)
	for i := 0; i < len(txs); i++ {
		if txHash == txs[i].Hash() {
			idx = byte(i)
		}
	}


	txPath := []byte{idx}
	txRLP, _ := rlp.EncodeToBytes(tx)
	txProof := Proof(txTrie, txPath[:])
	receiptRLP, _ := rlp.EncodeToBytes(blockReceipts[txPath[0]])
	receiptProof := Proof(receiptTrie, txPath[:])

	var decodedTx, decodedTxProof, decodedReceipt, decodedReceiptProof []interface{}

	rlp.DecodeBytes(txRLP, &decodedTx)
	rlp.DecodeBytes(txProof, &decodedTxProof)
	rlp.DecodeBytes(receiptRLP, &decodedReceipt)
	rlp.DecodeBytes(receiptProof, &decodedReceiptProof)

    proof := make([]interface{}, 0)
    proof = append(proof, txPath, decodedTx, decodedTxProof, decodedReceipt, decodedReceiptProof)

	return rlp.EncodeToBytes(proof)
}
