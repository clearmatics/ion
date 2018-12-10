// Copyright (c) 2018 Clearmatics Technologies Ltd

package utils

import (
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/rlp"
)

// Header used to marshall blocks into a string based struct
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

// EncodePrefix calculate prefix of the entire signed block
func EncodePrefix(blockHeader Header) (prefix []byte) {
	blockHeader.Extra = blockHeader.Extra[:len(blockHeader.Extra)-130]
	blockInterface := GenerateInterface(blockHeader)
	encodedPrefixBlock := EncodeBlock(blockInterface)

	return encodedPrefixBlock[1:3]
}

// EncodeExtraData calculate prefix of the extraData with the signature
func EncodeExtraData(blockHeader Header) (prefix []byte) {
	blockHeader.Extra = blockHeader.Extra[:len(blockHeader.Extra)-130]
	encExtra, err := hex.DecodeString(blockHeader.Extra[2:])
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	encodedExtraData := EncodeBlock(encExtra)

	return encodedExtraData[0:1]
}

// GenerateInterface Creates an interface for a block
func GenerateInterface(blockHeader Header) (rest interface{}) {
	blockInterface := []interface{}{}
	s := reflect.ValueOf(&blockHeader).Elem()

	// Append items into the interface
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i).String()
		// fmt.Printf("\n%s", f)

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
func EncodeBlock(blockInterface interface{}) (h []byte) {
	h, _ = rlp.EncodeToBytes(blockInterface)

	return h
}
