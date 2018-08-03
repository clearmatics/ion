// Copyright (c) 2018 Clearmatics Technologies Ltd

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Settings
type Setup struct {
	AddrTo       string `json:"rpc-to"`
	AccountTo    string `json:"account-to"`
	KeystoreTo   string `json:"keystore-to"`
	AddrFrom     string `json:"rpc-from"`
	AccountFrom  string `json:"account-from"`
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
