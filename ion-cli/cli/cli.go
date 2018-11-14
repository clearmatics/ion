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
	config config.Setup,
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
	authTo, keyTo := config.InitUser(config.KeystoreTo, config.PasswordTo)
	authTo.Value = big.NewInt(0)     // in wei
	authTo.GasLimit = uint64(100000) // in units
	authTo.GasPrice = gasPrice

	// Create an authorized transactor and spend 1 unicorn
	authFrom, keyFrom := config.InitUser(config.KeystoreFrom, config.PasswordFrom)
	authFrom.Value = big.NewInt(0)     // in wei
	authFrom.GasLimit = uint64(100000) // in units
	authFrom.GasPrice = gasPrice

	//---------------------------------------------------------------------------------------------
	// 	RPC Client Specific Commands
	//---------------------------------------------------------------------------------------------

	shell.AddCmd(&ishell.Cmd{
		Name: "getBlock",
		Help: "use: \tgetBlock [TO/FROM] [integer] \n\t\t\t\tdescription: Returns block header specified from chain [TO/FROM]",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
				c.Println("Error: Incorrect Arguments!")
			} else if c.Args[0] == "TO" {
				c.Println("Connecting to: " + config.AddrTo)
				getBlock(ethclientTo, c.Args[1])
			} else if c.Args[0] == "FROM" {
				c.Println("Connecting to: " + config.AddrFrom)
				getBlock(ethclientFrom, c.Args[1])
			}
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
        Name: "getTxProof",
        Help: "use: \tgetTxProof [Transaction Hash] \n\t\t\t\tdescription: Returns the proof of a specific transaction held within a Patricia trie",
        Func: func(c *ishell.Context) {
            c.Println("Connecting to: " + config.AddrTo + " and " + config.AddrFrom)
            c.ShowPrompt(false)
            defer c.ShowPrompt(true) // yes, revert after login.

            // Get the chainId
            bytesChainId := common.HexToHash(config.ChainId)

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

            c.Printf("  Path:           0x%x\n
                        TxValue:        0x%x\n
                        TxNodes:        0x%x\n
                        ReceiptValue:   0x%x\n
                        ReceiptNodes:   0x%x\n", tx.Path, txValue, txNodes, receiptValue, receiptNodes)
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
			c.Println("Connecting to: " + config.AddrTo)
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)

			// Get the chainId
			bytesChainId := common.HexToHash(config.ChainId)

			// Get the validators array
			c.Print("Enter Validators: ")
			validatorString := c.ReadLine()
			valArray := strings.Fields(validatorString)
			// valArray := strings.Fields("0x42eb768f2244c8811c63729a21a3569731535f06 0x6635f83421bf059cd8111f180f0727128685bae4 0x7ffc57839b00206d1ad20c69a1981b489f772031 0xb279182d99e65703f0076e4812653aab85fca0f0 0xd6ae8250b8348c94847280928c79fb3b63ca453e 0xda35dee8eddeaa556e4c26268463e26fb91ff74f 0xfc18cbc391de84dbd87db83b20935d3e89f5dd91")
			var validators []common.Address
			for _, val := range valArray {
				validators = append(validators, common.HexToAddress(val))
			}

			// Get genesis hash
			c.Print("Enter Genesis Hash: ")
			genesis := c.ReadLine()
			bytesGenesis := common.HexToHash(genesis)
			// bytesGenesis := common.HexToHash("0x100dc525cdcb7933e09f10d4019c38d342253a0aa32889fbbdbc5f2406c7546c")

			tx := contract.RegisterChain(
				ctx,
				ethclientTo,
				keyTo.PrivateKey,
				Validation,
				common.HexToAddress(config.Validation),
				bytesChainId,
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
			c.Println("Connecting to: " + config.AddrTo)
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			// Get the chainId
			bytesChainId, err := utils.StringToBytes32(config.ChainId)
			if err != nil {
				c.Printf("Error: %s", err)
				return
			}

			// Get the block number
			c.Print("Enter Block Number: ")
			blockNum := c.ReadLine()
			// blockNum := "2776659"
			c.Printf("RLP encode block:\nNumber:\t\t%s", blockNum)

			signedBlock, unsignedBlock := calculateRlpEncoding(ethclientFrom, blockNum)
			tx := contract.SubmitBlock(
				ctx,
				ethclientTo,
				keyTo.PrivateKey,
				Validation,
				common.HexToAddress(config.Validation),
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
			c.Println("Connecting to: " + config.AddrTo)
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			// Get the chainId
			bytesChainId, err := utils.StringToBytes32(config.ChainId)
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
				common.HexToAddress(config.AddrTo),
				common.HexToAddress(config.Validation),
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
			c.Println("Connecting to: " + config.AddrTo)
			// Get the chainId
			bytesChainId := common.HexToHash(config.ChainId)

			result := contract.LatestValidBlock(
				ctx,
				ethclientTo,
				Validation,
				common.HexToAddress(config.AddrTo),
				common.HexToAddress(config.Validation),
				bytesChainId,
			)

			c.Println("Checking for latest valid block:")
			c.Printf("\nBlockHash:\t0x%x\nChainId:\t%s\n", result, config.ChainId)
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
			c.Println("Connecting to: " + config.AddrFrom)

			tx := contract.Fire(
				ctx,
				ethclientFrom,
				keyFrom.PrivateKey,
				Trigger,
				common.HexToAddress(config.Trigger),
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
			c.Println("Connecting to: " + config.AddrTo + " and " + config.AddrFrom)
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			// Get the chainId
			bytesChainId := common.HexToHash(config.ChainId)

			// Get the transaction hash
			c.Print("Enter Transaction Hash: ")
			txHash := c.ReadLine()
			bytesTxHash := common.HexToHash(txHash)
			// bytesTxHash := common.HexToHash("0x5da684940b4fd9dec708cc159dc504aa01e90d40bb76a2b73299aee13aa72098")

			// Get the blockHash
			c.Print("Enter Block Hash: ")
			blockHash := c.ReadLine()
			bytesBlockHash := common.HexToHash(blockHash)
			// bytesBlockHash := common.HexToHash("0x74d37aa3c96bc98903451d0baf051b87550191aa0d92032f7406a4984610b046")

			// Generate the proof
			txPath, txValue, txNodes, receiptValue, receiptNodes := utils.GenerateProof(
				ctx,
				clientFrom,
				bytesTxHash,
			)

			// Execute
			tx := contract.VerifyExecute(
				ctx,
				ethclientTo,
				keyFrom.PrivateKey,
				Function,
				common.HexToAddress(config.Function),
				bytesChainId,
				bytesBlockHash,
				common.HexToAddress(config.Trigger), // TRIG_DEPLOYED_RINKEBY_ADDR,
				txPath,                                 // TEST_PATH,
				txValue,                                // TEST_TX_VALUE,
				txNodes,                                // TEST_TX_NODES,
				receiptValue,                           // TEST_RECEIPT_VALUE,
				receiptNodes,                           // TEST_RECEIPT_NODES,
				common.HexToAddress(config.AccountFrom), // TRIG_CALLED_BY,
				nil,
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
