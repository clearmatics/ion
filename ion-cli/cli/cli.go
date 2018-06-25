// Copyright (c) 2018 Clearmatics Technologies Ltd

package cli

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ion/ion-cli/Validation"
	"github.com/ion/ion-cli/config"
)

const key = `{"address":"2be5ab0e43b6dc2908d5321cf318f35b80d0c10d","crypto":{"cipher":"aes-128-ctr","ciphertext":"0b11aa865046778a1b16a9b8cb593df704e3fe09f153823d75442ad1aab66caa","cipherparams":{"iv":"4aa66b789ee2d98cf77272a72eeeaa50"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"b957fa7b7577240fd3791168bbe08903af4c8cc62c304f1df072dc2a59b1765e"},"mac":"197a06eb0449301d871400a6bdf6c136b6f7658ee41e3f2f7fd81ca11cd954a3"},"id":"a3cc1eae-3e36-4659-b759-6cf416216e72","version":3}`

// Launch - definition of commands and creates the iterface
func Launch(setup config.Setup) {
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	// Connect to the RPC Client
	client, err := ethclient.Dial("http://" + setup.AddrTo + ":" + setup.PortTo)
	if err != nil {
		log.Fatalf("could not create RPC client: %v", err)
	}

	// Initialise the contract
	address := common.HexToAddress(setup.Ion)
	validation, err := Validation.NewValidation(address, client)
	if err != nil {
		log.Fatal(err)
	}

	// Get a suggested gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Create an authorized transactor and spend 1 unicorn
	keyTo := config.ReadString(setup.KeystoreTo)
	auth, err := bind.NewTransactor(strings.NewReader(keyTo), "password1")
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	printInfo(setup)

	// Get the latest block number
	shell.AddCmd(&ishell.Cmd{
		Name: "latestBlock",
		Help: "Returns latest block number, arguments: latestBlock",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			c.Println("Get latest block number:")
			latestBlock(client)
			c.Println("===============================================================")
		},
	})

	// Get block N
	shell.AddCmd(&ishell.Cmd{
		Name: "getBlock",
		Help: "Returns a block header, arguments: getBlock [integer]",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			if len(c.Args) == 0 {
				c.Println("Input argument required, e.g.: getBlock 10")
			} else if len(c.Args) > 1 {
				c.Println("Only enter single argument")
			} else {
				getBlock(client, c.Args[0])
			}
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getValidators",
		Help: "Queries the validator contract for the whitelist of validators",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			result, err := validation.GetValidators(&bind.CallOpts{})
			if err != nil {
				fmt.Printf("Error: %s", err)
				return
			}
			c.Println("Validators Whitelist:")
			c.Printf("%x\n", result)

			c.Println("===============================================================")
		},
	})

	// Get block N and spew out the RLP encoded block
	shell.AddCmd(&ishell.Cmd{
		Name: "rlpBlock",
		Help: "Returns RLP encoded block header, arguments: rlpBlock [integer]",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			if len(c.Args) == 0 {
				c.Println("Input argument required, e.g.: rlpBlock 10")
			} else if len(c.Args) > 1 {
				c.Println("Only enter single argument")
			} else {
				c.Println("RLP encode block: " + c.Args[0])
				rlpEncodeBlock(client, c.Args[0])
			}
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "submitRlpBlock",
		Help: "Returns the RLP block header, signed block prefix, extra data prefix and submits to validation contract, arguments: submitRlpBlock [integer]",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			if len(c.Args) == 0 {
				c.Println("Choose a block.")
			} else if len(c.Args) > 1 {
				c.Println("Too many arguments entered.")
			} else {
				c.Println("RLP encode block: " + c.Args[0])
				encodedBlock, prefixBlock, prefixExtra := calculateRlpEncoding(client, c.Args[0])
				res, err := validation.ValidateBlock(auth, encodedBlock, prefixBlock, prefixExtra)
				if err != nil {
					c.Printf("Error: %s", err)
					return
				}
				c.Printf("Transaction Hash: 0x%x\n", res.Hash())
			}
			c.Println("===============================================================")
		},
	})

	// run shell
	shell.Run()
}

func printInfo(setup config.Setup) {
	// display welcome info.
	fmt.Println("===============================================================")
	fmt.Println("Ion Command Line Interface\n")
	fmt.Println("RPC Client [to]:")
	fmt.Println("Listening on: " + setup.AddrTo + ":" + setup.PortTo)
	fmt.Println("User Account: " + setup.AccountTo)
	fmt.Println("Ion Contract: " + setup.Ion)
	fmt.Println("\nRPC Client [from]: ")
	fmt.Println("Listening on: " + setup.AddrFrom + ":" + setup.PortFrom)
	fmt.Println("User Account: " + setup.AccountFrom)
	fmt.Println("===============================================================")
}

func strToHex(input string) (output string) {
	val, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("please input decimal:", err)
		return
	}
	output = strconv.FormatInt(int64(val), 16)

	return "0x" + output
}
