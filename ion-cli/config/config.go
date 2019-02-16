// Copyright (c) 2018 Clearmatics Technologies Ltd

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/clearmatics/autonity/accounts/abi/bind"
	"github.com/clearmatics/autonity/accounts/keystore"
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
	Function     string `json:"function-addr"`
}

type Account struct {
	Auth *bind.TransactOpts
	Key  *keystore.Key
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

func InitUser(privkeystore string, password string) (auth *bind.TransactOpts, userkey *keystore.Key, err error) {
	// retrieve private key
	keyjson, err := ioutil.ReadFile(privkeystore)
	if err != nil {
		return nil, nil, err
	}

	userkey, err = keystore.DecryptKey(keyjson, password)
	if err != nil {
		return nil, nil, err
	}

	// Create an authorized transactor
	key := ReadString(privkeystore)
	auth, err = bind.NewTransactor(strings.NewReader(key), password)
	if err != nil {
		return nil, nil, err
	}

	return
}
