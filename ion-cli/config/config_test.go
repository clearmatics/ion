// Copyright (c) 2018 Clearmatics Technologies Ltd

package config_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"./config"
)

func Test_Read_ValidSetupJson(t *testing.T) {
	path := findPath() + "../setup.json"
	setup := config.Read(path)

	assert.Equal(t, "8501", setup.PortTo)
	assert.Equal(t, "127.0.0.1", setup.AddrTo)
	assert.Equal(t, "8502", setup.PortFrom)
	assert.Equal(t, "127.0.0.1", setup.AddrFrom)
	assert.Equal(t, "0xb9fd43a71c076f02d1dbbf473c389f0eacec559f", setup.Ion)
}

func findPath() string {
	_, path, _, _ := runtime.Caller(0)
	pathSlice := strings.Split(path, "/")
	return strings.Trim(path, pathSlice[len(pathSlice)-1])
}
