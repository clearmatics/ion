// Copyright (c) 2018 Clearmatics Technologies Ltd

package cli_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ion/ion-cli/cli"
	"github.com/stretchr/testify/assert"
)

func Test_EncodePrefix(t *testing.T) {
	prefixString := "0214"
	expectedPrefix, _ := hex.DecodeString(prefixString)

	// read a fake block
	raw, err := ioutil.ReadFile("../block.json")
	if err != nil {
		fmt.Println("cannot find test block.json file:", err)
		return
	}

	// Marshall fake block into the Header
	var blockHeader cli.Header
	json.Unmarshal(raw, &blockHeader)
	prefix := cli.EncodePrefix(blockHeader)

	assert.Equal(t, expectedPrefix, prefix)

}

func Test_EncodeExtraDataPrefix(t *testing.T) {
	prefixString := "a0"
	expectedPrefix, _ := hex.DecodeString(prefixString)

	// read a fake block
	raw, err := ioutil.ReadFile("../block.json")
	if err != nil {
		fmt.Println("cannot find test block.json file:", err)
		return
	}

	// Marshall fake block into the Header
	var blockHeader cli.Header
	json.Unmarshal(raw, &blockHeader)
	prefix := cli.EncodeExtraData(blockHeader)

	assert.Equal(t, expectedPrefix, prefix)

}
