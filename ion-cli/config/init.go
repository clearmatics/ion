// Copyright (c) 2018 Clearmatics Technologies Ltd

package config

import (
	"log"

	contract "github.com/clearmatics/ion/ion-cli/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func InitClient(addr string) (clientTo *ethclient.Client) {
	// Connect to the RPC Client
	clientTo, err := ethclient.Dial(addr)
	if err != nil {
		log.Fatalf("could not create RPC client: %v", err)
	}

	return
}

func InitValidationContract(setup Setup, client *ethclient.Client) (Ion *contract.Ion) {
	// Initialise the contract
	address := common.HexToAddress(setup.Ion)
	Ion, err := contract.NewIon(address, client)
	if err != nil {
		log.Fatal(err)
	}

	return
}
