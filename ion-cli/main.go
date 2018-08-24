// Copyright (c) 2018 Clearmatics Technologies Ltd

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/clearmatics/ion/ion-cli/cli"
	"github.com/clearmatics/ion/ion-cli/config"
	contract "github.com/clearmatics/ion/ion-cli/contracts"
	"github.com/clearmatics/ion/ion-cli/utils"
)

var configFile = flag.String("config", "setup.json", "Description:\n path to the configuration file")

func main() {
	flag.Parse()

	if *configFile != "" {
		setup := config.ReadSetup(*configFile)

		clientTo := utils.ClientRPC(setup.AddrTo)
		clientFrom := utils.ClientRPC(setup.AddrFrom)

		// Compile contracts to use in sending transactions
		Validation := contract.CompileContract("Validation")
		Function := contract.CompileContract("Function")
		Trigger := contract.CompileContract("Trigger")
		printInfo(setup)

		// Launch the CLI
		cli.Launch(
			setup,
			clientTo,
			clientFrom,
			Validation,
			Trigger,
			Function,
		)

	} else {
		fmt.Print("Error: empty config!\n")
		os.Exit(3)
	}

}

func printInfo(setup config.Setup) {
	// display welcome info.
	fmt.Println("===============================================================")
	fmt.Println("Ion Command Line Interface\n")
	fmt.Println("RPC Client [TO]:")
	fmt.Println("\tListening on:\t\t" + setup.AddrTo)
	fmt.Println("\tUser Account:\t\t" + setup.AccountTo)
	fmt.Println("\tRPC ChainId:\t\t" + setup.ChainId)
	fmt.Println("\tValidation Contract:\t" + setup.Validation)
	fmt.Println("\tIon Contract:\t\t" + setup.Ion)
	fmt.Println("\tFunction Contract:\t" + setup.Function)
	fmt.Println("\nRPC Client [FROM]:")
	fmt.Println("\tListening on:\t\t" + setup.AddrFrom)
	fmt.Println("\tUser Account:\t\t" + setup.AccountFrom)
	fmt.Println("\tTrigger Contract:\t" + setup.Trigger)
	fmt.Println("===============================================================")
}
