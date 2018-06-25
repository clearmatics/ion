// Copyright (c) 2018 Clearmatics Technologies Ltd

package main

import (
	"flag"
	"fmt"
	"os"

	"./cli"
	"./config"
)

var configFile = flag.String("config", "setup.json", "Description:\n path to the configuration file")

func main() {
	flag.Parse()

	if *configFile != "" {
		setup := config.Read(*configFile)

		// Launch the CLI
		cli.Launch(setup)
	} else {
		fmt.Print("Error: empty config!\n")
		os.Exit(3)
	}

}
