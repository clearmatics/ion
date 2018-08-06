// Copyright (c) 2018 Clearmatics Technologies Ltd

package cli

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/abiosoft/ishell"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/clearmatics/ion/ion-cli/config"
	contract "github.com/clearmatics/ion/ion-cli/contracts"
	"github.com/clearmatics/ion/ion-cli/utils"
)

// Launch - definition of commands and creates the iterface
func Launch(setup config.Setup, clientFrom *ethclient.Client, Validation *contract.Validation, Ion *contract.Ion, Trigger *contract.Trigger) {
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	// Get a suggested gas price
	gasPrice, err := clientFrom.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Create an authorized transactor and spend 1 unicorn
	keyTo := config.ReadString(setup.KeystoreTo)
	authTo, err := bind.NewTransactor(strings.NewReader(keyTo), "password1")
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	authTo.Value = big.NewInt(0) // in wei
	authTo.From = common.HexToAddress(setup.AccountTo)
	authTo.GasLimit = uint64(100000) // in units
	authTo.GasPrice = gasPrice

	// Create an authorized transactor and spend 1 unicorn
	keyFrom := config.ReadString(setup.KeystoreFrom)
	authFrom, err := bind.NewTransactor(strings.NewReader(keyFrom), "password1")
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	authFrom.Value = big.NewInt(0)                         // in wei
	authFrom.From = common.HexToAddress(setup.AccountFrom) // in wei
	authFrom.GasLimit = uint64(100000)                     // in units
	authFrom.GasPrice = gasPrice

	//---------------------------------------------------------------------------------------------
	// 	RPC Client Specific Commands
	//---------------------------------------------------------------------------------------------

	shell.AddCmd(&ishell.Cmd{
		Name: "latestBlock",
		Help: "use: latestBlock  \n\t\t\t\tdescription: Returns number of latest block mined/sealed",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			c.Println("Connecting to: " + setup.AddrFrom)
			c.Println("Get latest block number:")
			lastBlock := latestBlock(clientFrom)
			c.Printf("latest block: %v\n", lastBlock.Number)

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

	//---------------------------------------------------------------------------------------------
	// 	Validation Specific Commands
	//---------------------------------------------------------------------------------------------

	shell.AddCmd(&ishell.Cmd{
		Name: "registerChainValidation",
		Help: "use: registerChainValidation \n\t\t\t\tdescription: Register new chain with validation contract",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			c.Println("Connecting to: " + setup.AddrTo)
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			// Get the chainId
			bytesChainId, err := utils.StringToBytes32(setup.ChainId)
			if err != nil {
				c.Printf("Error: %s", err)
				return
			}

			// Get the validators array
			c.Print("Enter Validators: ")
			validatorString := c.ReadLine()
			valArray := strings.Fields(validatorString)
			var validators []common.Address
			for _, val := range valArray {
				validators = append(validators, common.HexToAddress(val))
			}

			// Get genesis hash
			c.Print("Enter Genesis Hash: ")
			genesis := c.ReadLine()
			bytesGenesis, err := utils.StringToBytes32(genesis)
			if err != nil {
				c.Printf("Error: %s", err)
				return
			}

			ionAddress := common.HexToAddress(setup.Ion)
			tx, err := Validation.RegisterChain(authTo, bytesChainId, ionAddress, validators, bytesGenesis)
			if err != nil {
				c.Printf("Error: %s", err)
				return
			}
			c.Printf("\nTransaction Result:\n%x\n", tx)
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "checkBlockValidation",
		Help: "use: checkBlockValidation [blockHash]\n\t\t\t\tdescription: Returns true for validated blocks",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			c.Println("Connecting to: " + setup.AddrTo)
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			// Get the chainId
			bytesChainId, err := utils.StringToBytes32(setup.ChainId)
			if err != nil {
				c.Printf("Error: %s", err)
				return
			}

			// Get the blockHash
			c.Print("Enter BlockHash: ")
			blockHash := c.ReadLine()
			bytesBlockHash, err := utils.StringToBytes32(blockHash)
			if err != nil {
				c.Printf("Error: %s", err)
				return
			}

			result, err := Validation.MBlockhashes(&bind.CallOpts{}, bytesChainId, bytesBlockHash)
			if err != nil {
				c.Printf("Error: %s", err)
				return
			}
			c.Println("Checking for valid block:")
			c.Printf("ChainId:\t%x\nBlockHash:\t%x\nValid:\t\t%v\n", bytesChainId, bytesBlockHash, result)

			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "submitBlockValidation",
		Help: "use: submitBlockValidation [integer] \n\t\t\t\tdescription: Returns the RLP block header, signed block prefix, extra data prefix and submits to validation contract",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			// Get the chainId
			bytesChainId, err := utils.StringToBytes32(setup.ChainId)
			if err != nil {
				c.Printf("Error: %s", err)
				return
			}

			// Get the block number
			c.Print("Enter Block Number: ")
			blockNum := c.ReadLine()
			c.Printf("RLP encode block:\nNumber:\t\t%s", blockNum)

			signedBlock, unsignedBlock := calculateRlpEncoding(clientFrom, blockNum)
			res, err := Validation.SubmitBlock(authTo, bytesChainId, unsignedBlock, signedBlock)
			if err != nil {
				c.Printf("Error: %s", err)
				return
			}
			c.Printf("\nTransaction Hash:\n0x%x\n", res.Hash())
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "latestValidationBlock",
		Help: "use: latestValidationBlock \n\t\t\t\tdescription: Returns hash of the last block submitted to the validation contract",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			c.Println("Connecting to: " + setup.AddrTo)
			// Get the chainId
			bytesChainId, err := utils.StringToBytes32(setup.ChainId)
			if err != nil {
				c.Printf("Error: %s", err)
				return
			}

			result, err := Validation.MLatestblock(&bind.CallOpts{}, bytesChainId)
			if err != nil {
				c.Printf("Error: %s", err)
				return
			}
			c.Printf("Latest Block Submitted:\nBlockHash:\t0x%x\nChainId:\t%s\n", result, setup.ChainId)
			c.Println("===============================================================")
		},
	})

	//---------------------------------------------------------------------------------------------
	// 	Trigger Specific Commands
	//---------------------------------------------------------------------------------------------
	shell.AddCmd(&ishell.Cmd{
		Name: "triggerEvent",
		Help: "use: triggerEvent \n\t\t\t\tdescription: Returns hash of the last block submitted to the validation contract",
		Func: func(c *ishell.Context) {
			c.Println("===============================================================")
			c.Println("Connecting to: " + setup.AddrFrom)
			lastBlock := latestBlock(clientFrom)

			result, err := Trigger.Fire(authFrom)
			if err != nil {
				c.Printf("Error: %s", err)
				return
			}
			c.Printf("Triggered Event:\nResult:\t%+v\n", result.Hash)
			block := lastBlock.Number
			blockNumber, _ := strconv.ParseUint(block.String(), 0, 64)
			// s := []uint32{}
			ch := make(chan *contract.TriggerTriggered)
			opts := &bind.WatchOpts{}
			opts.Start = &blockNumber
			_, err = Trigger.WatchTriggered(opts, ch)
			if err != nil {
				log.Fatalf("Failed WatchTriggered: %v", err)
			}
			var newEvent *contract.TriggerTriggered = <-ch
			fmt.Println(newEvent.Caller)

			c.Println("===============================================================")
		},
	})

	// shell.AddCmd(&ishell.Cmd{
	// 	Name: "verifyAndExecute",
	// 	Help: "use: verifyAndExecute [Transaction Hash] [Block Number] \n\t\t\t\tdescription: Returns the proof of a specific transaction held within a Patricia trie",
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
	//---------------------------------------------------------------------------------------------

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
