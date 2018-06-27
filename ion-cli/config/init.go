// Copyright (c) 2018 Clearmatics Technologies Ltd

package config

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ion/ion-cli/Validation"
)

func InitClient(port string, addr string) (clientTo *ethclient.Client) {
	// Connect to the RPC Client
	clientTo, err := ethclient.Dial("http://" + addr + ":" + port)
	if err != nil {
		log.Fatalf("could not create RPC client: %v", err)
	}

	return
}

func InitValidationContract(setup Setup, client *ethclient.Client) (validation *Validation.Validation) {
	// Initialise the contract
	address := common.HexToAddress(setup.Ion)
	validation, err := Validation.NewValidation(address, client)
	if err != nil {
		log.Fatal(err)
	}

	return
}
