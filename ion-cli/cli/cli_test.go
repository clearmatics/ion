// Copyright (c) 2018 Clearmatics Technologies Ltd

package cli_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"../config"
)

func Test_Read_ValidSetupJson(t *testing.T) {
	path := findPath() + "../setup.json"
	setup := config.ReadSetup(path)

	assert.Equal(t, "8501", setup.PortTo)
	assert.Equal(t, "127.0.0.1", setup.AddrTo)
	assert.Equal(t, "8502", setup.PortFrom)
	assert.Equal(t, "127.0.0.1", setup.AddrFrom)
}

func findPath() string {
	_, path, _, _ := runtime.Caller(0)
	pathSlice := strings.Split(path, "/")
	return strings.Trim(path, pathSlice[len(pathSlice)-1])
}
