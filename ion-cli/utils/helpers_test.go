// Copyright (c) 2018 Clearmatics Technologies Ltd

package utils_test

import (
	"fmt"
	"testing"

	"github.com/clearmatics/ion/ion-cli/utils"
	"github.com/stretchr/testify/assert"
)

var TESTSTRING = "aa912ad61a8aa3e2d1144e4c76b746720e41682122a8b77eff890099a0ff6284"
var TESTSTRINGPREFIX = "0xaa912ad61a8aa3e2d1144e4c76b746720e41682122a8b77eff890099a0ff6284"
var TESTSTRINGFAIL = "0xaa912ad61a8aa3e2d1144e4c76b746720e41682122a8b77eff890099a0ff628432"

// Encodes a hex string as bytes array specifically Bytes32
func Test_StringToBytes32(t *testing.T) {
	result, _ := utils.StringToBytes32(TESTSTRING)

	// Transform bytes back into string
	stringResult := fmt.Sprintf("%x", result)
	assert.Equal(t, TESTSTRING, stringResult)
}

// Encodes a hex string with prefix as bytes array specifically Bytes32
func Test_StringToBytes32_Prefix(t *testing.T) {
	result, _ := utils.StringToBytes32(TESTSTRINGPREFIX[2:])

	// Transform bytes back into string
	stringResult := fmt.Sprintf("%x", result)
	assert.Equal(t, TESTSTRING, stringResult)
}

// Encodes a hex string with prefix as bytes array specifically Bytes32
func Test_StringToBytes32_IncorrectInput(t *testing.T) {
	_, err := utils.StringToBytes32(TESTSTRINGFAIL)

	assert.NotEqual(t, nil, err)
}
