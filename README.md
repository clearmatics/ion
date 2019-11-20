# Ion Interoperability Framework <img align="right" src="https://raw.githubusercontent.com/wiki/clearmatics/ion/images/ionlogo.png" height="120px" />
![Ion Version](https://img.shields.io/badge/ion-v2.2.0-brightgreen.svg)
[![Build Status](https://travis-ci.org/clearmatics/ion.svg?branch=master)](https://travis-ci.org/clearmatics/ion)
[![Solidity Version](https://img.shields.io/badge/solidity-v0.5.12-blue.svg)](https://solidity.readthedocs.io/en/v0.5.12/installing-solidity.html)
[![LGPLv3](https://img.shields.io/badge/license-LGPL%20v3-brightgreen.svg)](./LICENSE)
[![Gitter](https://img.shields.io/badge/%E2%8A%AA%20GITTER%20-JOIN%20CHAT%20%E2%86%92-orange.svg)](https://gitter.im/clearmatics/ion)

The Ion Interoperability Framework is a library that provides an interface for the development of general cross-chain smart contracts. It is part of [Clearmatics'](http://clearmatics.com) http://autonity.io project.


## Introduction

We strive towards a more interconnected fabric of systems, and to this end, methods for inter-system and cross-chain communications become paramount in facilitating this fluid ecosystem.

Ion is a system and function-agnostic framework for building cross-interacting smart contracts between blockchains and/or systems. It does not restrict itself to certain methods of interoperation and is not opinionated on what specific functions it should be built for and as such is an open protocol.
Atomic swaps and decentralised exchanges can be built on Ion and facilitate the free movement of value across different blockchains. These are just two of the possible use-cases that can be developed on top of the Ion framework.

We envision Ion to evolve to become a library of tools that developers can use on any system to build cross-chain smart contracts to interoperate with any other system.

### Contents
* [Getting Started](#getting-started)
* [Interoperate with Rinkeby!](#interoperate-with-rinkeby)
    * [Ethereum to Rinkeby](#ethereum-to-rinkeby)
    * [Hyperledger Fabric to Rinkeby](#hyperledger-fabric-to-rinkeby)
* [Develop on Ion](#develop-on-ion)
    * [Ethereum to Ethereum Interface](#ethereum-to-ethereum-interface)
    * [Hyperledger Fabric to Ethereum Interface](#hyperledger-fabric-to-ethereum-interface)
* [Ion CLI](#ion-cli)
* [Contribute!](#contribute)
* [Ionosphere](#ionosphere)

## Getting Started

Clone the repository and ensure that all the components work out of the box.

### With docker

```
docker build -t ion/dev .
docker run -ti --name ion ion/dev

# To run through the following test example you will need a separate terminal window
docker exec -ti ion /bin/bash
```

### Without docker

The following minimum versions of `node` and `golang` are required.

* [`nodejs`](https://nodejs.org/en/) v10.15.0
* [`golang`](https://golang.org/) 1.8

---

Run:

```
$ npm install
```
```
$ npm run testrpc
```
```
$ npm run test
```

to test the full stack of contracts including our example flow.

The tests should pass as below:
```
    ...

      ✓ Successful Add Block (546ms)
      ✓ Fail Add Block from unregistered chain
      ✓ Fail Add Block from non-ion (44ms)
      ✓ Fail Add Block with malformed data (70ms)
      ✓ Fail Add Same Block Twice (644ms)


  74 passing (40s)

```

With that you've just interoperated your test RPC client with the Rinkeby testnet! Our repository includes some example contracts that show you how to build smart contracts that interoperate with another chain and what mechanism that looks like.

We'll now use these example contracts to show you exactly how interoperation with Rinkeby looks like.

## Interoperate with Rinkeby!

This is a quick tutorial using our example contracts included to be able to verify a state transition in a block and call a function that depends on it. We'll demonstrate that you can use the following instructions below to interoperate from the listed systems with Rinkeby.

On the Rinkeby test network, we've already deployed a contract and executed a transaction there `Trigger.sol`. The example we will run you through will attempt to interact with that transaction by proving that it occurred on that chain and use the transaction data in a subsequent 'interactive' transaction on a local network. The transaction on Rinkeby has called the `fire()` function of the `Trigger.sol` contract, which emits an event containing the address of the caller. We will attempt to use the caller address by extracting it from the Rinkeby block using merkle proofs.

### Ethereum to Rinkeby

We've already deployed some contracts to the Rinkeby test network for you to play around with!

Ion: `0x3c70A876808ae953917ddf9d95f364614a59B941`

Clique: `0x07a435c7b9df1F331505DdC05165473BEBeAFCdB`

Ethereum Block Store: `0xe812064CCA52B42F6C1D5345Bc40fb0683eAfF15`

We will deploy our own instance of the `Function.sol` contract and pass proofs to verify a transaction that we will depend on in order to execute a function in the contract. If the proofs verify correctly then the function should emit an event to indicate that it has been executed.

Procedure:
1. We'll need the CLI here, if you are using the docker container build the CLI with `cd ion-cli/ && make build` else follow instructions to build the CLI [here](./ion-cli/).
2. `./ion-cli` Starts the CLI
3. `>>> connectToClient https://rinkeby.infura.io` Connect to the Rinkeby Testnet
4. `>>> addAccount me ./keystore/UTC--2018-11-14T13-34-31.599642840Z--b8844cf76df596e746f360957aa3af954ef51605` Add an account to be signing transactions with. We've included one that already has Rinkeby ETH for you :) Password to the keystore is `test`. If you arrived late to the party and there is no ETH left, tough luck, try creating your own account and requesting ETH from a faucet. Alternatively you can run this exact thread of commands on a `ganache-cli` instance but make sure you connect to the correct endpoint in step 2.
5. `>>> addContractInstance function /absolute/path/to/contracts/functional/Function.sol` Add your functional contract instance which compiles your contract. Must be passed an absolute path.
6. `>>> deployContract function me 1000000` Deploy your contract to Rinkeby! This will return an address that the contract is deployed at if successful. This contract has a constructor that requires two parameters to be supplied when prompted:
    * `_storeAddr`: `0xe812064CCA52B42F6C1D5345Bc40fb0683eAfF15`
    * `_verifierAddr`: `0xf973eB920fDB5897d79394F2e49430dCB9aA4ea1`
7. `>>> transactionMessage function verifyAndExecute me <deployed_address> 0 1000000` Call the function. This requires you to supply the deployed contract instance address. Here you will need to supply the following data as an input to the function when prompted:
    * `_chainId`: `0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177`
    * `_blockHash`: `0x6e13edb9c701353743106de578730b3191d344a05c2e40cfd747bedc912f12cc`
    * `_contractEmittedAddress`: `0xA2e4a61a3D2ce626Ba9B3e927cfFDB0e4E0bd06d`
    * `_proof`:
    
			0xf907430af86806843b9aca00830186a094a2e4a61a3d2ce626ba9b3e927cffdb0e4e0bd06d8084457094cc1ca082c3adf1cb22c7260686fd56ea8dc66fbb95717f2e02e4c3702e329fbd57cdada0552e7f3ef6e9f932e08acc7cc15ea6a92d7c11b8ac4a4c0765ce446b3415915ff90236f851a0fb63b3ee11a1ae3b9eb765d44ff793bea3f4dcc1b3851d3d499abeb7858675a180808080808080a07768d7f0c5cf3656a1b885a44f42aa9ac25e728a0ffd42064387f806a9d4c26b8080808080808080f9017180a01766dce2b77f929553bdb672e197ac6bd7a6e1af3bd6631a4761a0f9303de264a011691ee69c053c698656fd0d1d9109fdff8f54c555892e9a7595d9e4b13cdc81a034825307fc63fa81e00dfdb9567cbe216930cc54152008fec1c1a3c0cead50eba03b07413394305087c802b1b33ba2f4adee8ef781e9c44faa0024cfcd8bec8818a04b39fd7ee87b2ca8d0445c7c14e83e281b45df8ee0717e47ddd1586be6aa27aaa06eae234d88c3b3f1dd0696681dde9e3369bdf03b6f56d04ec57c97e65b8e4ae1a0ade94ca299f8a1effdc75566e92fd728c79e856dbaf1a1536e7cf650eb952a4da0c963e18570b2c6c64c47cb4b4ee0bc66b2d7c709585a096abdc22c82665e67b4a0b976edf4c20ccb46e24370ce40d056fa4b30b37f052095cf8e6826ea520376eca024b45b59e4246a74e0a532a172ae6eda1ab7f9f7575389ebcb49f0b3730bc8e2a07e55916b539f96fff5cb000dd364526fb29f46895ca69c6bebb62f200447cec88080808080f86d20b86af86806843b9aca00830186a094a2e4a61a3d2ce626ba9b3e927cffdb0e4e0bd06d8084457094cc1ca082c3adf1cb22c7260686fd56ea8dc66fbb95717f2e02e4c3702e329fbd57cdada0552e7f3ef6e9f932e08acc7cc15ea6a92d7c11b8ac4a4c0765ce446b3415915ff901640183122a96b9010000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020080000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000f85af85894a2e4a61a3d2ce626ba9b3e927cffdb0e4e0bd06de1a027a9902e06885f7c187501d61990eae923b37634a8d6dda55a04dc7078395340a0000000000000000000000000b8844cf76df596e746f360957aa3af954ef51605f90335f851a08e337e40227c1f29eb9be66dcce5738b0ebd7e3fb8f22f3bde4c2c562c47a3e980808080808080a02285cf183e1144e639fd557c61b8f7c9a3c407ea44d900819c24e0bdd043e8c28080808080808080f9017180a072146023fa33354e1aea072971ed85a27934706b44aa09cf1caf7f3fce3f53c4a02843057917ddf23376e770bcadde52ebbddea15eff5cdee7ed48e3b1705105c6a005d7e4bee60ad1bb393c7747e4770c0967abe35d0fe711ed44653665e1ba77d6a004eca1f6f4ca52e78956f277aba13ee78c408d065c12d15094eb9bbcc4d8fc19a09bd8a88adc645ca6cfef48ab9342e0c08e001a5b06c8a075f87795e7d4397152a01bb09aa9e8d5d2c86672a3f2ba2ebb57d82cc02a871097fa144ff6686be301b1a0fad6ed01fdfa90d8568187f3aa23cdb09a512a0438d38dd6fb665e0ffdb78f31a0fa1dd46146f1740adf146d8560ad67d5d99b881d459c6e194abdbc333ef43393a0768dae614a1f5dc5d2d5a5618065c001e808d28bea10b1ae7e8ab1f3aab8bf4aa08bf7adbbd243a63deffa64afd7e90aab7b58311ef7b4c4a824e9fb611786444ca0a93de51ab6adb387212faa7880788c217d2d56cf5baf60c72ca9ec11947afdbf8080808080f9016b20b90167f901640183122a96b9010000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020080000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000f85af85894a2e4a61a3d2ce626ba9b3e927cffdb0e4e0bd06de1a027a9902e06885f7c187501d61990eae923b37634a8d6dda55a04dc7078395340a0000000000000000000000000b8844cf76df596e746f360957aa3af954ef51605
    * `_expectedAddress`: `0xb8844cf76df596e746f360957aa3af954ef51605`
8. Success! The transaction should have been created and sent. Check on [Etherscan](https://rinkeby.etherscan.io/address/0xb8844cf76df596e746f360957aa3af954ef51605) (or any other method) whether the transaction was successful or not. It should have succeeded!

What occurred above is a transaction `0xb5850f2aa95f504f77e93b5ea09c279f94120079980866f554d55ff1451cdd26` was mined to a block which emitted a `Triggered` event. The block containing this transaction was then submitted to Ion deployed to Rinkeby (the original transaction itself was also executed on Rinkeby) and validated against the PoA Clique of Rinkeby. These blocks belong to a uniquely identifable chain with ID of `0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177`. Deployment of your own functional smart contract references a specific verifier that is designed to make proofs for the `Triggered` event. Calling the `verifyAndExecute` function with a proof of the transaction in the specific block then makes checks against the submitted block data for the transaction and the event. Upon successful discovery, your function executes.

#### Try out your own functions!

Take a look at `Function.sol` and you'll find a very simple `execute()` function. Try adding your own logic or event emissions here and follow the same procedure above and you'll be able to execute your own arbitrary code with a dependence on a particular state transition.

You can now also attempt to write your own functional smart contracts using the similar skeleton to the `Function.sol` contract.

Note that all the data submitted as a proof to the function call is generated merkle patricia proofs for a particular transaction at `0xcd68852f99928ab11adbc72ec473ec6526dac3b1b976c852745c47900f6b8e30` that was also executed on the Rinkeby Testnet. This transaction emits a `Triggered` event which contains the address of the caller when emitted. The `Function.sol` contract simply consumes events of this type, verifies that the transaction occurred in a block, that the event parameters are as expected and then executes a function.

### Hyperledger Fabric to Rinkeby

Once again we've already deployed some contracts to the Rinkeby test network for you to play around with!

Ion: `0x3c70A876808ae953917ddf9d95f364614a59B941`

Base Validation: `0x07d18C468C63Fca68198776D799316612840b1A0`

Fabric Block Store: `0x7cc9155EB4a12783bE5aBa9dcaA698d695D19A7D`

We will deploy our own instance of the `FabricFunction.sol` contract and retrieve data from a submitted fabric block to use in an Ethereum contract function call. The Fabric block submitted contains two key-value pairs currently, `A: 0` and `B: 3`. We'll show that we can retrieve the value of `B` and emit this in an event thus demonstrating usage of Fabric block state in an Ethereum transaction.

Procedure:
1. We'll need the Ion CLI here. Build the CLI by following instructions found on the [repository](https://github.com/clearmatics/ion-cli)
2. `./ion-cli` Starts the CLI
3. `>>> connectToClient https://rinkeby.infura.io` Connect to the Rinkeby Testnet
4. `>>> addAccount me ./keystore/UTC--2018-11-14T13-34-31.599642840Z--b8844cf76df596e746f360957aa3af954ef51605` Add an account to be signing transactions with. We've included one that already has Rinkeby ETH for you :) Password to the keystore is `test`. If you arrived late to the party and there is no ETH left, tough luck, try creating your own account and requesting ETH from a faucet. Alternatively you can run this exact thread of commands on a `ganache-cli` instance but make sure you connect to the correct endpoint in step 2.
5. `>>> addContractInstance fabfunc /absolute/path/to/contracts/functional/FabricFunction.sol` Add your functional contract instance which compiles your contract. Must be passed an absolute path.
6. `>>> deployContract fabfunc me 4000000` Deploy your contract to Rinkeby! This will return an address that the contract is deployed at if successful. This contract has a constructor that requires a parameter to be supplied when prompted:
    * `_storeAddr`: `0x7cc9155EB4a12783bE5aBa9dcaA698d695D19A7D`
7. `>>> transactionMessage fabfunc retrieveAndExecute me <deployed_address> 0 1000000` Call the function. This requires you to supply the deployed contract instance address. Here you will need to supply the following data as an input to the function when prompted:
    * `_chainId`: `0x05cd2d1264200118dd4878c3de1050d49a1c47fa67fecf42038fc0728e38cc7b`
    * `_channelId`: `orgchannel`
    * `_key`: `B`
8. Success! The transaction should have been created and sent. Check on [Etherscan](https://rinkeby.etherscan.io/address/0xb8844cf76df596e746f360957aa3af954ef51605) (or any other method) whether the transaction was successful or not. It should have succeeded!

The above leverages Fabric block state that has been submitted to the Rinkeby chain from an external Fabric network with unique chain ID of `0x05cd2d1264200118dd4878c3de1050d49a1c47fa67fecf42038fc0728e38cc7b`. In this chain there exists a channel `orgchannel` with a block containing transactions that have mutated the ledger. The function executed queries the `FabricStore.sol` contract deployed to Rinkeby at `0x7cc9155EB4a12783bE5aBa9dcaA698d695D19A7D` and retrieves the value at key `B` in channel `orgchannel`. It then emits an event including the block number, transaction number and the value as parameters.

#### Try out your own functions!

Take a look at `FabricFunction.sol` and you'll find a very simple `execute()` function. This function currently simply takes data and emits it. However you can use this state in however you choose. Try replacing the `execute` function body with your own logic to perform a transaction using Fabric block state. It's important to note that chaincode, the smart contract equivalent in Fabric, writes values as raw strings which means the stored data is arbitrary. This means that your functional contracts should be knowledgeable of these data formats to be able to use them effectively in your smart contracts on Ethereum.

## Develop on Ion

Develop your own cross-chain smart contracts using the Ion interface!

### Core Concept

The core of Ion revolves around the concept of dependence on state. As such, to create functional cross-chain smart contracts, the idea is to be able to make the execution of a smart contract function be dependent on some state transition to have occurred. Thus we create smart contracts that perform verification of state transitions before allowing the execution of code and this can be done many ways depending on the systems that intend to interoperate.

This results in a framework that should provide a simple interface to make proofs of state and as such is comprised of two core components:
* State Validation
* State Transition

The two core layers of the Ion framework will need different implementations for each method/mechanism by which validation and verification can be achieved across two interoperating systems. The two layers, validation and state verification, are analogous to the functions of current systems, consensus and state transition respectively. Thus the Ion framework aims to provide interfaces to make interoperation between any ledger governed by any consensus mechanism with another through the development of such interfaces.

#### State Validation

In order for two systems or chains to interoperate there must be some notion of the passing of state between the systems. The passing of this state ideally must be trustless and thus the State Validation layer handles this. It's purpose is to provide a mechanism by which passed state is checked for validity and correctness. Since we draw dependence on the state of another system to trigger arbitrary code execution we must ensure that any state that is passed is indeed correct.

#### State Transition

Once state has been passed from one system to another it can be used as a dependency for code execution. This code execution should have conditions to be contingent on a certain piece of data from another system/chain. In practice this could be checking the balance of an account on another chain or asserting that some transaction has been fulfilled.

The State Transition layer should provide a mechanism to allow checks to be made against the stored state from another system. These checks should involve discerning and/or confirming a certain piece of data or event in the state of another chain and using the successful verification to trigger the execution of code.


### Ethereum to Ethereum Interface

With Ethereum, interoperation between chains is mainly a question of validation as they share the EVM. We currently have made an implementation of the Clique proof-of-authority consensus mechanism used by the Rinkeby Testnet for validation. We achieve state verifications via event consumption. Using the presence of an event in a transaction, we can verify if the expected computation was done and to only do something if the verification succeeds.

To write a smart contract that depends on particular state transitions there are pre-requisites:
* Event Verifier contracts

For any event signature, a corresponding Event Verifier contract must be written which encodes the mechanism that extracts the relevant event to check expected parameters as part of the verification.

[`TriggerEventVerifier.sol`](./contracts/functional/TriggerEventVerifier.sol) is a very simple event verifier:
* Holds the event signature it decodes
* Encodes a verification function that takes expected fields as input to check against that included in the event

All Event Verifier contracts should perform the same way. The differences will simply be in the event signature and parameters checks of the event.

[`Function.sol`](./contracts/functional/Function.sol) provides a very simple example of how an event-consuming contract is written. Changing the event to be consumed by referencing a different Event Verifier allows you to draw dependence on a different state transition.

### Hyperledger Fabric to Ethereum Interface

Fabric state is written as key-value pairs. When Fabric blocks are submitted to an Ethereum chain we store a lot of data about the block. There are two types of information we store that can be used in Etheruem transactions:

* Ledger State
* Past transitions

Since the only data stored in a Fabric block is a state transition from past value to new value for a given key from a specific chaincode, it is hard to access state directly and only storing this simply allows us to check the presence of an expected transaction in the past.

As such we keep both a copy of the ledger and the entire list of state transitions. This allows us two different accessors of Fabric state:

* Query the current value of a key
* Verify that a transaction has happened at a given time/block

With this we can build use-case smart contracts in Ethereum to be able to query state from a Fabric block for usage, or when passed with expected data, verify if a transaction in a specified block mutated the ledger in an expected way. The example [`FabricFunction.sol`](./contracts/functional/FabricFunction.sol) shows how you would query the storage contract to make use of current Fabric ledger state.


### Testing

Test-driven development and unit-testing all individual components of your smart contracts are extremely important in developing cross-chain contracts. There are two main steps to testing:

* Core functionality of smart contracts
* Integration with Ion Interface

Traditional tests that ensure that your smart contract is operating in the way that you intended is always required. However with the added use of the Ion interface, you'll need to write tests that make sure they both integrate well with the verification mechanisms and still behave in the expected way.

Study the tests in the repository to discover how we've unit-tested the entire integrated stack of smart contracts.

## Ion CLI

The Ion Command-Line Interface is a tool built in Golang to provide functions to facilitate the use of the Ion framework.

We have a current work-in-progress for a [CLI for Hyperledger Fabric](https://github.com/Shirikatsu/fabric-examples/tree/format-block/fabric-cli) forked from another project.

The Command-Line Interface reference can be found [here](./ion-cli)

## How Ion works

Please see our [wiki](https://github.com/clearmatics/ion/wiki) for more detailed information about the design philosophy of Ion.

### Setting up an environment

To start developing on Ion, you'll need access to at least two different networks.

* [Set up a Hyperledger Fabric network](https://github.com/clearmatics/simpleshares)

### Deploy the Ion stack

Use the [CLI](https://github.com/clearmatics/ion-cli) to deploy the Ion stack to your own chain.

1. Deploy `Ion.sol`.
2. Deploy storage and validation contracts.

First use the CLI to connect to a client and connect an account to submit transaction from:

```bash
$ ./ion_cli
===============================================================
Ion Command Line Interface

Use 'help' to list commands
===============================================================
>>> connectToClient <your_rpc_endpoint>
Connecting to client...

Connected!
===============================================================
>>> addAccount <given_account_name> <path/to/keystore>
Please provide your key decryption password.

Account added succesfully.
```

Deploy `Ion.sol`:

```bash
>>> addContractInstance ion </absolute/path/to/Ion.sol>
Compiling contract...
Creating contract instance...
Added!
===============================================================

>>> deployContract ion <account_name> <gas_limit>
Enter input data for parameter _id:
<unique_id_for_chain_being_deployed_to>

Waiting for contract to be deployed
Deployed contract at: 0x...
===============================================================
```

Deploy storage contract, in this example we deploy the Ethereum block store contract:

```bash
>>> addContractInstance ethstore </absolute/path/to/EthereumStore.sol>
Compiling contract...
Creating contract instance...
Added!

===============================================================
>>> deployContract ethstore <account_name> <gas_limit>
Enter input data for parameter _ionAddr:
<deployed_ion_address>

Waiting for contract to be deployed
Deployed contract at: 0x...
===============================================================
```

Deploy validation and register contract, in this example we deploy the Ethereum Clique validation contract:

```bash
>>> addContractInstance clique </absolute/path/to/Clique.sol>
Compiling contract...
Creating contract instance...
Added!
===============================================================

>>> deployContract clique <account_name> <gas_limit>
Enter input data for parameter _ionAddr:
<deployed_ion_address>

Waiting for contract to be deployed
Deployed contract at: 0x...
===============================================================

>>> transactionMessage clique register me <deployed_clique_address> 0 4000000
Marshalling ABI
JSONify ABI
Packing Args to ABI
Retrieving public key
Creating transaction
Signing transaction
SENDING TRANSACTION
Waiting for transaction to be mined...
Transaction hash: 0x...
===============================================================
>>> transactionMessage clique RegisterChain <account_name> <deployed_clique_address> 0 <gas_limit>
Enter input data for parameter _chainId:
<unique_id_of_interoperating_chain>

Enter input data for parameter _validators:
<comma_separated_list_of_validator_addresses>

Enter input data for parameter _genesisBlockHash:
<genesis_block_hash>

Enter input data for parameter _storeAddr:
<deployed_block_store_address>

Marshalling ABI
JSONify ABI
Packing Args to ABI
Retrieving public key
Creating transaction
Signing transaction
SENDING TRANSACTION
Waiting for transaction to be mined...
Transaction hash: 0x...
===============================================================
```

Deployment of other validation/storage contracts would be performed similarly.


# Contribute!

We would love contributors to help evolve Ion into a universal framework for interoperability.

Functional use-case smart contracts should not live in this repository. Please create use-cases in your own repositories and we'll include a link to them in our Ion-based contract catalogue.

The repository is segmented into two main sections that require work:
* Validation
* Storage

#### Validation

Each system requires a mechanism to be able to prove the correctness/validity of any data it holds, and this mechanism must be encoded by a Validation contract. Thus each method by which data could be validated must have its own contract that describes it. For example, to validate blocks from a proof-of-authority chain, we must replicate the verification mechanism of that specific implementation.

Validation contracts for the consensus mechanism of an interoperating chain would be required in order to interact with it.

These should live in the `contracts/validation/` directory.

#### Storage

Each system holds its data in different formats, and subsequently proving that the data exists would be different. Thus a different storage contracts must be written that decode and store any arbitrary data formats for use on any other system. For example, proving a transaction exists in an Ethereum block is different from proving a UTXO in a Bitcoin block.

Storage contracts for the data format and state verification mechanisms of an interoperating chain would be required in order to interact with it.

These should live in the `contracts/storage/` directory.

With the above, we aim to expand our system-specific implementations for verification of both data validity and state transitions to allow the easier development of smart contracts using these interfaces.

# Ionosphere

All Ion-extensions or applications built on Ion are part of the Ionosphere! If you have built a cross-chain smart contract use-case please add a reference to your project here.

## Applications

* [Transact-Ion](https://github.com/Shirikatsu/Transact-Ion), an atomic swap contract built on Ion.
* [Simple Shares](https://github.com/clearmatics/simpleshares), a DvP model between Fabric and Ethereum using Ion.
* [web3j Ion](https://github.com/web3j/ion), a simple start application using Ion with web3j.
