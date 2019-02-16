// Copyright (c) 2018 Clearmatics Technologies Ltd
package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	ethereum "github.com/clearmatics/autonity"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/ethclient"
	"github.com/clearmatics/autonity/rpc"
)

// Client gets client or fails if no connection
func Client(url string) *ethclient.Client {
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal("Client failed to connect: ", err)
	} else {
		fmt.Println("Connected to: ", url)
	}
	return client
}

// GetBlockTxReceipts get the receipts for all the transactions in a block
func GetBlockTxReceipts(ec *ethclient.Client, block *types.Block) []*types.Receipt {
	var receiptsArr []*types.Receipt
	for _, tx := range block.Transactions() {
		receipt, err := ec.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal("TransactionReceipt ERROR:", err)
		}
		receiptsArr = append(receiptsArr, receipt)
	}
	return receiptsArr
}

// -------
// Since you can't get a block by giving it the transaction hash in go-ethereum
// the only solution was to replicate their code and add that feature to it
// really annoying!!!
// -------

// ClientRPC RPC Client gets an RPC client (useful to get the block number out of a transaction)
func ClientRPC(url string) *rpc.Client {
	c, err := rpc.DialContext(context.Background(), url)
	if err != nil {
		log.Fatal("RPC Client failed to connect: ", err)
	}
	return c
}

type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}

// BlockNumberByTransactionHash gets a block number by a transaction hash in that block
func BlockNumberByTransactionHash(ctx context.Context, c *rpc.Client, txHash common.Hash) (*string, *types.Transaction, error) {
	var json *rpcTransaction
	var err error
	err = c.CallContext(ctx, &json, "eth_getTransactionByHash", txHash)
	if err != nil {
		return nil, nil, err
	} else if json == nil {
		return nil, nil, ethereum.NotFound
	} else if _, r, _ := json.tx.RawSignatureValues(); r == nil {
		return nil, nil, fmt.Errorf("server returned transaction without signature")
	}
	return json.BlockNumber, json.tx, nil
}
