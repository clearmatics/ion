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

func GenerateProof(ctx context.Context, client *rpc.Client, txHash common.Hash) (txTriggerPath []byte, txTriggerRLP []byte, txTriggerProofArr []byte, receiptTrigger []byte, receiptTriggerProofArr []byte) {
	blockNumberStr, txTrigger, err := BlockNumberByTransactionHash(ctx, client, txHash)
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
	tx := block.Transactions()
	txTrie := TxTrie(tx)
	blockReceipts := GetBlockTxReceipts(clientETH, block)
	receiptTrie := ReceiptTrie(blockReceipts)

	// Calculate transaction index)
	for i := 0; i < len(tx); i++ {
		if txHash == tx[i].Hash() {
			idx = byte(i)
		}
	}

	txTriggerPath = append(txTriggerPath, idx)
	txTriggerRLP, _ = rlp.EncodeToBytes(txTrigger)
	txTriggerProofArr = Proof(txTrie, txTriggerPath[:])
	receiptTrigger, _ = rlp.EncodeToBytes(blockReceipts[txTriggerPath[0]])
	receiptTriggerProofArr = Proof(receiptTrie, txTriggerPath[:])

	return
}
