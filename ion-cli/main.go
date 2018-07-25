// Copyright (c) 2018 Clearmatics Technologies Ltd

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/clearmatics/ion/ion-cli/cli"
	"github.com/clearmatics/ion/ion-cli/config"
)

var configFile = flag.String("config", "setup.json", "Description:\n path to the configuration file")

func main() {
	flag.Parse()

	if *configFile != "" {
		setup := config.ReadSetup(*configFile)

		clientTo := config.InitClient(setup.AddrTo)
		clientFrom := config.InitClient(setup.AddrFrom)

		Validation := config.InitValidationContract(setup, clientTo)

		printInfo(setup)

		// Launch the CLI
		cli.Launch(setup, clientFrom, Validation)
	} else {
		fmt.Print("Error: empty config!\n")
		os.Exit(3)
	}

}

func printInfo(setup config.Setup) {
	// display welcome info.
	fmt.Println("===============================================================")
	fmt.Println("Ion Command Line Interface\n")
	fmt.Println("RPC Client [to]:")
	fmt.Println("Listening on: " + setup.AddrTo)
	fmt.Println("user Account: " + setup.AccountTo)
	fmt.Println("Ion Contract: " + setup.Ion)
	fmt.Println("\nRPC Client [from]: ")
	fmt.Println("Listening on: " + setup.AddrFrom)
	fmt.Println("user Account: " + setup.AccountFrom)
	fmt.Println("===============================================================")
}
