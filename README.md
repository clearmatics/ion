# Ion Interoperability Framework
The Ion Interoperability Framework is a library that provides an interface for the development of general cross-chain smart contracts.

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

You will need [`nodejs`](https://nodejs.org/en/) and [`golang`](https://golang.org/) installed.

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

    Check Roots Proof
	Gas used to submit check roots proof = 124536 gas
      ✓ Successful Check Roots Proof (138ms)
      ✓ Fail Roots Proof with wrong chain id (106ms)
      ✓ Fail Roots Proof with wrong block hash (100ms)
      ✓ Fail Roots Proof with wrong tx nodes (192ms)
      ✓ Fail Roots Proof with wrong receipt nodes (132ms)


  69 passing (43s)

```

With that you've just interoperated your test RPC client with the Rinkeby testnet! Our repository includes some example contracts that show you how to build smart contracts that interoperate with another chain and what mechanism that looks like.

We'll now use these example contracts to show you exactly how interoperation with Rinkeby looks like.

## Interoperate with Rinkeby!

This is a quick tutorial using our example contracts included to be able to verify a state transition in a block and call a function that depends on it. We'll demonstrate that you can use the following instructions below to interoperate from the listed systems with Rinkeby.

### Ethereum to Rinkeby

We've already deployed some contracts to the Rinkeby test network for you to play around with!

Ion: `0x9d6E5614E2a9714e5dC66d9704338b189cc60829`

Clique: `0x4453E9980b55ccf9c0Bf7DCeA5cA95098c59B473`

Base Validation: `0x4E377dE05cA963B4b1f735Ca499C07eFED371760`

Ethereum Block Store: `0x820b05784BDFbA088831bCde20458226B2Db05DC`

Fabric Block Store: `0x1cc3c07Bbeea94d0224878AfaB53514f1C7F99C0`

We will deploy our own instance of the `Function.sol` contract and pass proofs to verify a transaction that we will depend on in order to execute a function in the contract. If the proofs verify correctly then the function should emit an event to indicate that it has been executed.

Procedure:
1. We'll need the CLI here, [build the CLI](./ion-cli/).
2. `./ion-cli` Starts the CLI
3. `>>> connectToClient https://rinkeby.infura.io` Connect to the Rinkeby Testnet
4. `>>> addAccount me ./keystore/UTC--2018-11-14T13-34-31.599642840Z--b8844cf76df596e746f360957aa3af954ef51605` Add an account to be signing transactions with. We've included one that already has Rinkeby ETH for you :) Password to the keystore is `test`. If you arrived late to the party and there is no ETH left, tough luck, try creating your own account and requesting ETH from a faucet. Alternatively you can run this exact thread of commands on a `ganache-cli` instance but make sure you connect to the correct endpoint in step 2.
5. `>>> addContractInstance function /absolute/path/to/contracts/functional/Function.sol` Add your functional contract instance which compiles your contract. Must be passed an absolute path.
6. `>>> deployContract function me 1000000` Deploy your contract to Rinkeby! This will return an address that the contract is deployed at if successful. This contract has a constructor that requires two parameters to be supplied when prompted:
    * `_storeAddr`: `0x820b05784BDFbA088831bCde20458226B2Db05DC`
    * `_verifierAddr`: `0xf973eB920fDB5897d79394F2e49430dCB9aA4ea1`
7. `>>> transactionMessage function verifyAndExecute me <deployed_address> 0 1000000` Call the function. This requires you to supply the deployed contract instance address. Here you will need to supply the following data as an input to the function when prompted:
    * `_chainId`: `0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177`
    * `_blockHash`: `0xf88ef06bc1a9c60457d8a4b65c4020dae2ef7f3287076a4d2d481a1bcb8e3148`
    * `_contractEmittedAddress`: `0x5dF43D6eaDc3EE940eCbf66a114486f3eF853da3`
    * `_path`: `0x04`
    * `_tx`: 
    
			0xf86808843b9aca00830186a0945df43d6eadc3ee940ecbf66a114486f3ef853da38084457094cc1ca08231f8f3c7c32c425d43418053eea8f3de09a64e40833329d8ca94d118498f72a00373c76a75251dc0da8f06658e3a261a31ec85759102266ee9110c022f9d45d5
    * `_txNodes`:
    
    		0xf90236f851a0c478a441c408d00ad410c89a76a635913325eb62a8650ccbc96e8998d50e36dc80808080808080a0edb4c44cbd3957a9226f30b982449d487ff09ff0e933a25f1c005d9df93289c38080808080808080f9017180a07ad13782edc2465b8d6c914d6a18368597b6c906b420d5e3a3dd5dbc408166fba0600296a0213fdce5c37d0520f8e23c32c90bad0e8d9163edb6357fe794546e09a05e5cdcbc193fe8a2966f8d5ecc6d94bd527c37b3f6ae0a3530c9f22cd8efc1b4a09fcfbbffd3fc9b6d7bfc57ec78811f425a24c3828d2ac203c06c7540dd382514a0a0b78e007d6cbcd9fd048ebc6397ca0d067397fb5d61f32fdc118377eaf1039ba032a122caf55ecfa596d4671830e058d6116902ada20591c72f43b299d6b0cb95a06c1f1b6f656336723110875d56d38470520fc82eea9939956ea7c77cd98de7e4a0bb2cffa7c59f46bfadcb40f32a8294ed45c991e6b6e858fc1a81c1ae8e546800a0916fe7ee83eed2a2864dcb5989522af36983674b9a89b8ab856f0e92236383a7a0fa44c96e10751887a97d8855c941805d5de5a7bbb31dfbcdb7100144d35ac4cca0c43a7ef5f5eb44ddd864b3fb9e7be9953e49c5ca089c034175d98b80e5f461718080808080f86d20b86af86808843b9aca00830186a0945df43d6eadc3ee940ecbf66a114486f3ef853da38084457094cc1ca08231f8f3c7c32c425d43418053eea8f3de09a64e40833329d8ca94d118498f72a00373c76a75251dc0da8f06658e3a261a31ec85759102266ee9110c022f9d45d5
    
    * `_receipt`: 
    
    		0xf90164018311aac9b9010000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000100000000000000000000800000000000f85af858945df43d6eadc3ee940ecbf66a114486f3ef853da3e1a027a9902e06885f7c187501d61990eae923b37634a8d6dda55a04dc7078395340a0000000000000000000000000b8844cf76df596e746f360957aa3af954ef51605
    * `_receiptNodes`: 
    
    		0xf90335f851a08fbb95708b2169b98ca70955d2280ea41a4490918c7097b9fbb0dc02b6c1021d80808080808080a07e520e72a52d285315cad163ae68ff5f5b9d7b2efb8bf38488428249580b8aee8080808080808080f9017180a03f63d06b509ff798f0c456e58402326006f83b1dfbddf00ee61810754489817ca05f0c69a424cccf549fa2a1e460cab220f18238fac997214049480a1c48f320eaa0faee39d4012a0db36c610ac5f41306041dee4488d47bc41242afb40582e2b8ffa0d298be00622dc72ef2d1a2df8a259bcf769d36e4502ef8b1b2dbde2d19736c3ba04f78ffe7b29e2192dee6197eeeb19f56801dfe092c3fc485f8e2ca0141dddb1ca0e1c200d90941dbf37ac5d0f37e40479de4b2fe1a0a1b398c9f64cf80d6010eb0a034d7af644f909a94dc41bfa08ae9d0ea1c0d053e0c520378e8a5f2efd6422b62a068858e77fadd9ae975d9aff2e23e40b0f50dab934654a261930a40101ce370a0a068be4c01c98beead34d05516b5236e6196d85125b89e2632f23bdf16e7b16fe4a09f0c3ad7917f3f28a66d151270d2cd4ece50fecb2bd769c38512f1e9c5f9355ea05e1856c7908a4428edb582898261e563b2356fa3d1134357b47799a3e470bcc18080808080f9016b20b90167f90164018311aac9b9010000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000100000000000000000000800000000000f85af858945df43d6eadc3ee940ecbf66a114486f3ef853da3e1a027a9902e06885f7c187501d61990eae923b37634a8d6dda55a04dc7078395340a0000000000000000000000000b8844cf76df596e746f360957aa3af954ef51605
    * `_expectedAddress`: `0xb8844cf76df596e746f360957aa3af954ef51605`
8. Success! The transaction should have been created and sent. Check on [Etherscan](https://rinkeby.etherscan.io/address/0xb8844cf76df596e746f360957aa3af954ef51605) (or any other method) whether the transaction was successful or not. It should have succeeded!

#### Try out your own functions!

Take a look at `Function.sol` and you'll find a very simple `execute()` function. Try adding your own logic or event emissions here and follow the same procedure above and you'll be able to execute your own arbitrary code with a dependence on a particular state transition.

You can now also attempt to write your own functional smart contracts using the similar skeleton to the `Function.sol` contract.

Note that all the data submitted as a proof to the function call is generated merkle patricia proofs for a particular transaction at `0xcd68852f99928ab11adbc72ec473ec6526dac3b1b976c852745c47900f6b8e30` that was also executed on the Rinkeby Testnet. This transaction emits a `Triggered` event which contains the address of the caller when emitted. The `Function.sol` contract simply consumes events of this type, verifies that the transaction occurred in a block, that the event parameters are as expected and then executes a function.

### Hyperledger Fabric to Rinkeby

Once again we've already deployed some contracts to the Rinkeby test network for you to play around with!

Ion: `0x9d6E5614E2a9714e5dC66d9704338b189cc60829`

Base Validation: `0x4E377dE05cA963B4b1f735Ca499C07eFED371760`

Fabric Block Store: `0x1cc3c07Bbeea94d0224878AfaB53514f1C7F99C0`

We will deploy our own instance of the `FabricFunction.sol` contract and retrieve data from a submitted fabric block to use in an Ethereum contract function call. The Fabric block submitted contains two key-value pairs currently, `A: 0` and `B: 3`. We'll show that we can retrieve the value of `B` and emit this in an event thus demonstrating usage of Fabric block state in an Ethereum transaction.

Procedure:
1. We'll need the CLI here, [build the CLI](./ion-cli/).
2. `./ion-cli` Starts the CLI
3. `>>> connectToClient https://rinkeby.infura.io` Connect to the Rinkeby Testnet
4. `>>> addAccount me ./keystore/UTC--2018-11-14T13-34-31.599642840Z--b8844cf76df596e746f360957aa3af954ef51605` Add an account to be signing transactions with. We've included one that already has Rinkeby ETH for you :) Password to the keystore is `test`. If you arrived late to the party and there is no ETH left, tough luck, try creating your own account and requesting ETH from a faucet. Alternatively you can run this exact thread of commands on a `ganache-cli` instance but make sure you connect to the correct endpoint in step 2.
5. `>>> addContractInstance fabfunc /absolute/path/to/contracts/functional/FabricFunction.sol` Add your functional contract instance which compiles your contract. Must be passed an absolute path.
6. `>>> deployContract fabfunc me 4000000` Deploy your contract to Rinkeby! This will return an address that the contract is deployed at if successful. This contract has a constructor that requires a parameter to be supplied when prompted:
    * `_storeAddr`: `0x1cc3c07Bbeea94d0224878AfaB53514f1C7F99C0`
7. `>>> transactionMessage fabfunc retrieveAndExecute me <deployed_address> 0 1000000` Call the function. This requires you to supply the deployed contract instance address. Here you will need to supply the following data as an input to the function when prompted:
    * `_chainId`: `0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177`
    * `_channelId`: `0xf88ef06bc1a9c60457d8a4b65c4020dae2ef7f3287076a4d2d481a1bcb8e3148`
    * `_key`: `0x5dF43D6eaDc3EE940eCbf66a114486f3eF853da3`
8. Success! The transaction should have been created and sent. Check on [Etherscan](https://rinkeby.etherscan.io/address/0xb8844cf76df596e746f360957aa3af954ef51605) (or any other method) whether the transaction was successful or not. It should have succeeded!

The function executed queries the `FabricStore.sol` contract deployed to Rinkeby at `0x1cc3c07Bbeea94d0224878AfaB53514f1C7F99C0` and retrieves the value at key `B` in channel `orgchannel`. It then emits an event including the block number, transaction number and the value as parameters.

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

* [Set up an Ethereum network]()
* [Set up a Hyperledger Fabric network]()

### Deploy the Ion stack

Use to CLI to deploy the Ion stack to your own chain.

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

The repository is segmented into three main sections that require work:
* Validation
* Storage
* CLI

#### Validation

Each system requires a mechanism to be able to prove the correctness/validity of any data it holds, and this mechanism must be encoded by a Validation contract. Thus each method by which data could be validated must have its own contract that describes it. For example, to validate blocks from a proof-of-authority chain, we must replicate the verification mechanism of that specific implementation.

Validation contracts for the consensus mechanism of an interoperating chain would be required in order to interact with it.

These should live in the `contracts/validation/` directory.

#### Storage

Each system holds its data in different formats, and subsequently proving that the data exists would be different. Thus a different storage contracts must be written that decode and store any arbitrary data formats for use on any other system. For example, proving a transaction exists in an Ethereum block is different from proving a UTXO in a Bitcoin block.

Storage contracts for the data format and state verification mechanisms of an interoperating chain would be required in order to interact with it.

These should live in the `contracts/storage/` directory.

#### CLI

As the developments of the above two layers progress, there may be requirements to extend the CLI to ease the use of those functions i.e. block retrieval methods, data formatting etc. With that we encourage that the CLI is extended in tandem with additions contributed to the validation or storage sections for processes that may be more cumbersome to perform manually.



With the above, we aim to expand our system-specific implementations for verification of both data validity and state transitions to allow the easier development of smart contracts using these interfaces.

# Ionosphere

All Ion-extensions or applications built on Ion are part of the Ionosphere! If you have built a cross-chain smart contract use-case please add a reference to your project here.

## Applications

* [Transact-Ion](https://github.com/Shirikatsu/Transact-Ion), an atomic swap contract built on Ion.