// Copyright (c) 2018 Clearmatics Technologies Ltd
package cli

import (
	"context"
	"fmt"
	"log"
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

	//ethclientTo := ethclient.NewClient(clientTo)
	ethclientFrom := ethclient.NewClient(clientFrom)

	var ethClient *EthClient = nil
	var contracts map[string]*contract.ContractInstance = make(map[string]*contract.ContractInstance)
	var accounts map[string]*config.Account = make(map[string]*config.Account)
    contracts = contracts
    accounts = accounts

	// Get a suggested gas price
	gasPrice, err := ethclientFrom.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Create an authorized transactor and corrsponding privateKey
	authTo, keyTo := config.InitUser(configuration.KeystoreTo, configuration.PasswordTo)
	authTo.Value = big.NewInt(0)     // in wei
	authTo.GasLimit = uint64(100000) // in units
	authTo.GasPrice = gasPrice

	// Create an authorized transactor and spend 1 unicorn
	authFrom, keyFrom := config.InitUser(configuration.KeystoreFrom, configuration.PasswordFrom)
	authFrom.Value = big.NewInt(0)     // in wei
	authFrom.GasLimit = uint64(100000) // in units
	authFrom.GasPrice = gasPrice

    keyTo = keyTo
    keyFrom = keyFrom

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
		Help: "use: \taddContractInstances [name] [path/to/solidity/contract] [deployed address] \n\t\t\t\tdescription: Compiles a contract for use",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 3 {
                c.Println("Usage: \taddContractInstances [name] [path/to/solidity/contract] [deployed address] \n")
			} else {
                compiledContract := contract.CompileContractAt(c.Args[1])
                _, Abi := contract.GetContractBytecodeAndABI(compiledContract)
                abistruct, err := abi.JSON(strings.NewReader(Abi))
                if err != nil {
                    c.Err(err)
                }
                contracts[c.Args[0]] = &contract.ContractInstance{Contract: compiledContract, Address: common.HexToAddress(c.Args[2]), Abi: &abistruct}
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
		Help: "use: \taddAccount [name] [path/to/keystore] [password] \n\t\t\t\tdescription: Add account to be used for transactions",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 3 {
                c.Println("Usage: \taddAccount [name] [path/to/keystore] [password] \n")
			} else {
                auth, key := config.InitUser(c.Args[1], c.Args[2])
                account := &config.Account{Auth: auth, Key: key}
                accounts[c.Args[0]] = account
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
		Name: "messageCallFunction",
		Help: "use: \tmessageCallFunction [contract name] [function name] [from account name] [amount] [gasLimit] \n\t\t\t\tdescription: Connects to an RPC client to be used",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 5 {
                c.Println("Usage: \tmessageCallFunction [contract name] [function name] [from account name] [amount] [gasLimit] \n")
			} else {
			    if ethClient == nil {
			        c.Println("Please connect to a Client before invoking this function.\nUse \tconnectToClient [rpc url] \n")
			        return
			    }

                instance := contracts[c.Args[0]]
                methodName := c.Args[1]
                account := accounts[c.Args[2]]

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
                amount, ok := amount.SetString(c.Args[3], 10)
                if !ok {
                    c.Err(errors.New("Please enter an integer for <amount>"))
                }
                gasLimit, err := strconv.ParseUint(c.Args[4], 10, 64)
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
                //inputs := []interface{}{[]bool{true, false}, []bool{false,false,true}}

                c.Printf("Inputs = %s\n", inputs)
                c.Printf("First index type: %s\n", reflect.TypeOf(inputs[0]))

                tx, err := contract.TransactionContract(
                    ctx,
                    ethClient.client,
                    account.Key.PrivateKey,
                    instance.Contract,
                    instance.Address,
                    amount,
                    gasLimit,
                    c.Args[1],
                    inputs...
                )
                 if err != nil {
                    c.Println(err)
                    return
                 } else {
                    c.Printf("Transaction hash: %s", tx)
                 }
			}
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getBlockByNumber",
		Help: "use: \tgetBlockByNumber [rpc url] [integer] \n\t\t\t\tdescription: Returns block header specified from chain [TO/FROM]",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
                if ethClient != nil {
			        getBlockByNumber(ethClient, c.Args[0])
                } else {
				    c.Println("Usage: \tgetBlock [rpc url] [integer] \n")
                }
			} else {
			    c.Println("Connecting to client...\n")
			    newClient := getClient(c.Args[0])
			    getBlockByNumber(newClient, c.Args[1])
			}
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getBlockByHash",
		Help: "use: \tgetBlockByNumber [rpc url] [hash] \n\t\t\t\tdescription: Returns block header specified from chain [TO/FROM]",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
                if (ethClient != nil) {
			        getBlockByHash(ethClient, c.Args[0])
                } else {
				    c.Println("Usage: \tgetBlock [rpc url] [hash] \n")
                }
			} else {
			    c.Println("Connecting to client...\n")
			    newClient := getClient(c.Args[0])
			    getBlockByHash(newClient, c.Args[1])
			}
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
        Name: "getProof",
        Help: "use: \tgetProof [rpc url] [Transaction Hash] \n\t\t\t\tdescription: Returns a merkle patricia proof of a specific transaction and its receipt in a block",
        Func: func(c *ishell.Context) {
            if len(c.Args) != 2 {
                if (ethClient != nil) {
                    getProof(ethClient, c.Args[0])
                } else {
                    c.Println("Usage: \tgetBlock [rpc url] [hash] \n")
                }
            } else {
                getProof(getClient(c.Args[0]), c.Args[1])
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

        c.Printf("Name: %s\n", argument.Name)
        c.Printf("Type: %s\n", argument.Type)
        c.Printf("Type: %s\n", argument.Type.Type)
        c.Printf("Kind: %s\n", argument.Type.Kind)
        c.Printf("Elem: %s\n", argument.Type.Elem)
        if argument.Type.Elem != nil {
            c.Printf("Elem Type: %s\n", argument.Type.Elem.Type)
            c.Printf("Elem Kind: %s\n", argument.Type.Elem.Kind)
        }

        c.Printf("Reflect byte type: %s\n", reflect.TypeOf([]byte{}))

        if argument.Type.Kind == reflect.Array || argument.Type.Kind == reflect.Slice {
            c.Println("Argument is array\n")

            // Occurs if argument is address which resolves Kind to array with no element
            if argument.Type.Type == reflect.TypeOf(common.Address{}) {
                item, err := utils.ConvertToType(input, &argument.Type)
                if err != nil {
                    c.Err(err)
                }
                args = append(args, item)
                continue
            }

            elementType := argument.Type.Elem
            // bytes = []byte{} argument type = slice, no element, type equates to []uint8
            // byte[] = [][1]byte{} argument type = slice, element type = array, type equates to [][1]uint8
            // byte = bytes1
            // bytesn = [n]byte{} 0 < n < 33, argument type = array, no element, type equates to [n]uint8
            // bytesn[] = [][n]byte{} argument type = slice, element type = array, type equares to [][n]uint8
            // bytesn[m] = [m][n]byte{} argument type = array, element type = array, type equates to [m][n]uint8
            // Many annoying cases of byte arrays

            c.Printf("Argument type: %s\n", argument.Type.Type)
            c.Printf("Byte array type: %s\n", reflect.TypeOf([]byte{}))
            if argument.Type.Type == reflect.TypeOf([]byte{}) {
                c.Printf("Element is byte\n")
                bytes, err := hex.DecodeString(input)
                if err != nil {
                    c.Err(err)
                }
                args = append(args, bytes)
            } else {
                c.Printf("Element type: %s\n", elementType.Type)
                c.Printf("Element kind: %s\n", elementType.Kind)
                c.Printf("Element size: %s\n", elementType.Size)
                c.Printf("Element element: %s\n", elementType.Elem)
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
                /*case reflect.TypeOf([]byte{}):
                    convertedArray := make([]common.Address, 0, size)
                    for _, item := range array {
                        a := common.HexToAddress(item)
                        convertedArray = append(convertedArray, a)
                    }
                    args = append(args, convertedArray)*/
                default:
                    errStr := fmt.Sprintf("Type %s not found", elementType.Type)
                    return nil, errors.New(errStr)
                }
            }
        } else {
            item, err := utils.ConvertToType(input, &argument.Type)
            if err != nil {
                c.Err(err)
            }
            args = append(args, item)
        }
    }

    return
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