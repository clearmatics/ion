// Copyright (c) 2018 Clearmatics Technologies Ltd

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Settings
type Setup struct {
	PortTo       string `json:"rpc-port-to"`
	AddrTo       string `json:"rpc-addr-to"`
	AccountTo    string `json:"account-to"`
	KeystoreTo   string `json:"keystore-to"`
	PortFrom     string `json:"rpc-port-from"`
	AddrFrom     string `json:"rpc-addr-from"`
	AccountFrom  string `json:"account-from"`
	KeystoreFrom string `json:"keystore-from"`
	Ion          string `json:"ion-addr"`
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
