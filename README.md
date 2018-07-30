# Ion Interoperability Protocol

The Ion Interoperability Protocol provides mechanisms to perform atomic swaps and currency transfers across multiple turing-complete blockchains.

## Ion State Verification Scheme

## Block Validation Scheme
Block validation scheme is a set of smart contracts which ensures that block headers submitted to the contract are mined/sealed by an approved partie(s). The motivation behind this is to update the state of a blockchain onto another blockchain. To do this we first need to know which blocks are valid - dependent on the definition of a valid by the underlying consensus algorithm. As deterministic finality is a requirement we seek Clique PoA and Istanbul PBFT consensus is to be used on the chain from which the state is being taken.

For a full description and roadmap of the project please refer to the Clearmatics [Ion-Stage-2 Wiki](https://github.com/clearmatics/ion/wiki/Ion-Stage-2---Proposal#validation).

## Running the Project
A Clique PoA chain is launched and then the block headers are taken and updated to the validation contract which is deployed on a ganache chain. However in the javascript tests the contract is deployed on the PoA chain itself for sake of simplicity, this has no bearing on the functionality of the project.

Running the project requires initialisation of the following components:
  * Validation smart contract and tests
  * Two separate blockchain: Clique _proof of authority_ network, Ganache test network
  * Golang Ion CLI used to interact with the smart contract(s)

In order to use the smart contracts and run the tests it is necessary to first initialise the test networks.

***Note:*** that as the contract searches for specific parts of the block header that only exist in Clique and IBFT, Ganache or Ethash chains cannot be used for the _root_ blockchain from which headers are extracted.

### Requirements
* golang version: 1.9.x

### Installing Ion
Having cloned and entered the repo:
```
$ git clone git@github.com:clearmatics/ion.git
$ cd /path/to/ion
```

Now run the command:
```
$ tree -L 1
```

Which hopefully returns this:
```
.
├── abi
├── CODE_OF_CONDUCT.md
├── contracts
├── CONTRIBUTING.md
├── docs
├── ion-cli
├── LICENSE
├── Makefile
├── migrations
├── package.json
├── package-lock.json
├── README.md
├── test
└── truffle.js
```

### Testing Contracts
In order to test the Solidity contracts using the Javascript tests a testrpc must be run. As the validation contract relies upon receiving signatures in the `extraData` field of the block header it is not sufficient to run an instance of ganache-cli, rather a Clique PoA chain must be initialised.

To use the tests please follow the instructions on how to run a single validator Clique chain given [here](https://github.com/maxrobot/network-geth). Additionally you must ensure that the account sealing blocks is identical to that defined in the `validation.js` test itself `0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d`.

Having launched a single-validator clique chain with the sealer `0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d`, run the tests as follows.
```
$ npm install
$ npm run test
```

### Ion Command Line Interface
The Ion CLI is a tool for handling the passing of block headers between to blockchains written in Golang to leverage the extensive ethereum libraries. It is not a critical part of the Ion infrastructure rather is just an open utility that people can extend and use as they see fit.

In its current form the Ion CLI allows the user to connect to two separate blockchains, via RPC, and submit block headers to a validation contract.

#### Testing Ion CLI
```
$ cd ion-cli
$ make build
```
In order to run the basic unit tests for the Ion CLI run,
```
$ make test
```
If the tests pass successfully then the CLI can be run.

Additional integration tests can be run however it requires launching a Clique PoA chain as above. To run the integration tests launch command,
```
$ make integration-test
```

#### Running Ion CLI
As mentioned in the project description this simple implementation of the validation contract is active only on a single blockchain, however the CLI is simulating the passing of the headers to and from as if it were between separate chains.

Having followed the instructions on how to setup a Clique blockchain, which is hosted on `127.0.0.1:8501`, we run a ganache-cli in another terminal on `127.0.0.1:8545` and deploy the contract to the ganache blockchain,
```
$ npm run testrpc
$ npm run deploy
```

Following this we can attach to the Ion Command Line Interface,
```
$ cd /path/to/ion/ion-cli
$ make build
```
Assuming a successful build we must create a setup file which contains the connection of the separate blockchains, user accounts, account keystores, and the address of the deployed validation contract. Change the default values in the example setup.json then run the `ion-cli` poiinting to the modified setup file.
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
running help displays the available commands.
```
>>> help

Commands:
  clear                      clear the screen
  exit                       exit the program
  generateTxProof            use: generateTxProof [Transaction Hash] [Block Number] 
                             description: Returns the proof of a specific transaction held within a Patricia trie
  getBlock                   use: getBlock [integer] 
                             description: Returns block header specified
  getValidationBlock         use: latestValidationBlock 
                             description: Returns hash of the last block submitted to the validation contract
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


