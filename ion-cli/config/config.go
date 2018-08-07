// Copyright (c) 2018 Clearmatics Technologies Ltd

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// Settings
type Setup struct {
	AddrTo       string `json:"rpc-to"`
	AccountTo    string `json:"account-to"`
	PasswordTo   string `json:"password-to"`
	KeystoreTo   string `json:"keystore-to"`
	AddrFrom     string `json:"rpc-from"`
	AccountFrom  string `json:"account-from"`
	PasswordFrom string `json:"password-from"`
	KeystoreFrom string `json:"keystore-from"`
	ChainId      string `json:"validation-chainid"`
	Validation   string `json:"validation-addr"`
	Ion          string `json:"ion-addr"`
	Trigger      string `json:"trigger-addr"`
}

// Takes path to a JSON and returns a struct of the contents
func ReadSetup(config string) (setup Setup) {
	raw, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Print(err, "\n")
	}

	err = json.Unmarshal(raw, &setup)

	return
}

// Takes path to a JSON and returns a string of the contents
func ReadString(path string) (contents string) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err, "\n")
	}

	contents = string(raw)

	return

}

func InitUser(privkeystore string, password string) (auth *bind.TransactOpts, userkey *keystore.Key) {
	// retrieve private key
	keyjson, err := ioutil.ReadFile(privkeystore)
	if err != nil {
		fmt.Println("Error failed to read keystore: %v", err)
	}

	userkey, err = keystore.DecryptKey(keyjson, password)
	if err != nil {
		fmt.Println("Error json key failed to decrypt: %v", err)
	}

	// Create an authorized transactor
	key := ReadString(privkeystore)
	auth, err = bind.NewTransactor(strings.NewReader(key), password)
	if err != nil {
		log.Fatalf("Error failed to create authorized transactor: %v", err)
	}

	return
}
