// Copyright (c) 2018 Clearmatics Technologies Ltd
package cli

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/abiosoft/ishell"

	"github.com/clearmatics/ion/ion-cli/config"
	contract "github.com/clearmatics/ion/ion-cli/contracts"
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
					inputs...,
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
