// Copyright (c) 2018 Clearmatics Technologies Ltd

package main

import (
    "fmt"
	"github.com/clearmatics/ion/ion-cli/cli"
)

func main() {
    // Launch the CLI
    printWelcome()
    cli.Launch()
}

func printWelcome() {
	// display welcome info.
	fmt.Println("===============================================================")
	fmt.Print("Ion Command Line Interface\n\n")
	fmt.Println("Use 'help' to list commands")
	fmt.Println("===============================================================")
}
