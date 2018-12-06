// Copyright (c) 2018 Clearmatics Technologies Ltd

package config_test

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"

	"github.com/clearmatics/ion/ion-cli/config"
)

func Test_ReadValidSetupJson(t *testing.T) {
	path := findPath() + "./test.json"
	setup := config.ReadSetup(path)

	assert.Equal(t, "127.0.0.1:8545", setup.AddrTo)
	assert.Equal(t, "127.0.0.1:8501", setup.AddrFrom)
	assert.Equal(t, "0xb9fd43a71c076f02d1dbbf473c389f0eacec559f", setup.Ion)
}

func Test_ReadValidKeystore(t *testing.T) {
	path := findPath() + "./UTC--2018-06-05T09-31-57.109288703Z--2be5ab0e43b6dc2908d5321cf318f35b80d0c10d"
	contents := config.ReadString(path)

	const val = "{\"address\":\"2be5ab0e43b6dc2908d5321cf318f35b80d0c10d\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"0b11aa865046778a1b16a9b8cb593df704e3fe09f153823d75442ad1aab66caa\",\"cipherparams\":{\"iv\":\"4aa66b789ee2d98cf77272a72eeeaa50\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"b957fa7b7577240fd3791168bbe08903af4c8cc62c304f1df072dc2a59b1765e\"},\"mac\":\"197a06eb0449301d871400a6bdf6c136b6f7658ee41e3f2f7fd81ca11cd954a3\"},\"id\":\"a3cc1eae-3e36-4659-b759-6cf416216e72\",\"version\":3}"

	assert.Equal(t, val, contents)

}

func Test_InitUser(t *testing.T) {
	keystore := "./UTC--2018-06-05T09-31-57.109288703Z--2be5ab0e43b6dc2908d5321cf318f35b80d0c10d"
	password := "password1"
	expectedFrom := common.HexToAddress("2be5ab0e43b6dc2908d5321cf318f35b80d0c10d")
	expectedPrivateKey := "e176c157b5ae6413726c23094bb82198eb283030409624965231606ec0fbe65b"

	auth, userkey, err := config.InitUser(keystore, password)

	assert.Equal(t, auth.From, expectedFrom)
	privateKey := fmt.Sprintf("%x", crypto.FromECDSA(userkey.PrivateKey))
	assert.Equal(t, privateKey, expectedPrivateKey)

}

func findPath() string {
	_, path, _, _ := runtime.Caller(0)
	pathSlice := strings.Split(path, "/")
	return strings.Trim(path, pathSlice[len(pathSlice)-1])
}
