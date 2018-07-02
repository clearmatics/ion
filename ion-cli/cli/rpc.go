// Copyright (c) 2018 Clearmatics Technologies Ltd

package cli

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
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

func latestBlock(client *ethclient.Client) {
	// var lastBlock Block
	lastBlock, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		fmt.Println("can't get latest block:", err)
		return
	}
	// Print events from the subscription as they arrive.
	fmt.Printf("latest block: %v\n", lastBlock.Number)
}

func getBlock(client *ethclient.Client, block string) {
	// var blockHeader header
	blockNum := new(big.Int)
	blockNum.SetString(block, 10)

	lastBlock, err := client.HeaderByNumber(context.Background(), blockNum)
	if err != nil {
		fmt.Println("can't get requested block:", err)
		return
	}
	// Marshal into a JSON
	b, err := json.MarshalIndent(lastBlock, "", " ")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println("Block:", block)
	fmt.Println(string(b))
}

// func calculateRlpEncoding(client *ethclient.Client, block string) {
func calculateRlpEncoding(client *ethclient.Client, block string) (rlpBlock []byte, prefixBlock []byte, prefixExtra []byte) {
	var blockHeader header
	blockNum := new(big.Int)
	blockNum.SetString(block, 10)

	lastBlock, err := client.HeaderByNumber(context.Background(), blockNum)
	if err != nil {
		fmt.Println("can't get requested block:", err)
		return
	}

	// Encode the orginal block header
	rlpBlock, err = rlp.EncodeToBytes(&lastBlock)
	if err != nil {
		fmt.Println("can't RLP encode requested block:", err)
		return
	}
	fmt.Printf("\nEncoded Block Header:\n%+x\n", rlpBlock)

	// Marshal block into a JSON to allow manipulation of specific fields
	b, err := json.MarshalIndent(lastBlock, "", " ")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	err = json.Unmarshal([]byte(b), &blockHeader)
	if err != nil {
		latestBlock(client)
		fmt.Printf("Error: %s", err)
		return
	}

	// Generate an interface to encode the blockheader without the signature in the extraData
	prefixBlock = encodeBlock(blockHeader)
	fmt.Printf("\nSigned Block Header Prefix:\n%+x\n", prefixBlock)

	// Generate an interface to encode the blockheader without the signature in the extraData
	prefixExtra = encodeExtraData(blockHeader)
	fmt.Printf("\nExtraData Field Prefix:\n%+x\n", prefixExtra)

	return rlpBlock, prefixBlock, prefixExtra

}

// EncodePrefix calculate prefix of the entire signed block
func encodePrefix(blockHeader header) (prefix []byte) {
	blockHeader.Extra = blockHeader.Extra[:len(blockHeader.Extra)-130]
	blockInterface := GenerateInterface(blockHeader)
	encodedPrefixBlock := encodeBlock(blockInterface)

	return encodedPrefixBlock[1:3]
}

// EncodeExtraData calculate prefix of the extraData with the signature
func encodeExtraData(blockHeader header) (prefix []byte) {
	blockHeader.Extra = blockHeader.Extra[:len(blockHeader.Extra)-130]
	encExtra, err := hex.DecodeString(blockHeader.Extra[2:])
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	encodedExtraData := encodeBlock(encExtra)

	return encodedExtraData[0:1]
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
