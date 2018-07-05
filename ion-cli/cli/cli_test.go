// Copyright (c) 2018 Clearmatics Technologies Ltd

package cli_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/clearmatics/ion/ion-cli/config"
)

func Test_Read_ValidSetupJson(t *testing.T) {
	path := findPath() + "./test.json"
	setup := config.ReadSetup(path)

	assert.Equal(t, "127.0.0.1:8545", setup.AddrTo)
	assert.Equal(t, "0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d", setup.AccountTo)
	assert.Equal(t, "../poa-network/node1/keystore/UTC--2018-06-05T09-31-57.109288703Z--2be5ab0e43b6dc2908d5321cf318f35b80d0c10d", setup.KeystoreTo)
	assert.Equal(t, "0xb9fd43a71c076f02d1dbbf473c389f0eacec559f", setup.Ion)
	assert.Equal(t, "127.0.0.1:8501", setup.AddrFrom)
	assert.Equal(t, "0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d", setup.AccountFrom)
	assert.Equal(t, "../poa-network/node1/keystore/UTC--2018-06-05T09-31-57.109288703Z--2be5ab0e43b6dc2908d5321cf318f35b80d0c10d", setup.KeystoreFrom)
}

func findPath() string {
	_, path, _, _ := runtime.Caller(0)
	pathSlice := strings.Split(path, "/")
	return strings.Trim(path, pathSlice[len(pathSlice)-1])
}
