# Ion Interoperability Protocol

The Ion Interoperability Protocol provides mechanisms to perform atomic swaps and currency transfers
across multiple turing-complete blockchains.

## Block Validation Scheme
Block validation scheme is a set of smart contract which aim to be able to ensure that block headers submitted to the contract are sealed by an approved partie(s). The motivation behind this is to update the state of a Clique PoA or Istanbul PBFT chain onto any other blockchain. To do this we first need to know which blocks are valid - dependent on the definition of a valid by the underlying consensus algorithm.

At this stage a Clique PoA chain is launched and then the block headers updated to the smart contract which is deployed on itself. However in the near future this will be extended to update the Clique header onto any other chain. For a full description and roadmap of the project please refer to the Clearmatics [Ion-Stage-2 Wiki](https://github.com/clearmatics/ion/wiki/Ion-Stage-2---Proposal#validation).

## Running the Project
This project contains a number of components:
  * Validation smart contract and tests
  * Clique _proof of authority_ test network
  * Command Line Interface for interacting with the smart contract

In order to use the smart contracts and run the tests it is necessary to first initialise the test network.

***Note:*** that as the contract searches for specific parts of the block header that only exist in Clique, Ganache or Ethash chains cannot be used.

### Initialise the Test Network
The instructions are based on the tutorial of [Salanfe](https://hackernoon.com/setup-your-own-private-proof-of-authority-ethereum-network-with-geth-9a0a3750cda8) but has the more complicated parts already initialised.

First install an instance of [geth](https://geth.ethereum.org/downloads/) and clone the repo.

Having cloned and entered the repo:
```
$ git clone git@github.com:maxrobot/validation.git
$ cd /path/to/validation
```

Now run the command:
```
$ tree -L 1
```

Which hopefully returns this:
```
├── abi
├── build
├── contracts
├── migrations
├── package.json
├── poa-network
├── README.md
├── src
├── test
└── truffle.js
```

#### Initialise Nodes
Network files are found in the `/path/to/validation/poa-network` directory. Enter the poa-network directory and initialise the two nodes which will be sealing blocks:
```
$ cd /path/to/validation/poa-network
$ geth --datadir node1/ init genesis.json
$ geth --datadir node2/ init genesis.json
```

#### Launch the Bootnode
The boot node tells the peers how to connect with each other. In another terminal instance run:
```
$ bootnode -nodekey boot.key -verbosity 9 -addr :30310
$ INFO [06-07|12:16:21] UDP listener up                          self=enode://dcb1dbf8d710eb7d10e0e2db1e6d3370c4b048efe47c7a85c4b537b60b5c11832ef25f026915b803e928c1d93f01b853131e412c6308c4c6141d1504c78823c8@[::]:30310
```
As the peers communicate this terminal should fill with logs.

***Note:*** bootnode may not be found if go-ethereum/geth is not installed fully.

#### Start and Attach to the Nodes
Each node must be launched either as a background operation or on separate terminal instances. Thus from the poa-network directory for node 1 run:
```
$ geth --datadir node1/ --syncmode 'full' --port 30311 --rpc --rpcaddr 'localhost' --rpcport 8501 --bootnodes 'enode://dcb1dbf8d710eb7d10e0e2db1e6d3370c4b048efe47c7a85c4b537b60b5c11832ef25f026915b803e928c1d93f01b853131e412c6308c4c6141d1504c78823c8@127.0.0.1:30310' --networkid 1515 --gasprice '1' -unlock '0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d' --password node1/password.txt --mine
```
then attach:
```
$ geth attach node1/geth.ipc
```
 and again for node 2:
```
$ geth --datadir node2/ --syncmode 'full' --port 30312 --rpc --rpcaddr 'localhost' --rpcport 8502 --bootnodes 'enode://dcb1dbf8d710eb7d10e0e2db1e6d3370c4b048efe47c7a85c4b537b60b5c11832ef25f026915b803e928c1d93f01b853131e412c6308c4c6141d1504c78823c8@127.0.0.1:30310' --networkid 1515 --gasprice '0' -unlock '0x8671e5e08d74f338ee1c462340842346d797afd3' --password node2/password.txt --mine
```
attaching:
```
$ geth attach node2/geth.ipc
```
***Note:*** that IPC has been used to attach to the nodes, this allows the clique module to be used.

### Testing Contracts
After launching the network, from the root of the repository:
```
$ npm install
$ truffle test
```

### Ion Command Line Interface
The Ion CLI is a tool for handling the passing of block headers between to blockchains written in Golang to leverage the extensive ethereum libraries. It is not a critical part of the Ion infrastructure rather is just an open utility that people can extend and use as they see fit.

In its current form the Ion CLI allows the user to connect to two separate blockchains, via RPC, and submit block headers to a validation contract.

#### Running the CLI
As mentioned in the project description this simple implementation of the validation contract is active only on a single blockchain, however the CLI is simulating the passing of the headers to and from as if it were between separate chains.

Having followed the instructions on how to setup a Clique blockchain, which is hosted on `127.0.0.1:8502`, and running another chain on `127.0.0.1:8501` we can attach the CLI.

```
$ cd /path/to/validation/src
$ make build
```
Assuming a successful build the tool can be run,
```
$ ./ion-cli [/path/to/setup.json]
===============================================================
Ion Command Line Interface

RPC Client [to]:
Listening on: 127.0.0.1:8501
User Account: 0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d
Ion Contract: 0xb9fd43a71c076f02d1dbbf473c389f0eacec559f

RPC Client [from]: 
Listening on: 127.0.0.1:8502
User Account: 0x8671e5e08d74f338ee1c462340842346d797afd3
===============================================================
>>>
```


