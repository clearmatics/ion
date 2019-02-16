// Copyright (c) 2018 Clearmatics Technologies Ltd

package utils

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/clearmatics/autonity/accounts/abi/bind"
	"github.com/clearmatics/autonity/ethclient"
)

func GetNonce(client *ethclient.Client, auth *bind.TransactOpts) {
	// Find the correct tx nonce
	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		log.Fatalf("Failed to calculate nonce: %v", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
}

// Needed to convert transaction strings into the correct format
func StringToBytes32(input string) (output [32]byte, err error) {
	// Check string length is correct 64
	if len(input) == 64 {
		inputBytes, err := hex.DecodeString(input)
		if err != nil {
			log.Fatalf("Failed to encode string as bytes: %v", err)
		}

		copy(output[:], inputBytes[:len(output)])

		return output, nil
	} else if len(input) == 66 && input[:2] == "0x" {
		inputBytes, err := hex.DecodeString(input[2:])
		if err != nil {
			log.Fatalf("Failed to encode string as bytes: %v", err)
		}

		copy(output[:], inputBytes[:len(output)])

		return output, nil
	} else {
		return [32]byte{}, fmt.Errorf("Failed to encode string as bytes32, incorrect string input")
	}

}
