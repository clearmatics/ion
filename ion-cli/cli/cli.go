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
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/clearmatics/ion/ion-cli/config"
	contract "github.com/clearmatics/ion/ion-cli/contracts"
)

// Launch - definition of commands and creates the iterface
func Launch(setup config.Setup, clientFrom *ethclient.Client, Ion *contract.Ion) {
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	// Get a suggested gas price
	gasPrice, err := clientFrom.SuggestGasPrice(context.Background())
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

	shell.AddCmd(&ishell.Cmd{
		Name: "latestBlock",
		Help: "use: latestBlock  \n\t\t\t\tdescription: Returns number of latest block mined/sealed",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			c.Println("Connecting to: " + setup.AddrFrom)
			c.Println("Get latest block number:")
			latestBlock(clientFrom)
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getBlock",
		Help: "use: getBlock [integer] \n\t\t\t\tdescription: Returns block header specified",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			c.Println("Connecting to: " + setup.AddrFrom)
			if len(c.Args) == 0 {
				c.Println("Input argument required, e.g.: getBlock 10")
			} else if len(c.Args) > 1 {
				c.Println("Only enter single argument")
			} else {
				getBlock(clientFrom, c.Args[0])
			}
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "ionChainId",
		Help: "use: ionChainId \n\t\t\t\tdescription: Returns id of Ion chain",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			c.Println("Connecting to: " + setup.AddrTo)
			if len(c.Args) == 0 {
				result, err := Ion.ChainId(&bind.CallOpts{})
				if err != nil {
					c.Printf("Error: %s", err)
					return
				}
				c.Printf("Result:\t%x\n", result)
			} else if len(c.Args) > 0 {
				c.Println("Only enter single argument")
			}
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "ionRegisteredChains",
		Help: "use: ionRegisteredChains \n\t\t\t\tdescription: Returns array of all registered chains",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			c.Println("Connecting to: " + setup.AddrTo)
			if len(c.Args) == 0 {
				result, err := Ion.RegisteredChains(&bind.CallOpts{}, big.NewInt(0))
				if err != nil {
					c.Printf("Error: %s", err)
					return
				}
				c.Printf("Result:\t%s\n", result)
			} else if len(c.Args) > 0 {
				c.Println("Only enter single argument")
			}
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "ionRegisterChain",
		Help: "use: ionRegisterChain \n\t\t\t\tdescription: Register chain with Ion contract",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			c.Println("Connecting to: " + setup.AddrTo)
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			// Get the chainId
			c.Print("New ChainId: ")
			chainId := c.ReadLine()

			c.Printf("Running Command:\t%s\t%s", chainId
			c.Println("===============================================================")
		},
	})

	// shell.AddCmd(&ishell.Cmd{
	// 	Name: "getValidators",
	// 	Help: "use: getValidators \n\t\t\t\tdescription: Returns the whitelist of validators from validator contract",
	// 	Func: func(c *ishell.Context) {
	// 		c.Println("===============================================================")
	// 		c.Println("Connecting to: " + setup.AddrFrom)
	// 		result, err := Validation.GetValidators(&bind.CallOpts{})
	// 		if err != nil {
	// 			fmt.Printf("Error: %s", err)
	// 			return
	// 		}
	// 		c.Println("Validators Whitelist:")
	// 		c.Printf("%x\n", result)

	// 		c.Println("===============================================================")
	// 	},
	// })

	// shell.AddCmd(&ishell.Cmd{
	// 	Name: "latestValidationBlock",
	// 	Help: "use: latestValidationBlock \n\t\t\t\tdescription: Returns hash of the last block submitted to the validation contract",
	// 	Func: func(c *ishell.Context) {
	// 		c.Println("===============================================================")
	// 		c.Println("Connecting to: " + setup.AddrTo)
	// 		result, err := Validation.LatestBlock(&bind.CallOpts{})
	// 		if err != nil {
	// 			fmt.Printf("Error: %s", err)
	// 			return
	// 		}
	// 		c.Println("Last Block Submitted:")
	// 		c.Printf("0x%x\n", result)

	// 		c.Println("===============================================================")
	// 	},
	// })

	// shell.AddCmd(&ishell.Cmd{
	// 	Name: "getValidationBlock",
	// 	Help: "use: latestValidationBlock \n\t\t\t\tdescription: Returns hash of the last block submitted to the validation contract",
	// 	Func: func(c *ishell.Context) {
	// 		c.Println("===============================================================")
	// 		c.Println("Connecting to: " + setup.AddrFrom)
	// 		blockNum := new(big.Int)
	// 		blockNum.SetString(c.Args[0], 10)
	// 		result, err := Validation.GetBlock(&bind.CallOpts{}, blockNum)
	// 		if err != nil {
	// 			fmt.Printf("Error: %s", err)
	// 			return
	// 		}
	// 		c.Println("Last Block Submitted:")
	// 		c.Printf("Block:\n\t%x\n", result.BlockHeight)
	// 		c.Printf("Hash:\n\t0x%x\n", result.BlockHash)
	// 		c.Printf("Parent Hash:\n\t0x%x\n", result.PrevBlockHash)

	// 		c.Println("===============================================================")
	// 	},
	// })

	// shell.AddCmd(&ishell.Cmd{
	// 	Name: "submitValidationBlock",
	// 	Help: "use: submitValidationBlock [integer] \n\t\t\t\tdescription: Returns the RLP block header, signed block prefix, extra data prefix and submits to validation contract",
	// 	Func: func(c *ishell.Context) {
	// 		c.Println("===============================================================")
	// 		c.Println("Connecting to: " + setup.AddrTo)
	// 		if len(c.Args) == 0 {
	// 			c.Println("Select a block")
	// 		} else if len(c.Args) > 1 {
	// 			c.Println("Too many arguments entered.")
	// 		} else {
	// 			c.Println("RLP encode block: " + c.Args[0])
	// 			encodedBlock, prefixBlock, prefixExtra := calculateRlpEncoding(clientFrom, c.Args[0])
	// 			res, err := Validation.ValidateBlock(auth, encodedBlock, prefixBlock, prefixExtra)
	// 			if err != nil {
	// 				c.Printf("Error: %s", err)
	// 				return
	// 			}
	// 			c.Printf("\nTransaction Hash:\n0x%x\n", res.Hash())
	// 		}
	// 		c.Println("===============================================================")
	// 	},
	// })

	// shell.AddCmd(&ishell.Cmd{
	// 	Name: "generateTxProof",
	// 	Help: "use: generateTxProof [Transaction Hash] [Block Number] \n\t\t\t\tdescription: Returns the proof of a specific transaction held within a Patricia trie",
	// 	Func: func(c *ishell.Context) {
	// 		c.Println("===============================================================")
	// 		c.Println("Connecting to: " + setup.AddrTo)
	// 		if len(c.Args) == 0 {
	// 			c.Println("Select a block")
	// 		} else if len(c.Args) > 2 {
	// 			c.Println("Too many arguments entered.")
	// 		} else {
	// 			c.Println("RLP encode block: " + c.Args[0])
	// 			rootHash, idx, leaf, proof := GenerateTxProof(clientFrom, c.Args[0], c.Args[1])
	// 			c.Printf("\nRoot Hash:\n% 0x\nTransaction Index:\n% 0x\nTransaction Leaf:\n% 0x\nProof:\n% 0x\n", rootHash, idx, leaf, proof)
	// 		}
	// 		c.Println("===============================================================")
	// 	},
	// })

	//---------------------------------------------------------------------------------------------
	// 	Ion Specific Commands
	//
	//---------------------------------------------------------------------------------------------
	// shell.AddCmd(&ishell.Cmd{
	// 	Name: "ionRegisterChain",
	// 	Help: "use: ionRegisterChain [Transaction Hash] [Block Number] \n\t\t\t\tdescription: Returns the proof of a specific transaction held within a Patricia trie",
	// 	Func: func(c *ishell.Context) {
	// 		c.Println("===============================================================")
	// 		c.Println("Connecting to: " + setup.AddrTo)
	// 		if len(c.Args) == 0 {
	// 			c.Println("Select a block")
	// 		} else if len(c.Args) > 2 {
	// 			c.Println("Too many arguments entered.")
	// 		} else {
	// 			c.Println("RLP encode block: " + c.Args[0])
	// 			rootHash, idx, leaf, proof := GenerateTxProof(clientFrom, c.Args[0], c.Args[1])
	// 			c.Printf("\nRoot Hash:\n% 0x\nTransaction Index:\n% 0x\nTransaction Leaf:\n% 0x\nProof:\n% 0x\n", rootHash, idx, leaf, proof)
	// 		}
	// 		c.Println("===============================================================")
	// 	},
	// })

	// run shell
	shell.Run()
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
