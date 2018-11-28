# Ion Interoperability Framework
The Ion Interoperability Framework is a library that provides an interface for the development of general cross-chain smart contracts.

## Getting Started

Clone the repository

## Interoperate with Rinkeby!

We've already deployed some contracts to the Rinkeby test network for you to play around with!

Trigger: 0xA2e4a61a3D2ce626Ba9B3e927cfFDB0e4E0bd06d

Procedure:
./ion-cli
>>> addContractInstance ion /path/to/Ion.sol
>>> addContractInstance ethstore /path/to/EthereumStore.sol
>>> addContractInstance clique /path/to/Clique.sol
>>> addAccount name_of_your_account /path/to/keystore.json
>>> connectToClient https://your.endpoint:port
>>> deployContract ion name_of_your_account gaslimit
>>> deployContract ethstore name_of_your_account gaslimit
>>> deployContract clique name_of_your_account gaslimit

## Develop on Ion


## CLI


## How it works


# Contribute!




The Ion Interoperability Framework provides a interface to perform interoperability across multiple turing-complete blockchains. Using the generalised framework specific use cases, such as atomic swaps, can be developed.

In order to perform cross-chain interop the Ion framework verifies specific transactions executed on one blockchain A on another blockchain B. Being able to verify these transactions cross-chain requires submission of valid block headers, to the blockchain where a transaction is to be proven, and then performing a number of Patricia trie proofs of the transaction, receipts and logs. Smart contracts can then be built to execute only if the proof is verified, this is known as continuous execution.

To facilitate this Ion has three main components:
  * Block Storage Contracts
  * Modular Validation Scheme
  * Ion Framework Contracts

The Ion hub contract is the core component of the framework, the contract persists key data of all valid blocks submitted to the validation contracts required to verify a state transition. To prove a transaction has occurred on an external blockchain the Ion framework contract require (for EVM-based chains):

* Block Header:
  * Tx root hash
  * Receipt root hash
* Block Hash

For other blockchains, different data will be required to prove state transitions and as such contracts must be written that adhere to those system-specific mechanisms.

For each block submitted to the validation contracts this information is appended to the Ion hub contract, thus creating a generic interface for Ion framework contracts to receive valid block data. Ion provides a generalised interoperability framework and thus block validation is designed to be modular, in order to allow for interop between chains with any consensus mechanism.

When external blockchains are registered with a specific validation contract, having its own consensus specific validation mechanism, the validation contract adds the chain to the Ion hub. Subsequently all blocks successfully submitted to the validation contract are added to the Ion hub.

