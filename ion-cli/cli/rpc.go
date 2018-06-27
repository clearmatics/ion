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

type Header struct {
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
	// var blockHeader Header
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

// func rlpEncodeBlock(client *rpc.Client, block string) {
func rlpEncodeBlock(client *ethclient.Client, block string) {
	var blockHeader Header
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
	err = json.Unmarshal([]byte(b), &blockHeader)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	// fmt.Printf("%+v\n", blockHeader)
	blockInterface := GenerateInterface(blockHeader)
	encodedBlock := EncodeBlock(blockInterface)
	fmt.Printf("%+x\n", encodedBlock)
}

// func calculateRlpEncoding(client *ethclient.Client, block string) {
func calculateRlpEncoding(client *ethclient.Client, block string) (rlpBlock []byte, prefixBlock []byte, prefixExtra []byte) {
	var blockHeader Header
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
	blockHeader.Extra = blockHeader.Extra[:len(blockHeader.Extra)-130]
	blockInterface := GenerateInterface(blockHeader)
	encodedPrefixBlock := EncodeBlock(blockInterface)
	prefixBlock = encodedPrefixBlock[1:3]
	fmt.Printf("\nSigned Block Header Prefix:\n%+x\n", prefixBlock)

	// Generate an interface to encode the blockheader without the signature in the extraData
	encExtra, _ := hex.DecodeString(blockHeader.Extra[2:])
	encodedExtraData := EncodeBlock(encExtra)
	prefixExtra = encodedExtraData[0:1]
	fmt.Printf("\nExtraData Field Prefix:\n%+x\n", prefixExtra)

	return rlpBlock, prefixBlock, prefixExtra

}

// Creates an interface for a block
func GenerateInterface(blockHeader Header) (rest interface{}) {
	blockInterface := []interface{}{}
	s := reflect.ValueOf(&blockHeader).Elem()
	// fmt.Println(s)

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
		// fmt.Printf("%+v\n", s.Field(i))
		// fmt.Printf("%+x\n", f)
		// fmt.Printf("%+x\n", element)
		blockInterface = append(blockInterface, element)
	}

	return blockInterface
}

// Encodes a block
func EncodeBlock(blockInterface interface{}) (h []byte) {
	h, _ = rlp.EncodeToBytes(blockInterface)

	return h
}
