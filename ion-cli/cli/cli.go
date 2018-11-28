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
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/abiosoft/ishell"

	"github.com/clearmatics/ion/ion-cli/config"
	contract "github.com/clearmatics/ion/ion-cli/contracts"
	"github.com/clearmatics/ion/ion-cli/utils"
)

// Launch - definition of commands and creates the interface
func Launch() {
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
			    client, err := getClient(c.Args[0])
			    if err != nil {
			        c.Println("Could not connect to client.\n")
			        return
			    }
			    ethClient = client
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
                c.Println("Usage: \taddContractInstance [name] [path/to/solidity/contract]\n")
			} else {
                err := addContractInstance(c.Args[1], c.Args[0], contracts)
                if err != nil {
                    c.Println(err)
                    return
                }
                c.Println("Added!")
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
		Help: "use: \tdeployContract [contract name] [account name] [gas limit]\n\t\t\t\tdescription: Deploys specified contract instance to connected client",
		Func: func(c *ishell.Context) {
            if len(c.Args) != 3 {
                c.Println("Usage: \tdeployContract [contract name] [account name] [gas limit] \n")
            } else {
			    if ethClient == nil {
			        c.Println("Please connect to a Client before invoking this function.\nUse \tconnectToClient [rpc url] \n")
			        return
			    }
                contractInstance := contracts[c.Args[0]]
                if contractInstance == nil {
                    errStr := fmt.Sprintf("Contract instance %s not found.\nUse \taddContractInstance [name] [path/to/solidity/contract]\n", c.Args[0])
                    c.Println(errStr)
			        return
                }

                binStr, abiStr := contract.GetContractBytecodeAndABI(contractInstance.Contract)

                account := accounts[c.Args[1]]
                if account == nil {
                    errStr := fmt.Sprintf("Account %s not found.\nUse \taddAccount [name] [path/to/keystore] \n", c.Args[1])
			        c.Println(errStr)
			        return
                }

                gasLimit, err := strconv.ParseUint(c.Args[2], 10, 64)
                if err != nil {
                    c.Println(err)
                    return
                }

                constructorInputs, err := parseMethodParameters(c, contractInstance.Abi, "")
                if err != nil {
                    c.Printf("Error parsing constructor parameters: %s\n", err)
                    return
                }

                /*gasLimit = gasLimit
                constructorInputs = constructorInputs

                c.Printf("Contract Info: %s\n\n", contractInstance.Contract.Info)
                contractInfo := make(map[string]*compiler.ContractInfo)
                str, err := json.Marshal(contractInstance.Contract.Info)
                if err != nil {
                    c.Println(err)
                    return
                }
                err = json.Unmarshal([]byte(str), &contractInfo)

                if err != nil {
                    c.Println(err)
                    return
                }
                c.Printf("Unmarshalled: %+v\n\n", contractInfo)*/


                payload := contract.CompilePayload(binStr, abiStr, constructorInputs...)

                tx, err := contract.DeployContract(
                    ctx,
                    ethClient.client,
                    account.Key.PrivateKey,
                    payload,
                    nil,
                    gasLimit,
                )
                if err != nil {
                    c.Println(err)
                    return
                }

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
		Name: "linkAndDeployContract",
		Help: "use: \tdeployContract [contract name] [account name] [gas limit]\n\t\t\t\tdescription: Deploys specified contract instance to connected client",
		Func: func(c *ishell.Context) {
            if len(c.Args) != 3 {
                c.Println("Usage: \tdeployContract [contract name] [account name] [gas limit] \n")
            } else {
			    if ethClient == nil {
			        c.Println("Please connect to a Client before invoking this function.\nUse \tconnectToClient [rpc url] \n")
			        return
			    }
                contractInstance := contracts[c.Args[0]]
                if contractInstance == nil {
                    errStr := fmt.Sprintf("Contract instance %s not found.\nUse \taddContractInstance [name] [path/to/solidity/contract] \n", c.Args[0])
                    c.Println(errStr)
			        return
                }

                c.ShowPrompt(false)
                defer c.ShowPrompt(true)
                c.Println("Please provide comma separated list of libraries to link in the form <LibraryName>:<DeployedAddress> e.g. RLP:0x123456789")
                input := c.ReadLine()
                libraries := strings.Split(input, ",")
                library := make(map[string]common.Address)

                for _, lib := range libraries {
                    name := strings.Split(lib, ":")[0]
                    address := common.HexToAddress(strings.Split(lib, ":")[1])
                    library[name] = address
                }

                compiledContract, err := contract.CompileContractWithLibraries(contractInstance.Path, library)
                if err != nil {
                    c.Println(err)
                    return
                }

                binStr, abiStr := contract.GetContractBytecodeAndABI(compiledContract)

                account := accounts[c.Args[1]]
                if account == nil {
                    errStr := fmt.Sprintf("Account %s not found.\nUse \taddAccount [name] [path/to/keystore] \n", c.Args[1])
			        c.Println(errStr)
			        return
                }

                gasLimit, err := strconv.ParseUint(c.Args[2], 10, 64)
                if err != nil {
                    c.Println(err)
                    return
                }

                constructorInputs, err := parseMethodParameters(c, contractInstance.Abi, "")
                if err != nil {
                    c.Printf("Error parsing constructor parameters: %s\n", err)
                    return
                }

                payload := contract.CompilePayload(binStr, abiStr, constructorInputs...)

                tx, err := contract.DeployContract(
                    ctx,
                    ethClient.client,
                    account.Key.PrivateKey,
                    payload,
                    nil,
                    gasLimit,
                )
                if err != nil {
                    c.Println(err)
                    return
                }

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
		Name: "transactionMessage",
		Help: "use: \ttransactionMessage [contract name] [function name] [from account name] [deployed contract address] [amount] [gasLimit] \n\t\t\t\tdescription: Calls a contract function as a transaction.",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 6 {
                c.Println("Usage: \ttransactionMessage [contract name] [function name] [from account name] [deployed contract address] [amount] [gasLimit] \n")
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
                    errStr := fmt.Sprintf("Contract instance %s not found.\nUse \taddContractInstance [name] [path/to/solidity/contract] \n", c.Args[0])
                    c.Println(errStr)
			        return
                }
                if account == nil {
                    errStr := fmt.Sprintf("Account %s not found.\nUse \taddAccount [name] [path/to/keystore]\n", c.Args[2])
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

	/*shell.AddCmd(&ishell.Cmd{
		Name: "callMessage",
		Help: "use: \tcallMessage [contract name] [function name] [from account name] [deployed contract address] \n\t\t\t\tdescription: Connects to an RPC client to be used",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 4 {
                c.Println("Usage: \tcallMessage [contract name] [function name] [from account name] [deployed contract address] \n")
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
                    errStr := fmt.Sprintf("Account %s not found.\nUse \taddAccount [name] [path/to/keystore]\n", c.Args[2])
			        c.Println(errStr)
			        return
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

                var out interface{}

                out, err = contract.CallContract(
                    ctx,
                    ethClient.client,
                    instance.Contract,
                    account.Key.Address,
                    contractDeployedAddress,
                    c.Args[1],
                    out,
                    inputs...
                )
                 if err != nil {
                    c.Println(err)
                    return
                 } else {
                    c.Printf("Result: %s\n", out)
                 }
			}
			c.Println("===============================================================")
		},
	})*/

	shell.AddCmd(&ishell.Cmd{
		Name: "getTransactionByHash",
		Help: "use: \tgetTransactionByHash [optional rpc url] [hash]\n\t\t\t\tdescription: Returns transaction specified by hash from connected client or specified endpoint",
		Func: func(c *ishell.Context) {
		    var json []byte
		    var err error

            if len(c.Args) == 1 {
                if ethClient != nil {
			        _, json, err = getTransactionByHash(ethClient, c.Args[0])
                } else {
			        c.Println("Please connect to a Client before invoking this function.\nUse \tconnectToClient [rpc url] \n")
			        return
                }
            } else if len(c.Args) == 2 {
                client, err := getClient(c.Args[0])
                if err != nil {
                    c.Println(err)
                    return
                }
			    _, json, err = getTransactionByHash(client, c.Args[1])
            } else {
                c.Println("Usage: \tgetTransactionByHash [optional rpc url] [hash]\n")
                return
            }
            if err != nil {
                c.Println(err)
                return
            }
            c.Printf("Transaction: %s\n", json)
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getBlockByNumber",
		Help: "use: \tgetBlockByNumber [optional rpc url] [integer]\n\t\t\t\tdescription: Returns block header specified by height from connected client or from specified endpoint",
		Func: func(c *ishell.Context) {
		    var json []byte
		    var err error

            if len(c.Args) == 1 {
                if ethClient != nil {
			        _, json, err = getBlockByNumber(ethClient, c.Args[0])
                } else {
			        c.Println("Please connect to a Client before invoking this function.\nUse \tconnectToClient [rpc url] \n")
			        return
                }
            } else if len(c.Args) == 2 {
                client, err := getClient(c.Args[0])
                if err != nil {
                    c.Println(err)
                    return
                }
			    _, json, err = getBlockByNumber(client, c.Args[1])
            } else {
                c.Println("Usage: \tgetBlockByNumber [optional rpc url] [integer]\n")
                return
            }
            if err != nil {
                c.Println(err)
                return
            }
            c.Printf("Block: %s\n", json)
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getBlockByHash",
		Help: "use: \tgetBlockByHash [optional rpc url] [hash] \n\t\t\t\tdescription: Returns block header specified by hash from connected client or from specific endpoint",
		Func: func(c *ishell.Context) {
		    var json []byte
		    var err error

            if len(c.Args) == 1 {
                if ethClient != nil {
			        _, json, err = getBlockByHash(ethClient, c.Args[0])
                } else {
			        c.Println("Please connect to a Client before invoking this function.\nUse \tconnectToClient [rpc url] \n")
			        return
                }
            } else if len(c.Args) == 2 {
                client, err := getClient(c.Args[0])
                if err != nil {
                    c.Println(err)
                    return
                }
			    _, json, err = getBlockByHash(client, c.Args[1])
            } else {
                c.Println("Usage: \tgetBlockByHash [optional rpc url] [hash] \n")
                return
            }
            if err != nil {
                c.Println(err)
                return
            }
            c.Printf("Block: %s\n", json)
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
                client, err := getClient(c.Args[0])
                if err != nil {
                    c.Println(err)
                    return
                }
                getProof(client, c.Args[1])
            } else {
                c.Println("Usage: \tgetProof [optional rpc url] [Transaction hash] \n")
                return
            }
            c.Println("===============================================================")
        },
    })

	//---------------------------------------------------------------------------------------------
	// 	Clique Specific Commands
	//---------------------------------------------------------------------------------------------

	shell.AddCmd(&ishell.Cmd{
		Name: "getBlockByNumber_Clique",
		Help: "use: \tgetBlockByNumber_Clique [optional rpc url] [integer]\n\t\t\t\tdescription: Returns signed and unsigned RLP-encoded block headers by block number required for submission to Clique validation from connected client or specified endpoint",
		Func: func(c *ishell.Context) {
            if len(c.Args) == 1 {
                if ethClient != nil {
			        block, _, err := getBlockByNumber(ethClient, c.Args[0])
                    if err != nil {
                        c.Println(err)
                        return
                    }
                    signedBlock, unsignedBlock := RlpEncode(block)
                    c.Printf("Signed Block: %+x\n", signedBlock)
                    c.Printf("Unsigned Block: %+x\n", unsignedBlock)
                } else {
			        c.Println("Please connect to a Client before invoking this function.\nUse \tconnectToClient [rpc url] \n")
			        return
                }
            } else if len(c.Args) == 2 {
                client, err := getClient(c.Args[0])
                if err != nil {
                    c.Println(err)
                    return
                }
			    block, _, err := getBlockByNumber(client, c.Args[1])
                if err != nil {
                    c.Println(err)
                    return
                }
                signedBlock, unsignedBlock := RlpEncode(block)
                c.Printf("Signed Block:\n %+x\n", signedBlock)
                c.Printf("Unsigned Block:\n %+x\n", unsignedBlock)
            } else {
                c.Println("Usage: \tgetBlockByNumber_Clique [optional rpc url] [integer]\n")
                return
            }
			c.Println("===============================================================")
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "getBlockByHash_Clique",
		Help: "use: \tgetBlockByHash_Clique [optional rpc url] [hash] \n\t\t\t\tdescription: Returns signed and unsigned RLP-encoded block headers by block hash required for submission to Clique validation from connected client or specified endpoint",
		Func: func(c *ishell.Context) {
            if len(c.Args) == 1 {
                if ethClient != nil {
			        block, _, err := getBlockByHash(ethClient, c.Args[0])
                    if err != nil {
                        c.Println(err)
                        return
                    }
                    signedBlock, unsignedBlock := RlpEncode(block)
                    c.Printf("Signed Block: 0x%+x\n", signedBlock)
                    c.Printf("Unsigned Block: 0x%+x\n", unsignedBlock)
                } else {
			        c.Println("Please connect to a Client before invoking this function.\nUse \tconnectToClient [rpc url] \n")
			        return
                }
            } else if len(c.Args) == 2 {
                client, err := getClient(c.Args[0])
                if err != nil {
                    c.Println(err)
                    return
                }
			    block, _, err := getBlockByHash(client, c.Args[1])
                if err != nil {
                    c.Println(err)
                    return
                }
                signedBlock, unsignedBlock := RlpEncode(block)
                c.Printf("Signed Block:\n %+x\n", signedBlock)
                c.Printf("Unsigned Block:\n %+x\n", unsignedBlock)
            } else {
                c.Println("Usage: \tgetBlockByHash_Clique [optional rpc url] [hash]\n")
                return
            }
			c.Println("===============================================================")
		},
	})

	shell.Run()
}

func getClient(url string) (*EthClient, error) {
    rpc := utils.ClientRPC(url)
    eth := ethclient.NewClient(rpc)

    client := EthClient{client: eth, rpcClient: rpc, url: url}

    _, _, err := getBlockByNumber(&client, "0")

    return &client, err
}

func parseMethodParameters(c *ishell.Context, abiStruct *abi.ABI, methodName string) (args []interface{}, err error) {
    var inputParameters abi.Arguments
    if methodName != "" {
        inputParameters = abiStruct.Methods[methodName].Inputs
    } else {
        inputParameters = abiStruct.Constructor.Inputs
    }

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
    fmt.Println("Compiling contract...")
    compiledContract, err := contract.CompileContractAt(pathToContract)
    if err != nil {
        return err
    }
    _, Abi := contract.GetContractBytecodeAndABI(compiledContract)
    abistruct, err := abi.JSON(strings.NewReader(Abi))
    if err != nil {
        return err
    }
    fmt.Println("Creating contract instance...")
    contracts[contractName] = &contract.ContractInstance{Contract: compiledContract, Abi: &abistruct, Path: pathToContract}
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