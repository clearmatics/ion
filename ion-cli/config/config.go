// Copyright (c) 2018 Clearmatics Technologies Ltd

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Settings
type Setup struct {
	PortTo      string `json:"rpc-port-to"`
	AddrTo      string `json:"rpc-addr-to"`
	AccountTo   string `json:"account-to"`
	PortFrom    string `json:"rpc-port-from"`
	AddrFrom    string `json:"rpc-addr-from"`
	AccountFrom string `json:"account-from"`
	Ion         string `json:"ion-addr"`
}

func Read(config string) (setup Setup) {
	raw, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Print(err, "\n")
	}

	err = json.Unmarshal(raw, &setup)

	return setup
}
