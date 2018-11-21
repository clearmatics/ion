// Copyright (c) 2018 Clearmatics Technologies Ltd
package cli

import (
	"context"
	"fmt"
	"encoding/hex"
	"math/big"
	"strconv"
	"errors"
	"strings"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/abiosoft/ishell"

	"github.com/clearmatics/ion/ion-cli/config"
	contract "github.com/clearmatics/ion/ion-cli/contracts"
	"github.com/clearmatics/ion/ion-cli/utils"
)

// Launch - definition of commands and creates the interface
func Launch(
	configuration config.Setup,
	clientTo *rpc.Client,
	clientFrom *rpc.Client,
) {
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	// Create new context
	ctx := context.Background()

	var ethClient *EthClient = nil
	var contracts map[string]*contract.ContractInstance = make(map[string]*contract.ContractInstance)
	var accounts map[string]*config.Account = make(map[string]*config.Account)

	//---------------------------------------------------------------------------------------------
	// 	RPC Client Specific Commands
	//---------------------------------------------------------------------------------------------

	shell.AddCmd(&ishell.Cmd{
		Name: "connectToClient",
		Help: "use: \tconnectToClient [rpc url] \n\t\t\t\tdescription: Connects to an RPC client to be used",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
                c.Println("Usage: \tconnectToClient [rpc url] \n")
			} else {
			    c.Println("Connecting to client...\n")
			    ethClient = getClient(c.Args[0])
			    c.Println("Connected!")
			}
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "addContractInstance",
		Help: "use: \taddContractInstance [name] [path/to/solidity/contract]\n\t\t\t\tdescription: Compiles a contract for use",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
                c.Println("Usage: \taddContractInstances [name] [path/to/solidity/contract]\n")
			} else {
                err := addContractInstance(c.Args[1], c.Args[0], contracts)
                if err != nil {
                    c.Println(err)
                    return
                }
			}
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "listContracts",
		Help: "use: \tlistContracts \n\t\t\t\tdescription: List compiled contract instances",
		Func: func(c *ishell.Context) {
            for key := range contracts {
                c.Println(key)
            }
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "addAccount",
		Help: "use: \taddAccount [name] [path/to/keystore]\n\t\t\t\tdescription: Add account to be used for transactions",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
                c.Println("Usage: \taddAccount [name] [path/to/keystore]\n")
			} else {
                c.ShowPrompt(false)
                defer c.ShowPrompt(true)
			    c.Println("Please provide your key decryption password.")
			    input := c.ReadPassword()
                auth, key, err := config.InitUser(c.Args[1], input)
                if err != nil {
                    c.Println(err)
                    return
                }
                account := &config.Account{Auth: auth, Key: key}
                accounts[c.Args[0]] = account

                c.Println("Account added succesfully.")
			}
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "listAccounts",
		Help: "use: \tlistAccounts \n\t\t\t\tdescription: List all added accounts",
		Func: func(c *ishell.Context) {
            for key := range accounts {
                c.Println(key)
            }
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "deployContract",
		Help: "use: \tdeployContract [contract name] [account name] \n\t\t\t\tdescription: Deploys specified contract instance to connected client",
		Func: func(c *ishell.Context) {
            if len(c.Args) != 2 {
                c.Println("Usage: \tdeployContract [contract name] [account name] \n")
            } else {
			    if ethClient == nil {
			        c.Println("Please connect to a Client before invoking this function.\nUse \tconnectToClient [rpc url] \n")
			        return
			    }
                contractInstance := contracts[c.Args[0]]
                if contractInstance == nil {
                    errStr := fmt.Sprintf("Contract instance %s not found.\nUse \taddContractInstances [name] [path/to/solidity/contract] [deployed address] \n", c.Args[0])
                    c.Println(errStr)
			        return
                }

                binStr, abiStr := contract.GetContractBytecodeAndABI(contractInstance.Contract)

                account := accounts[c.Args[1]]
                if account == nil {
                    errStr := fmt.Sprintf("Account %s not found.\nUse \taddAccount [name] [path/to/keystore] [password] \n", c.Args[1])
			        c.Println(errStr)
			        return
                }

                tx := contract.CompileAndDeployContract(
                    ctx,
                    ethClient.client,
                    account.Key.PrivateKey,
                    binStr,
                    abiStr,
                    nil,
                    uint64(3000000),
                )

                c.Println("Waiting for contract to be deployed")
                addr, err := bind.WaitDeployed(ctx, ethClient.client, tx)
                if err != nil {
                    c.Println(err)
                    return
                }
                c.Printf("Deployed contract at: %s\n", addr.String())
            }
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "messageCallFunction",
		Help: "use: \tmessageCallFunction [contract name] [function name] [from account name] [deployed contract address] [amount] [gasLimit] \n\t\t\t\tdescription: Connects to an RPC client to be used",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 6 {
                c.Println("Usage: \tmessageCallFunction [contract name] [function name] [from account name] [deployed contract address] [amount] [gasLimit] \n")
			} else {
			    if ethClient == nil {
			        c.Println("Please connect to a Client before invoking this function.\nUse \tconnectToClient [rpc url] \n")
			        return
			    }

                instance := contracts[c.Args[0]]
                methodName := c.Args[1]
                account := accounts[c.Args[2]]
                contractDeployedAddress := common.HexToAddress(c.Args[3])

                if instance == nil {
                    errStr := fmt.Sprintf("Contract instance %s not found.\nUse \taddContractInstances [name] [path/to/solidity/contract] [deployed address] \n", c.Args[0])
                    c.Println(errStr)
			        return
                }
                if account == nil {
                    errStr := fmt.Sprintf("Account %s not found.\nUse \taddAccount [name] [path/to/keystore] [password] \n", c.Args[2])
			        c.Println(errStr)
			        return
                }

                amount := new(big.Int)
                amount, ok := amount.SetString(c.Args[4], 10)
                if !ok {
                    c.Err(errors.New("Please enter an integer for <amount>"))
                }
                gasLimit, err := strconv.ParseUint(c.Args[5], 10, 64)
                if err != nil {
                    c.Err(errors.New("Please enter an integer for <gasLimit>"))
                }

                if instance.Abi.Methods[methodName].Name == "" {
                    c.Printf("Method name \"%s\" not found for contract \"%s\"\n", methodName, c.Args[0])
                    return
                }

                inputs, err := parseMethodParameters(c, instance.Abi, methodName)
                if err != nil {
                    c.Printf("Error parsing parameters: %s\n", err)
                    return
                }

                tx, err := contract.TransactionContract(
                    ctx,
                    ethClient.client,
                    account.Key.PrivateKey,
                    instance.Contract,
                    contractDeployedAddress,
                    amount,
                    gasLimit,
                    c.Args[1],
                    inputs...
                )
                 if err != nil {
                    c.Println(err)
                    return
                 } else {
                    c.Println("Waiting for transaction to be mined...")
                    receipt, err := bind.WaitMined(ctx, ethClient.client, tx)
                    if err != nil {
                        c.Println(err)
                        return
                    }
                    c.Printf("Transaction hash: %s\n", receipt.TxHash.String())
                 }
			}
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getBlockByNumber",
		Help: "use: \tgetBlockByNumber [rpc url] [integer] \n\t\t\t\tdescription: Returns block header specified from chain [TO/FROM]",
		Func: func(c *ishell.Context) {
            if len(c.Args) == 1 {
                if ethClient != nil {
			        getBlockByNumber(ethClient, c.Args[0])
                } else {
			        c.Println("Please connect to a Client before invoking this function.\nUse \tconnectToClient [rpc url] \n")
			        return
                }
            } else if len(c.Args) == 2 {
			    getBlockByNumber(getClient(c.Args[0]), c.Args[1])
            } else {
                c.Println("Usage: \tgetBlock [rpc url] [integer] \n")
            }
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getBlockByHash",
		Help: "use: \tgetBlockByNumber [rpc url] [hash] \n\t\t\t\tdescription: Returns block header specified from chain [TO/FROM]",
		Func: func(c *ishell.Context) {
            if len(c.Args) == 1 {
                if ethClient != nil {
			        getBlockByHash(ethClient, c.Args[0])
                } else {
			        c.Println("Please connect to a Client before invoking this function.\nUse \tconnectToClient [rpc url] \n")
			        return
                }
            } else if len(c.Args) == 2 {
			    getBlockByHash(getClient(c.Args[0]), c.Args[1])
            } else {
                c.Println("Usage: \tgetBlock [optional rpc url] [hash] \n")
            }
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
        Name: "getProof",
        Help: "use: \tgetProof [optional rpc url] [Transaction Hash] \n\t\t\t\tdescription: Returns a merkle patricia proof of a specific transaction and its receipt in a block",
        Func: func(c *ishell.Context) {
            if len(c.Args) == 1 {
                if ethClient != nil {
                    getProof(ethClient, c.Args[0])
                } else {
			        c.Println("Please connect to a Client before invoking this function.\nUse \tconnectToClient [rpc url] \n")
			        return
                }
            } else if len(c.Args) == 2 {
                getProof(getClient(c.Args[0]), c.Args[1])
            } else {
                c.Println("Usage: \tgetBlock [optional rpc url] [hash] \n")
            }
            c.Println("===============================================================")
        },
    })

	//---------------------------------------------------------------------------------------------
	// 	Validation Specific Commands
	//---------------------------------------------------------------------------------------------
	/*
	shell.AddCmd(&ishell.Cmd{
		Name: "registerChainValidation",
		Help: "use: \tregisterChainValidation\n \t\t\t\t\tEnter Validators: [ADDRESS ADDRESS]\n \t\t\t\t\tEnter Genesis Hash: [HASH] \n\t\t\t\tdescription: Register new chain with validation contract",
		Func: func(c *ishell.Context) {
			c.Println("Connecting to: " + configuration.AddrTo)
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)

			// Get the chainId
			bytesChainId := common.HexToHash(configuration.ChainId)

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
				common.HexToAddress(configuration.Validation),
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
			c.Println("Connecting to: " + configuration.AddrTo)
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			// Get the chainId
			bytesChainId, err := utils.StringToBytes32(configuration.ChainId)
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
				common.HexToAddress(configuration.Validation),
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
			c.Println("Connecting to: " + configuration.AddrTo)
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			// Get the chainId
			bytesChainId, err := utils.StringToBytes32(configuration.ChainId)
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
				common.HexToAddress(configuration.AddrTo),
				common.HexToAddress(configuration.Validation),
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
			c.Println("Connecting to: " + configuration.AddrTo)
			// Get the chainId
			bytesChainId := common.HexToHash(configuration.ChainId)

			result := contract.LatestValidBlock(
				ctx,
				ethclientTo,
				Validation,
				common.HexToAddress(configuration.AddrTo),
				common.HexToAddress(configuration.Validation),
				bytesChainId,
			)

			c.Println("Checking for latest valid block:")
			c.Printf("\nBlockHash:\t0x%x\nChainId:\t%s\n", result, configuration.ChainId)
			c.Println("===============================================================")
		},
	})*/

	//---------------------------------------------------------------------------------------------
	// 	Trigger Specific Commands
	//---------------------------------------------------------------------------------------------
	/*
	shell.AddCmd(&ishell.Cmd{
		Name: "triggerEvent",
		Help: "use: \ttriggerEvent \n\t\t\t\tdescription: Returns hash of the last block submitted to the validation contract",
		Func: func(c *ishell.Context) {
			c.Println("Connecting to: " + configuration.AddrFrom)

			tx := contract.Fire(
				ctx,
				ethclientFrom,
				keyFrom.PrivateKey,
				Trigger,
				common.HexToAddress(configuration.Trigger),
			)

			c.Printf("Transaction Hash:\n0x%x\n", tx.Hash())
			c.Println("===============================================================")
		},
	})
*/

	//---------------------------------------------------------------------------------------------
	// 	Function Specific Commands
	//---------------------------------------------------------------------------------------------
	/*
	shell.AddCmd(&ishell.Cmd{
		Name: "verifyAndExecute",
		Help: "use: \tverifyAndExecute [Transaction Hash] \n\t\t\t\tdescription: Returns the proof of a specific transaction held within a Patricia trie",
		Func: func(c *ishell.Context) {
			c.Println("Connecting to: " + configuration.AddrTo + " and " + configuration.AddrFrom)
			c.ShowPrompt(false)
			defer c.ShowPrompt(true) // yes, revert after login.

			// Get the chainId
			bytesChainId := common.HexToHash(configuration.ChainId)

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
				common.HexToAddress(configuration.Function),
				bytesChainId,
				bytesBlockHash,
				common.HexToAddress(configuration.Trigger), // TRIG_DEPLOYED_RINKEBY_ADDR,
				txPath,                                 // TEST_PATH,
				txValue,                                // TEST_TX_VALUE,
				txNodes,                                // TEST_TX_NODES,
				receiptValue,                           // TEST_RECEIPT_VALUE,
				receiptNodes,                           // TEST_RECEIPT_NODES,
				common.HexToAddress(configuration.AccountFrom), // TRIG_CALLED_BY,
				nil,
			)

			c.Printf("Transaction Hash:\n0x%x\n", tx.Hash())
			c.Println("===============================================================")
		},
	})

    */
	// run shell
	shell.Run()
}

func getClient(url string) (client *EthClient) {
    rpc := utils.ClientRPC(url)
    eth := ethclient.NewClient(rpc)

    return &EthClient{client: eth, rpcClient: rpc, url: url}
}

func parseMethodParameters(c *ishell.Context, abiStruct *abi.ABI, methodName string) (args []interface{}, err error) {
    inputParameters := abiStruct.Methods[methodName].Inputs
    c.ShowPrompt(false)
    defer c.ShowPrompt(true)

    for i := 0; i < len(inputParameters); i++ {
        argument := inputParameters[i]
        c.Printf("Enter input data for parameter %s:\n", argument.Name)

        input := c.ReadLine()

        // bytes = []byte{} argument type = slice, no element, type equates to []uint8
        // byte[] = [][1]byte{} argument type = slice, element type = array, type equates to [][1]uint8
        // byte = bytes1
        // bytesn = [n]byte{} 0 < n < 33, argument type = array, no element, type equates to [n]uint8
        // bytesn[] = [][n]byte{} argument type = slice, element type = array, type equares to [][n]uint8
        // bytesn[m] = [m][n]byte{} argument type = array, element type = array, type equates to [m][n]uint8
        // Many annoying cases of byte arrays

        if argument.Type.Kind == reflect.Array || argument.Type.Kind == reflect.Slice {
            c.Println("Argument is array\n")

            // One dimensional byte array
            // Accepts all byte arrays as hex string with pre-pended '0x' only
            if argument.Type.Elem == nil {
                if argument.Type.Type == reflect.TypeOf(common.Address{}) {
                    // address solidity type
                    item, err := utils.ConvertToType(input, &argument.Type)
                    if err != nil {
                        c.Err(err)
                    }
                    args = append(args, item)
                    continue
                } else if argument.Type.Type == reflect.TypeOf([]byte{}) {
                    // bytes solidity type
                    bytes, err := hex.DecodeString(input[2:])
                    if err != nil {
                        c.Err(err)
                    }
                    args = append(args, bytes)
                    continue
                } else {
                    // Fixed byte array of size n; bytesn solidity type
                    // Any submitted bytes longer than the expected size will be truncated

                    bytes, err := hex.DecodeString(input[2:])
                    if err != nil {
                        c.Err(err)
                    }

                    // Fixed sized arrays can't be created with variables as size
                    switch argument.Type.Size {
                    case 1:
                        var byteArray [1]byte
                        copy(byteArray[:], bytes[:1])
                        args = append(args, byteArray)
                    case 2:
                        var byteArray [2]byte
                        copy(byteArray[:], bytes[:2])
                        args = append(args, byteArray)
                    case 3:
                        var byteArray [3]byte
                        copy(byteArray[:], bytes[:3])
                        args = append(args, byteArray)
                    case 4:
                        var byteArray [4]byte
                        copy(byteArray[:], bytes[:4])
                        args = append(args, byteArray)
                    case 5:
                        var byteArray [5]byte
                        copy(byteArray[:], bytes[:5])
                        args = append(args, byteArray)
                    case 6:
                        var byteArray [6]byte
                        copy(byteArray[:], bytes[:6])
                        args = append(args, byteArray)
                    case 7:
                        var byteArray [7]byte
                        copy(byteArray[:], bytes[:7])
                        args = append(args, byteArray)
                    case 8:
                        var byteArray [8]byte
                        copy(byteArray[:], bytes[:8])
                        args = append(args, byteArray)
                    case 9:
                        var byteArray [9]byte
                        copy(byteArray[:], bytes[:9])
                        args = append(args, byteArray)
                    case 10:
                        var byteArray [10]byte
                        copy(byteArray[:], bytes[:10])
                        args = append(args, byteArray)
                    case 11:
                        var byteArray [11]byte
                        copy(byteArray[:], bytes[:11])
                        args = append(args, byteArray)
                    case 12:
                        var byteArray [12]byte
                        copy(byteArray[:], bytes[:12])
                        args = append(args, byteArray)
                    case 13:
                        var byteArray [13]byte
                        copy(byteArray[:], bytes[:13])
                        args = append(args, byteArray)
                    case 14:
                        var byteArray [14]byte
                        copy(byteArray[:], bytes[:14])
                        args = append(args, byteArray)
                    case 15:
                        var byteArray [15]byte
                        copy(byteArray[:], bytes[:15])
                        args = append(args, byteArray)
                    case 16:
                        var byteArray [16]byte
                        copy(byteArray[:], bytes[:16])
                        args = append(args, byteArray)
                    case 17:
                        var byteArray [17]byte
                        copy(byteArray[:], bytes[:17])
                        args = append(args, byteArray)
                    case 18:
                        var byteArray [18]byte
                        copy(byteArray[:], bytes[:18])
                        args = append(args, byteArray)
                    case 19:
                        var byteArray [19]byte
                        copy(byteArray[:], bytes[:19])
                        args = append(args, byteArray)
                    case 20:
                        var byteArray [20]byte
                        copy(byteArray[:], bytes[:20])
                        args = append(args, byteArray)
                    case 21:
                        var byteArray [21]byte
                        copy(byteArray[:], bytes[:21])
                        args = append(args, byteArray)
                    case 22:
                        var byteArray [22]byte
                        copy(byteArray[:], bytes[:22])
                        args = append(args, byteArray)
                    case 23:
                        var byteArray [23]byte
                        copy(byteArray[:], bytes[:23])
                        args = append(args, byteArray)
                    case 24:
                        var byteArray [24]byte
                        copy(byteArray[:], bytes[:24])
                        args = append(args, byteArray)
                    case 25:
                        var byteArray [25]byte
                        copy(byteArray[:], bytes[:25])
                        args = append(args, byteArray)
                    case 26:
                        var byteArray [26]byte
                        copy(byteArray[:], bytes[:26])
                        args = append(args, byteArray)
                    case 27:
                        var byteArray [27]byte
                        copy(byteArray[:], bytes[:27])
                        args = append(args, byteArray)
                    case 28:
                        var byteArray [28]byte
                        copy(byteArray[:], bytes[:28])
                        args = append(args, byteArray)
                    case 29:
                        var byteArray [29]byte
                        copy(byteArray[:], bytes[:29])
                        args = append(args, byteArray)
                    case 30:
                        var byteArray [30]byte
                        copy(byteArray[:], bytes[:30])
                        args = append(args, byteArray)
                    case 31:
                        var byteArray [31]byte
                        copy(byteArray[:], bytes[:31])
                        args = append(args, byteArray)
                    case 32:
                        var byteArray [32]byte
                        copy(byteArray[:], bytes[:32])
                        args = append(args, byteArray)
                    default:
                        errStr := fmt.Sprintf("Error parsing fixed size byte array. Array of size %i incompatible", argument.Type.Size)
                        return nil, errors.New(errStr)
                    }
                    continue
                }

            }

            array := strings.Split(input, ",")
            argSize := argument.Type.Size
            size := len(array)
            if argSize != 0 {
                for size != argSize {
                    c.Printf("Please enter %i comma-separated list of elements:\n", argSize)
                    input = c.ReadLine()
                    array = strings.Split(input, ",")
                    size = len(array)
                }
            }

            size = len(array)

            elementType := argument.Type.Elem

            // Elements cannot be kind slice                                        only mean slice
            if elementType.Kind == reflect.Array && elementType.Type != reflect.TypeOf(common.Address{}) {
                // Is 2D byte array
                /* Nightmare to implement, have to account for:
                    * Slice of fixed byte arrays; bytes32[] in solidity for example, generally bytesn[]
                    * Fixed array of fixed byte arrays; bytes32[10] in solidity for example bytesn[m]
                    * Slice or fixed array of string; identical to above two cases as string in solidity is array of bytes

                    Since the upper bound of elements in an array in solidity is 2^256-1, and each fixed byte array
                    has a limit of bytes32 (bytes1, bytes2, ..., bytes31, bytes32), and Golang array creation takes
                    constant length values, we would have to paste the switch-case containing 1-32 fixed byte arrays
                    2^256-1 times to handle every possibility. Since arrays of arrays in seldom used, we have not
                    implemented it.
                */

                return nil, errors.New("2D Arrays unsupported. Use \"bytes\" instead.")


                /*
                slice := make([]interface{}, 0, size)
                err = addFixedByteArrays(array, elementType.Size, slice)
                if err != nil {
                    return nil, err
                }
                args = append(args, slice)
                continue
                */
            } else {
                switch elementType.Type {
                case reflect.TypeOf(bool(false)):
                    convertedArray := make([]bool, 0, size)
                    for _, item := range array {
                        b, err := utils.ConvertToBool(item)
                        if err != nil {
                            return nil, err
                        }
                        convertedArray = append(convertedArray, b)
                    }
                    args = append(args, convertedArray)
                case reflect.TypeOf(int8(0)):
                    convertedArray := make([]int8, 0, size)
                    for _, item := range array {
                        i, err := strconv.ParseInt(item, 10, 8)
                        if err != nil {
                            return nil, err
                        }
                        convertedArray = append(convertedArray, int8(i))
                    }
                    args = append(args, convertedArray)
                case reflect.TypeOf(int16(0)):
                    convertedArray := make([]int16, 0, size)
                    for _, item := range array {
                        i, err := strconv.ParseInt(item, 10, 16)
                        if err != nil {
                            return nil, err
                        }
                        convertedArray = append(convertedArray, int16(i))
                    }
                    args = append(args, convertedArray)
                case reflect.TypeOf(int32(0)):
                    convertedArray := make([]int32, 0, size)
                    for _, item := range array {
                        i, err := strconv.ParseInt(item, 10, 32)
                        if err != nil {
                            return nil, err
                        }
                        convertedArray = append(convertedArray, int32(i))
                    }
                    args = append(args, convertedArray)
                case reflect.TypeOf(int64(0)):
                    convertedArray := make([]int64, 0, size)
                    for _, item := range array {
                        i, err := strconv.ParseInt(item, 10, 64)
                        if err != nil {
                            return nil, err
                        }
                        convertedArray = append(convertedArray, int64(i))
                    }
                    args = append(args, convertedArray)
                case reflect.TypeOf(uint8(0)):
                    convertedArray := make([]uint8, 0, size)
                    for _, item := range array {
                        u, err := strconv.ParseUint(item, 10, 8)
                        if err != nil {
                            return nil, err
                        }
                        convertedArray = append(convertedArray, uint8(u))
                    }
                    args = append(args, convertedArray)
                case reflect.TypeOf(uint16(0)):
                    convertedArray := make([]uint16, 0, size)
                    for _, item := range array {
                        u, err := strconv.ParseUint(item, 10, 16)
                        if err != nil {
                            return nil, err
                        }
                        convertedArray = append(convertedArray, uint16(u))
                    }
                    args = append(args, convertedArray)
                case reflect.TypeOf(uint32(0)):
                    convertedArray := make([]uint32, 0, size)
                    for _, item := range array {
                        u, err := strconv.ParseUint(item, 10, 32)
                        if err != nil {
                            return nil, err
                        }
                        convertedArray = append(convertedArray, uint32(u))
                    }
                    args = append(args, convertedArray)
                case reflect.TypeOf(uint64(0)):
                    convertedArray := make([]uint64, 0, size)
                    for _, item := range array {
                        u, err := strconv.ParseUint(item, 10, 64)
                        if err != nil {
                            return nil, err
                        }
                        convertedArray = append(convertedArray, uint64(u))
                    }
                    args = append(args, convertedArray)
                case reflect.TypeOf(&big.Int{}):
                    convertedArray := make([]*big.Int, 0, size)
                    for _, item := range array {
                        newInt := new(big.Int)
                        newInt, ok := newInt.SetString(item, 10)
                        if !ok {
                            return nil, errors.New("Could not convert string to big.int")
                        }
                        convertedArray = append(convertedArray, newInt)
                    }
                    args = append(args, convertedArray)
                case reflect.TypeOf(common.Address{}):
                    convertedArray := make([]common.Address, 0, size)
                    for _, item := range array {
                        a := common.HexToAddress(item)
                        convertedArray = append(convertedArray, a)
                    }
                    args = append(args, convertedArray)
                default:
                    errStr := fmt.Sprintf("Type %s not found", elementType.Type)
                    return nil, errors.New(errStr)
                }
            }
        } else {
            switch argument.Type.Kind {
            case reflect.String:
                args = append(args, input)
            case reflect.Bool:
                b, err := utils.ConvertToBool(input)
                if err != nil {
                    return nil, err
                }
                args = append(args, b)
            case reflect.Int8:
                i, err := strconv.ParseInt(input, 10, 8)
                if err != nil {
                    return nil, err
                }
                args = append(args, int8(i))
            case reflect.Int16:
                i, err := strconv.ParseInt(input, 10, 16)
                if err != nil {
                    return nil, err
                }
                args = append(args, int16(i))
            case reflect.Int32:
                i, err := strconv.ParseInt(input, 10, 32)
                if err != nil {
                    return nil, err
                }
                args = append(args, int32(i))
            case reflect.Int64:
                i, err := strconv.ParseInt(input, 10, 64)
                if err != nil {
                    return nil, err
                }
                args = append(args, int64(i))
            case reflect.Uint8:
                u, err := strconv.ParseUint(input, 10, 8)
                if err != nil {
                    return nil, err
                }
                args = append(args, uint8(u))
            case reflect.Uint16:
                u, err := strconv.ParseUint(input, 10, 16)
                if err != nil {
                    return nil, err
                }
                args = append(args, uint16(u))
            case reflect.Uint32:
                u, err := strconv.ParseUint(input, 10, 32)
                if err != nil {
                    return nil, err
                }
                args = append(args, uint32(u))
            case reflect.Uint64:
                u, err := strconv.ParseUint(input, 10, 64)
                if err != nil {
                    return nil, err
                }
                args = append(args, uint64(u))
            case reflect.Ptr:
                newInt := new(big.Int)
                newInt, ok := newInt.SetString(input, 10)
                if !ok {
                    return nil, errors.New("Could not convert string to big.int")
                }
                if err != nil {
                    return nil, err
                }
                args = append(args, newInt)
            case reflect.Array:
                if argument.Type.Type == reflect.TypeOf(common.Address{}) {
                    address := common.HexToAddress(input)
                    args = append(args, address)
                } else {
                    return nil, errors.New("Conversion failed. Item is array type, cannot parse")
                }
            default:
                errStr := fmt.Sprintf("Error, type not found: %s", argument.Type.Kind)
                return nil, errors.New(errStr)
            }
        }
    }

    return
}

func checkClientExists(client *EthClient) bool {
    return client != nil
}

func addContractInstance(pathToContract string, contractName string, contracts map[string]*contract.ContractInstance) (error) {
    compiledContract, err := contract.CompileContractAt(pathToContract)
    if err != nil {
        return err
    }
    _, Abi := contract.GetContractBytecodeAndABI(compiledContract)
    abistruct, err := abi.JSON(strings.NewReader(Abi))
    if err != nil {
        return err
    }
    contracts[contractName] = &contract.ContractInstance{Contract: compiledContract, Abi: &abistruct}
    return nil
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