More details can be found in the [Ion Wiki](https://github.com/clearmatics/ion/wiki).

## Running the Project
A Clique PoA chain is launched and then the block headers are taken and updated to the validation contract which is deployed on a ganache chain. However in the javascript tests the contract is deployed on the PoA chain itself for sake of simplicity, this has no bearing on the functionality of the project.

Running the project requires initialisation of the following components:
  * Validation smart contract and tests
  * Two separate blockchain: Clique _proof of authority_ network, Ganache test network
  * Golang Ion CLI used to interact with the smart contract(s)

In order to use the smart contracts and run the tests it is necessary to first initialise the test networks.

**Note:** that as the contract searches for specific parts of the block header that only exist in Clique and IBFT, Ganache or Ethash chains cannot be used for the _root_ blockchain from which headers are extracted.

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

Having launched a Clique PoA chain, hosted on `127.0.0.1:8501` (the first account on the `eth.accounts` array should be unlocked in the node), run the tests as follows:
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
$ make test
```
Assuming a successful build and passing of the tests a setup file which contains the connection of the separate blockchains, user accounts, account keystores, and the address of the deployed validation contract should be created. Change the default values in the example setup.json then run the `ion-cli` poiinting to the modified setup file.
```
$ ./ion-cli --config [/path/to/setup.json]
===============================================================
Ion Command Line Interface

RPC Client [to]:
	Listening on:		http://127.0.0.1:8501
	User Account:		0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d
	RPC ChainId:		0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075
	Validation Contract:	0xb9fd43a71c076f02d1dbbf473c389f0eacec559f
	Ion Contract:		0x6aa4444974f60bf3a0bf074d3c194f88ae4d4613
	Function Contract:	0x49e71afdcaf62d7384f0b801c9e3c6e18d4a2597

RPC Client [from]:
	Listening on:		https://127.0.0.1:8545
	User Account:		0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d
	Trigger Contract:	0x61621bcf02914668f8404c1f860e92fc1893f74c
===============================================================
>>>
```
running help displays the available commands.
```
>>> help
Commands:
   checkBlockValidation         use:   checkBlockValidation
                                       Enter Blockhash: [HASH]
                                description: Returns true for validated blocks
   clear                        clear the screen
   exit                         exit the program
   getBlock                     use:   getBlock [TO/FROM] [integer] 
                                description: Returns block header specified from chain [TO/FROM]
   help                         display help
   latestBlock                  use:   latestBlock [TO/FROM] 
                                description: Returns number of latest block mined/sealed from chain [TO/FROM]
   latestValidatedBlock         use:   latestValidatedBlock 
                                description: Returns hash of the last block submitted to the validation contract
   registerChainValidation      use:   registerChainValidation
                                       Enter Validators: [ADDRESS ADDRESS]
                                       Enter Genesis Hash: [HASH] 
                                description: Register new chain with validation contract
   submitBlockValidation        use:   submitBlockValidation
                                       Enter Block Number: [INTEGER]
                                description: Returns the RLP block header, signed block prefix, extra data prefix and submits to validation contract
   triggerEvent                 use:   triggerEvent 
                                description: Returns hash of the last block submitted to the validation contract
   verifyAndExecute             use:   verifyAndExecute [Transaction Hash] 
                                description: Returns the proof of a specific transaction held within a Patricia trie
```

#### Ion CLI Walkthrough
A simple walkthrough of the Ion framework between a local testrpc and the Rinkeby test network is detailed. Users can trigger an event on Rinkeby with the account `0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d` and verify the transaction on the local testrpc where they deploy the event consuming contracts.

A brief overview of the steps required are:
* Launch Ion CLI
* Register rinkeby with the validation contract
* Trigger event on rinkeby test network - this tutorial will use the transaction `0x5da684940b4fd9dec708cc159dc504aa01e90d40bb76a2b73299aee13aa72098` - [check etherscan](https://rinkeby.etherscan.io/tx/0x5da684940b4fd9dec708cc159dc504aa01e90d40bb76a2b73299aee13aa72098)
* Submit relevant rinkeby blocks to the validation contract
* Verify event happened on testrpc and execute function

##### Step 1. Launch Ion CLI
Having followed the instructions to run and build the Ion CLI and testrpc, see [here](https://github.com/clearmatics/ion/tree/ion-stage-2#running-ion-cli), launch the CLI with the setup file `rinkeby.json`.
```
$ ./ion-cli --config rinkeby.json
===============================================================
Ion Command Line Interface

RPC Client [to]:
	Listening on:		http://127.0.0.1:8545
	User Account:		0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d
	RPC ChainId:		0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075
	Validation Contract:	0xb9fd43a71c076f02d1dbbf473c389f0eacec559f
	Ion Contract:		0x6aa4444974f60bf3a0bf074d3c194f88ae4d4613
	Function Contract:	0x49e71afdcaf62d7384f0b801c9e3c6e18d4a2597

RPC Client [from]:
	Listening on:		https://rinkeby.infura.io
	User Account:		0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d
	Trigger Contract:	0x61621bcf02914668f8404c1f860e92fc1893f74c
===============================================================
>>>
```

##### Step 2. Register Chain
In order to validate blocks first the chain must be registered in the validation contract on the testrpc:
```
>>> registerChainValidation
Connecting to: http://127.0.0.1:8545
Enter Validators: 0x42eb768f2244c8811c63729a21a3569731535f06 0x6635f83421bf059cd8111f180f0727128685bae4 0x7ffc57839b00206d1ad20c69a1981b489f772031 0xb279182d99e65703f0076e4812653aab85fca0f0 0xd6ae8250b8348c94847280928c79fb3b63ca453e 0xda35dee8eddeaa556e4c26268463e26fb91ff74f 0xfc18cbc391de84dbd87db83b20935d3e89f5dd91
Enter Genesis Hash: 0x100dc525cdcb7933e09f10d4019c38d342253a0aa32889fbbdbc5f2406c7546c

Returns:
Transaction: 0xcd5f4405260a6935b048e9136d211df99e594359abe04dcf975c730e5cf0d708
===============================================================
```

To see if it has successfully been registered we can check the the contract if this is a valid block:
```
>>> checkBlockValidation
Connecting to: http://127.0.0.1:8545
Enter BlockHash: 0x100dc525cdcb7933e09f10d4019c38d342253a0aa32889fbbdbc5f2406c7546c
Checking for valid block:
ChainId:	ab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075
BlockHash:	100dc525cdcb7933e09f10d4019c38d342253a0aa32889fbbdbc5f2406c7546c

Returns:
Valid:		true
===============================================================
```

##### Step 3. Submit Block to Validation
As our event happened in block 2776659 this block must be submitted to the validation contract prior to verification of our transaction.
```
>>> submitBlockValidation
Connecting to: http://127.0.0.1:8545
Enter Block Number: 2776659
RLP encode block:
Number:		2776659
Signed Block Header Prefix:
f9025ca0100dc525cdcb7933e09f10d4019c38d342253a0aa32889fbbdbc5f2406c7546ca01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347940000000000000000000000000000000000000000a0ad9b6d8c20a0631e2513968cdf3667dffabf9d2f6c1bf22a5990861192e1d266a053413d0e5fd5854665fab663ad8ffb0f5d06bf1907a5a0a3e45300de1ce23fcda0997750d465d96422e7692e281d16db709e6ff9c11c1d8229410340cf1598a8b6b9010000000000c0000010000000000000000000000000000000008000000000000000000000040000000000002004000000000000000002000000000000000004000000010000000008000000002800000002000000000004000000000001000001400000000002000000040000000000080080000000000000000000001000008000000000002000000000000000000000000000000000000000000002000000000000000000010000000010004000000000000020000000000000000000000000000200010300000000000005000040000000000004000000000000000042006000000000200000008000000000000000800000020400000000001080020080000001832a5e53836b33668349d003845b6ab5b9b861d68301080d846765746886676f312e3130856c696e75780000000000000000000cd4835e13d2204ad1fbc94d18a1a0373e92ffc63a021f8212a72dd1cff2e2057bded7534ed90df2492e60c7050db528715e3bc7734bef0a4c97542886cbc57c01a00000000000000000000000000000000000000000000000000000000000000000880000000000000000

Unsigned Block Header Prefix:
f9021aa0100dc525cdcb7933e09f10d4019c38d342253a0aa32889fbbdbc5f2406c7546ca01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347940000000000000000000000000000000000000000a0ad9b6d8c20a0631e2513968cdf3667dffabf9d2f6c1bf22a5990861192e1d266a053413d0e5fd5854665fab663ad8ffb0f5d06bf1907a5a0a3e45300de1ce23fcda0997750d465d96422e7692e281d16db709e6ff9c11c1d8229410340cf1598a8b6b9010000000000c0000010000000000000000000000000000000008000000000000000000000040000000000002004000000000000000002000000000000000004000000010000000008000000002800000002000000000004000000000001000001400000000002000000040000000000080080000000000000000000001000008000000000002000000000000000000000000000000000000000000002000000000000000000010000000010004000000000000020000000000000000000000000000200010300000000000005000040000000000004000000000000000042006000000000200000008000000000000000800000020400000000001080020080000001832a5e53836b33668349d003845b6ab5b9a0d68301080d846765746886676f312e3130856c696e7578000000000000000000a00000000000000000000000000000000000000000000000000000000000000000880000000000000000

Returns:
Transaction Hash: 0x965da7ac02d7ef58417e32b427b33ff2717d3cd6bb213fd05853a2c8d0fd092e
===============================================================
```

Again check that the block has successfully been validated:
```
>>> checkBlockValidation
Connecting to: http://127.0.0.1:8545
Enter BlockHash: 0x74d37aa3c96bc98903451d0baf051b87550191aa0d92032f7406a4984610b046
Checking for valid block:
ChainId:	ab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075
BlockHash:	74d37aa3c96bc98903451d0baf051b87550191aa0d92032f7406a4984610b046

Returns:
Valid:		true
===============================================================
```

##### Step 4. Verify and Execute Event Consuming Contract
Now the block that contains our event has been successfully submitted to the validation contract we can execute the contract which consumes the event:
```
>>> verifyAndExecute
Connecting to: http://127.0.0.1:8545
Enter Transaction Hash: 0x5da684940b4fd9dec708cc159dc504aa01e90d40bb76a2b73299aee13aa72098
Enter Block Hash: 0x74d37aa3c96bc98903451d0baf051b87550191aa0d92032f7406a4984610b046

Returns:
Transaction Hash: 0x9da94ab127a05a81c5a8ec159e5e103efe44403c57685d53da1b0126880a47fb
===============================================================
```

##### Step 5. Check Transaction Successfully Executed
Attach to the geth client of the RPC TO chain and run:
```
> eth.getTransactionReceipt("0x40538002f640c647bcf4feb922e31523c0cd2b46b38791dd7368e8cbe5bbba15")
{
  blockHash: "0x9da94ab127a05a81c5a8ec159e5e103efe44403c57685d53da1b0126880a47fb",
  blockNumber: 11,
  contractAddress: null,
  cumulativeGasUsed: 352099,
  gasUsed: 352099,
  logs: [{
      address: "0x9abefbe4cca994c5d1934dff50c6a863edcf5f52",
      blockHash: "0x9da94ab127a05a81c5a8ec159e5e103efe44403c57685d53da1b0126880a47fb",
      blockNumber: 11,
      data: "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac807574d37aa3c96bc98903451d0baf051b87550191aa0d92032f7406a4984610b0460000000000000000000000000000000000000000000000000000000000000002",
      logIndex: 0,
      topics: ["0xf0bc00f5b90f382e1bbca216713ca9e2e8e298f9d7717d30847905395f287046"],
      transactionHash: "0x40538002f640c647bcf4feb922e31523c0cd2b46b38791dd7368e8cbe5bbba15",
      transactionIndex: 0,
      type: "mined"
  }, {
      address: "0x9abefbe4cca994c5d1934dff50c6a863edcf5f52",
      blockHash: "0x9da94ab127a05a81c5a8ec159e5e103efe44403c57685d53da1b0126880a47fb",
      blockNumber: 11,
      data: "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac807574d37aa3c96bc98903451d0baf051b87550191aa0d92032f7406a4984610b0460000000000000000000000000000000000000000000000000000000000000000",
      logIndex: 1,
      topics: ["0xf0bc00f5b90f382e1bbca216713ca9e2e8e298f9d7717d30847905395f287046"],
      transactionHash: "0x40538002f640c647bcf4feb922e31523c0cd2b46b38791dd7368e8cbe5bbba15",
      transactionIndex: 0,
      type: "mined"
  }, {
      address: "0x9abefbe4cca994c5d1934dff50c6a863edcf5f52",
      blockHash: "0x9da94ab127a05a81c5a8ec159e5e103efe44403c57685d53da1b0126880a47fb",
      blockNumber: 11,
      data: "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac807574d37aa3c96bc98903451d0baf051b87550191aa0d92032f7406a4984610b0460000000000000000000000000000000000000000000000000000000000000001",
      logIndex: 2,
      topics: ["0xf0bc00f5b90f382e1bbca216713ca9e2e8e298f9d7717d30847905395f287046"],
      transactionHash: "0x40538002f640c647bcf4feb922e31523c0cd2b46b38791dd7368e8cbe5bbba15",
      transactionIndex: 0,
      type: "mined"
  }, {
      address: "0x93981af8db02c7ef40d0ed61caef2726a79eb903",
      blockHash: "0x9da94ab127a05a81c5a8ec159e5e103efe44403c57685d53da1b0126880a47fb",
      blockNumber: 11,
      data: "0x00",
      logIndex: 3,
      topics: ["0x68f46c45a243a0e9065a97649faf9a5afe1692f2679e650c2f853b9cd734cc0e"],
      transactionHash: "0x40538002f640c647bcf4feb922e31523c0cd2b46b38791dd7368e8cbe5bbba15",
      transactionIndex: 0,
      type: "mined"
  }],
  logsBloom: "0x00000000000000000000000000000000000000000001000000080000000000000000000800000000000000000000000000000000000000010000000000000000000040000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000100080000000000000000000000040000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000",
  status: "0x1",
  transactionHash: "0x40538002f640c647bcf4feb922e31523c0cd2b46b38791dd7368e8cbe5bbba15",
  transactionIndex: 0
}
``` 
`status: "0x1"` shows the transaction executed successfully!
