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

		Ion := config.InitIon(setup, clientTo)
		Validation := config.InitValidation(setup, clientTo)

		printInfo(setup)

		// Launch the CLI
		cli.Launch(setup, clientFrom, Validation, Ion)
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
	fmt.Println("Listening on:\t\t" + setup.AddrTo)
	fmt.Println("User Account:\t\t" + setup.AccountTo)
	fmt.Println("Validation Contract:\t" + setup.Validation)
	fmt.Println("Validation ChainId:\t" + setup.ChainId)
	fmt.Println("Ion Contract:\t\t" + setup.Ion)
	fmt.Println("\nRPC Client [from]:")
	fmt.Println("Listening on:\t\t" + setup.AddrFrom)
	fmt.Println("User Account:\t\t" + setup.AccountFrom)
	fmt.Println("===============================================================")
}
