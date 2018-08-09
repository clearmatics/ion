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
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/abiosoft/ishell"

	"github.com/clearmatics/ion/ion-cli/config"
	contract "github.com/clearmatics/ion/ion-cli/contracts"
	"github.com/clearmatics/ion/ion-cli/utils"
)

// Launch - definition of commands and creates the iterface
func Launch(
	setup config.Setup,
	clientTo *rpc.Client,
	clientFrom *rpc.Client,
	Validation *compiler.Contract,
	Trigger *compiler.Contract,
	Function *compiler.Contract,
) {
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	// Create new context
	ctx := context.Background()

	ethclientTo := ethclient.NewClient(clientTo)
	ethclientFrom := ethclient.NewClient(clientFrom)

	// Get a suggested gas price
	gasPrice, err := ethclientFrom.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Create an authorized transactor and corrsponding privateKey
	authTo, keyTo := config.InitUser(setup.KeystoreTo, setup.PasswordTo)
	authTo.Value = big.NewInt(0)     // in wei
	authTo.GasLimit = uint64(100000) // in units
	authTo.GasPrice = gasPrice

	// Create an authorized transactor and spend 1 unicorn
	authFrom, keyFrom := config.InitUser(setup.KeystoreFrom, setup.PasswordFrom)
	authFrom.Value = big.NewInt(0)     // in wei
	authFrom.GasLimit = uint64(100000) // in units
	authFrom.GasPrice = gasPrice

	//---------------------------------------------------------------------------------------------
	// 	RPC Client Specific Commands
	//---------------------------------------------------------------------------------------------

	shell.AddCmd(&ishell.Cmd{
		Name: "latestBlock",
		Help: "use: \tlatestBlock [TO/FROM] \n\t\t\t\tdescription: Returns number of latest block mined/sealed from chain [TO/FROM]",
		Func: func(c *ishell.Context) {
			if c.Args[0] == "TO" {
				c.Println("Connecting to: " + setup.AddrTo)
				c.Println("Get latest block number:")
				lastBlock := latestBlock(ethclientTo)
				c.Printf("latest block: %v\n", lastBlock.Number)
			} else if c.Args[0] == "FROM" {
				c.Println("Connecting to: " + setup.AddrFrom)
				c.Println("Get latest block number:")
				lastBlock := latestBlock(ethclientFrom)
				c.Printf("latest block: %v\n", lastBlock.Number)
			} else {
				c.Println("Please choose enter TO or FROM only!")
			}
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getBlock",
		Help: "use: \tgetBlock [TO/FROM] [integer] \n\t\t\t\tdescription: Returns block header specified from chain [TO/FROM]",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
				c.Println("Error: Incorrect Arguments!")
			} else if c.Args[0] == "TO" {
				c.Println("Connecting to: " + setup.AddrTo)
				getBlock(ethclientTo, c.Args[1])
			} else if c.Args[0] == "FROM" {
				c.Println("Connecting to: " + setup.AddrFrom)
				getBlock(ethclientFrom, c.Args[1])
			}
			c.Println("===============================================================")
		},
	})

	//---------------------------------------------------------------------------------------------
	// 	Validation Specific Commands
	//---------------------------------------------------------------------------------------------
	shell.AddCmd(&ishell.Cmd{
		Name: "registerChainValidation",
		Help: "use: \tregisterChainValidation\n \t\t\t\t\tEnter Validators: [ADDRESS ADDRESS]\n \t\t\t\t\tEnter Genesis Hash: [HASH] \n\t\t\t\tdescription: Register new chain with validation contract",
		Func: func(c *ishell.Context) {
			c.Println("Connecting to: " + setup.AddrTo)
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)

			// Get the chainId
			bytesChainId := common.HexToHash(setup.ChainId)

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
			bytesGenesis := common.HexToHash(genesis)

			tx := contract.RegisterChain(
				ctx,
				ethclientTo,
				keyTo.PrivateKey,
				Validation,
				common.HexToAddress(setup.Validation),
				bytesChainId,
				common.HexToAddress(setup.Ion),
				validators,
				bytesGenesis,
			)

			c.Printf("Transaction Hash:\n0x%x\n", tx.Hash())
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "submitBlockValidation",
		Help: "use: \tsubmitBlockValidation\n \t\t\t\t\tEnter Block Number: [INTEGER]\n\t\t\t\tdescription: Returns the RLP block header, signed block prefix, extra data prefix and submits to validation contract",
		Func: func(c *ishell.Context) {
			c.Println("Connecting to: " + setup.AddrTo)
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

			signedBlock, unsignedBlock := calculateRlpEncoding(ethclientFrom, blockNum)
			tx := contract.SubmitBlock(
				ctx,
				ethclientTo,
				keyTo.PrivateKey,
				Validation,
				common.HexToAddress(setup.Validation),
				bytesChainId,
				unsignedBlock,
				signedBlock,
			)

			c.Printf("Transaction Hash:\n0x%x\n", tx.Hash())
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "checkBlockValidation",
		Help: "use: \tcheckBlockValidation\n \t\t\t\t\tEnter Blockhash: [HASH]\n\t\t\t\tdescription: Returns true for validated blocks",
		Func: func(c *ishell.Context) {
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

			result := contract.ValidBlock(
				ctx,
				ethclientTo,
				Validation,
				common.HexToAddress(setup.AddrTo),
				common.HexToAddress(setup.Validation),
				bytesChainId,
				bytesBlockHash,
			)

			c.Println("Checking for valid block:")
			c.Printf("ChainId:\t%x\nBlockHash:\t%x\nValid:\t\t%v\n", bytesChainId, bytesBlockHash, result)
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "latestValidatedBlock",
		Help: "use: \tlatestValidatedBlock \n\t\t\t\tdescription: Returns hash of the last block submitted to the validation contract",
		Func: func(c *ishell.Context) {
			c.Println("Connecting to: " + setup.AddrTo)
			// Get the chainId
			bytesChainId := common.HexToHash(setup.ChainId)

			result := contract.LatestValidBlock(
				ctx,
				ethclientTo,
				Validation,
				common.HexToAddress(setup.AddrTo),
				common.HexToAddress(setup.Validation),
				bytesChainId,
			)

			c.Println("Checking for latest valid block:")
			c.Printf("\nBlockHash:\t0x%x\nChainId:\t%s\n", result, setup.ChainId)
			c.Println("===============================================================")
		},
	})

	//---------------------------------------------------------------------------------------------
	// 	Trigger Specific Commands
	//---------------------------------------------------------------------------------------------
	shell.AddCmd(&ishell.Cmd{
		Name: "triggerEvent",
		Help: "use: \ttriggerEvent \n\t\t\t\tdescription: Returns hash of the last block submitted to the validation contract",
		Func: func(c *ishell.Context) {
			c.Println("Connecting to: " + setup.AddrFrom)

			tx := contract.Fire(
				ctx,
				ethclientFrom,
				keyFrom.PrivateKey,
				Trigger,
				common.HexToAddress(setup.Trigger),
			)

			c.Printf("Transaction Hash:\n0x%x\n", tx.Hash())
			c.Println("===============================================================")
		},
	})

	//---------------------------------------------------------------------------------------------
	// 	Function Specific Commands
	//---------------------------------------------------------------------------------------------
	shell.AddCmd(&ishell.Cmd{
		Name: "verifyAndExecute",
		Help: "use: \tverifyAndExecute [Transaction Hash] \n\t\t\t\tdescription: Returns the proof of a specific transaction held within a Patricia trie",
		Func: func(c *ishell.Context) {
			c.Println("Connecting to: " + setup.AddrTo + " and " + setup.AddrFrom)
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			// Get the chainId
			bytesChainId := common.HexToHash(setup.ChainId)

			// Get the transaction hash
			c.Print("Enter Transaction Hash: ")
			txHash := c.ReadLine()
			bytesTxHash := common.HexToHash(txHash)

			// Get the blockHash
			c.Print("Enter Block Hash: ")
			blockHash := c.ReadLine()
			bytesBlockHash := common.HexToHash(blockHash)

			// Generate the proof
			txPath, txValue, txNodes, receiptValue, receiptNodes := utils.GenerateProof(
				ctx,
				clientFrom,
				bytesTxHash,
			)

			amount := big.NewInt(10000000)

			// Execute
			tx := contract.VerifyExecute(
				ctx,
				ethclientTo,
				keyTo.PrivateKey,
				Function,
				common.HexToAddress(setup.Function),
				bytesChainId,
				bytesBlockHash,
				common.HexToAddress(setup.Trigger), // TRIG_DEPLOYED_RINKEBY_ADDR,
				txPath,                              // TEST_PATH,
				txValue,                             // TEST_TX_VALUE,
				txNodes,                             // TEST_TX_NODES,
				receiptValue,                        // TEST_RECEIPT_VALUE,
				receiptNodes,                        // TEST_RECEIPT_NODES,
				common.HexToAddress(setup.AddrFrom), // TRIG_CALLED_BY,
				amount,
			)
			c.Printf("Transaction Hash:\n0x%x\n", tx.Hash())
			c.Println("===============================================================")
		},
	})

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
