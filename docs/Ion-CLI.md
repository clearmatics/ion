# Ion Command Line Interface
The Ion CLI is a tool which allows users to easily interact with the Ion project. Written in golang it allows rapid development of new commands and contracts by leveraging the [ishell](https://github.com/abiosoft/ishell) and [go-ethereum smart contract bindings](https://github.com/ethereum/go-ethereum/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts).

***Note:*** The Ion CLI is not a trusted part of the Ion infrastructure rather it is just a tool to facilitate users, who should verify it functionality prior to using any unknown code.

##  Running Ion CLI
In order to compile the Ion CLI run:
```
$ cd /path/to/validation/src
$ make build
$ make test
```

Given all tests pass, the Ion CLI can be run. Prior to running the user must ensure that the `setup.json` file has been modified to contain:

    * Address and port of foreign Clique chain rpc
    * Address and port of native chain rpc
    * User account on foreign Clique chain
    * User account on native chain
    * Address of the deployed validation contract on native

Once this has been setup correctly the CLI can be launched as follows:
```
$ ./ion-cli -config [/path/to/setup.json]
===============================================================
Ion Command Line Interface

RPC Client [to]:
Listening on: 127.0.0.1:8501
User Account: 0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d
Ion Contract: 0xb9fd43a71c076f02d1dbbf473c389f0eacec559f

RPC Client [from]: 
Listening on: 127.0.0.1:8545
User Account: 0x8671e5e08d74f338ee1c462340842346d797afd3
===============================================================
>>>
```

Running help displays the available commands:
```
>>> help

Commands:
  clear                      clear the screen
  exit                       exit the program
  getBlock                   use: getBlock [integer] 
                             description: Returns block header specified
  getValidators              use: getValidators 
                             description: Returns the whitelist of validators from validator contract
  help                       display help
  latestBlock                use: latestBlock  
                             description: Returns number of latest block mined/sealed
  latestValidationBlock      use: latestValidationBlock 
                             description: Returns hash of the last block submitted to the validation contract
  submitValidationBlock      use: submitValidationBlock [integer] 
                             description: Returns the RLP block header, signed block prefix, extra data prefix and submits to validation contract
```

### Tutorial


## Extending the Ion CLI
In order to add your contract to the Ion CLI first a golang version of the solidity smart contract needs to be created, to do this we follow the instructions from [go-ethereum smart contract bindings](https://github.com/ethereum/go-ethereum/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts).

We will add a contract called `Spoon.sol` to the CLI. This requires the generation of the `abi` and `bin` files. To do this run:
```
$ npm run genbin
$ npm run genabi
```
Now the latest versions of the `abi` and `bin` files will be found in the `/path/to/ion/contracts/` directory. Next generate the `.go` version of the desired smart contract using the `abigen` to do so run:
```
$ abigen --bin=/path/to/Spoon.bin --abi /path/to/Spoon.abi --pkg contract --type Spoon --out Spoon.go
```
next place the output `Spoon.go` in the package specific directory for your golang code. The contract can then be interfaced with simply through importing the contract package.

### Golang Smart Contract Interface
Given the exisiting Ion CLI framework any additional contracts should be placed in the `ion/ion-cli/contracts/` directory and appended to the contract package.

To use an instance of the Spoon contract insert:
```
func InitSpoonContract(setup Setup, client *ethclient.Client) (Spoon *contract.Spoon) {
	// Initialise the contract
	address := common.HexToAddress(setup.Ion)
	Spoon, err := contract.NewSpoon(address, client)
	if err != nil {
		log.Fatal(err)
	}

	return
}
```
