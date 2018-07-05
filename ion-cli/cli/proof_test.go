// Copyright (c) 2018 Clearmatics Technologies Ltd

package cli_test

import (
	"encoding/hex"
	"log"
	"testing"

	"github.com/clearmatics/ion/ion-cli/cli"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

func initTestClient() (client *ethclient.Client) {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func Test_GenerateTxProof(t *testing.T) {
	client := initTestClient()

	var transaction = "0xd828cadfcc7694d314058404506fc10a4dadac72aba68c286b34137da4156630"

	var expectedRootHash = common.HexToHash("0x363d32a031aea0b43d92989f96fa0d19c75f18d93fe86f1e11be79202a614bd0")
	expectedLeaf, _ := hex.DecodeString("f8cb8201e1850d09dc3000827c9f9466323324b77d72c65ea76caa918464836498ebd680b864b31d61b000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001000000000000000000000000561dedea8f2554241222e2f2eb221e625c7fb7b526a03b8772e36e5f5a4ffc4c1e6ef884c72208e6086b5d637591fe2dc7b8559aec9ba02c1da7488bdf45f01f511358b7b79e330b5e9e28bb3a10e6006e112148460801")
	expectedIdx, _ := hex.DecodeString("31") // 49 is 31 in hex bytes
	var blockNumber = "5904064"

	rootHash, idx, leaf, _ := cli.GenerateTxProof(client, transaction, blockNumber)
	assert.Equal(t, expectedRootHash, rootHash)
	assert.Equal(t, expectedIdx, idx)
	assert.Equal(t, expectedLeaf, leaf)
}
