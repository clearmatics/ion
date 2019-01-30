// Copyright (c) 2018 Clearmatics Technologies Ltd
package cli

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/clearmatics/ion/ion-cli/utils"
	"github.com/ethereum/go-ethereum/rpc"
)

// Header used to marshall blocks into a string based struct
type header struct {
	ParentHash  string `json:"parentHash"`
	UncleHash   string `json:"sha3Uncles"`
	Coinbase    string `json:"miner"`
	Root        string `json:"stateRoot"`
	TxHash      string `json:"transactionsRoot"`
	ReceiptHash string `json:"receiptsRoot"`
	Bloom       string `json:"logsBloom"`
	Difficulty  string `json:"difficulty"`
	Number      string `json:"number"`
	GasLimit    string `json:"gasLimit"`
	GasUsed     string `json:"gasUsed"`
	Time        string `json:"timestamp"`
	Extra       string `json:"extraData"`
	MixDigest   string `json:"mixHash"`
	Nonce       string `json:"nonce"`
}

type EthClient struct {
    client *ethclient.Client
    rpcClient *rpc.Client
    url string
}

func latestBlock(eth *EthClient) (lastBlock *types.Header) {
	// var lastBlock Block
	lastBlock, err := eth.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		fmt.Println("can't get latest block:", err)
		return nil
	}

	return
}

func getBlockByNumber(eth *EthClient, number string) (*types.Header, []byte, error) {
	// var blockHeader header
	blockNum := new(big.Int)
	blockNum.SetString(number, 10)

	block, err := eth.client.HeaderByNumber(context.Background(), blockNum)
	if err != nil {
		return nil, nil, err
	}
	// Marshal into a JSON
	b, err := json.MarshalIndent(block, "", " ")
	if err != nil {
		return nil, nil, err
	}
	return block, b, nil
}

func getBlockByHash(eth *EthClient, hash string) (*types.Header, []byte, error) {
	blockHash := common.HexToHash(hash)

	block, err := eth.client.HeaderByHash(context.Background(), blockHash)
	if err != nil {
		return nil, nil, err
	}
	// Marshal into a JSON
	b, err := json.MarshalIndent(block, "", " ")
	if err != nil {
		return nil, nil, err
	}
	return block, b, nil
}

func getTransactionByHash(eth *EthClient, hash string) (*types.Transaction, []byte, error) {
	txHash := common.HexToHash(hash)

	tx, _, err := eth.client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, nil, err
	}
	// Marshal into a JSON
	t, err := json.MarshalIndent(tx, "", " ")
	if err != nil {
		return nil, nil, err
	}
	return tx, t, nil
}

func getProof(eth *EthClient, transactionHash string) {
    // Get the transaction hash
    bytesTxHash := common.HexToHash(transactionHash)

    // Generate the proof
    proof, err := utils.GenerateProof(
        context.Background(),
        eth.rpcClient,
        bytesTxHash,
    )

    if err != nil {
        panic(err)
    }

    //fmt.Printf( "Path:           0x%x\n" +
    //            "TxValue:        0x%x\n" +
    //            "TxNodes:        0x%x\n" +
    //            "ReceiptValue:   0x%x\n" +
    //            "ReceiptNodes:   0x%x\n", txPath, txValue, txNodes, receiptValue, receiptNodes)

    fmt.Printf("Proof: 0x%x\n", proof)
}

func RlpEncode(blockHeader *types.Header) (rlpSignedBlock []byte, rlpUnsignedBlock []byte) {
	// Encode the orginal block header
	_, err := rlp.EncodeToBytes(&blockHeader)
	if err != nil {
		fmt.Println("can't RLP encode requested block:", err)
		return
	}

	// Generate an interface to encode the blockheader without the signature in the extraData
	rlpSignedBlock = encodeSignedBlock(blockHeader)
	rlpUnsignedBlock = encodeUnsignedBlock(blockHeader)

	return rlpSignedBlock, rlpUnsignedBlock

}

// EncodePrefix calculate prefix of the entire signed block
func encodeUnsignedBlock(lastBlock *types.Header) (encodedBlock []byte) {
	lastBlock.Extra = lastBlock.Extra[:len(lastBlock.Extra)-65]

	encodedBlock, err := rlp.EncodeToBytes(&lastBlock)
	if err != nil {
		fmt.Println("can't RLP encode requested block:", err)
		return
	}

	return encodedBlock

}

// EncodePrefix calculate prefix of the entire signed block
func encodeSignedBlock(lastBlock *types.Header) (encodedBlock []byte) {
	encodedBlock, err := rlp.EncodeToBytes(&lastBlock)
	if err != nil {
		fmt.Println("can't RLP encode requested block:", err)
		return
	}

	return encodedBlock
}

// GenerateInterface Creates an interface for a block
func GenerateInterface(blockHeader header) (rest interface{}) {
	blockInterface := []interface{}{}
	s := reflect.ValueOf(&blockHeader).Elem()

	// Append items into the interface
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i).String()

		// Remove the 0x prefix
		f = f[2:]

		// single character then pre-pending a 0 turns it into a byte
		if len(f)%2 != 0 {
			f = "0" + f
		}

		element, _ := hex.DecodeString(f)
		blockInterface = append(blockInterface, element)
	}

	return blockInterface
}

// Encodes a block
func encodeBlock(blockInterface interface{}) (h []byte) {
	h, _ = rlp.EncodeToBytes(blockInterface)

	return h
}
