# Ion Stage 2: Phase 1 Roadmap

Ion Stage 2 separates the cross chain payment use case from the interoperability solution. Focus has moved away from pure atomic exchange of value across chains, towards a mechanism to prove state across chains.

Given two blockchains A and B the state of B should be verifiable on A and vice-versa. To do this a smart-contract should be developed which effectively becomes a light client of the corresponding blockchain. This facilitates interoperability on a much larger scale than simply value transactions.

However the cross chain payment still serves as an illustrative example of how the solution would work in a specific scenario, but it is not part of the core solution to be developed. The Ion Stage 2 relies on a smart contract that allows storing and verification of the state of another chain. The verification of blocks from a foreign blockchain on a native chain should leverage the underlying consensus of the chain that is to be passed.

In Phase 1 we intend to tackle two different aspects:
* Storing and proving of state from a foreign chain
* Validation of passed foreign state

We assume the context of the validator set is an existing known set of nodes that also engage in the consensus protocol of a foreign chain so as to ensure the validity of signed state.

## Nomenclature
We define are a common set of terminology which will be used herein that have a specific meaning or context within the project.

* Native Chain: Refers to the chain where active interaction is taking place. This is where state from the foreign chain will be sent to for persistence and be verified against.
* Foreign Chain: Refers to the chain whose state is being persisted. Signed data about blocks from this chain will be passed to the native chain and stored.

The above naming scheme applies in the context of describing interaction flow in a one-way setting where state is passed from one chain to another. In more complex flow where both chains are actively interacted with, this naming convention may be omitted to reduce confusion.

* Proof: Refers to merkle proof-like mechanisms to assert the existence of an item or event in a block or root hash
* Validation: Refers to the signing and verifying of signatures of a piece of data, usually the block hash/header of a foreign chain
* State: Refers to the data captured and transferred between chains that facilitates the ability to prove a state transition of another chain. This will consist of another chain's block header and block hash.

## Targets
State-proving is our fundamental goal but we frame our use-case around performing a cross-chain PvP payment. Any programmable contract can interface with the Relay contract and we outline an initial example through our use-case of how this would be achieved.

In order to perform a cross-chain PvP payment we must:
  * Prove state of blockchain A on B
  * Verify the signatures of block signatories from a different blockchain
  * Settle a cross-chain transaction between two counterparties via the above mechanisms
  * Provide well-documented interfaces to allow users to easily interact with the project

### Assumptions
Listed here are the assumptions with which the project is being developed:
  * Ethereum-based blockchain
  * PBFT Consensus Protocol or other immediate-finality algorithms
  * Permissioned network
  * Validator set is known and assumed as correct

## Project Planning
Ion stage 2 will be developed using agile methodologies, with weekly sprints. Note that the sprint objective will remain dynamic and should change as the project continues.

### Sprint 1 - PoC Final Proposal Definition.
Date: 04.06.2018 - 08.06.2018

Description:
We aim to describe fully how the Phase 1 PoC would work, detailing in entirety the functionality of all smart-contracts to be developed.

Output:
  * Project specification.

### Sprint 2 - Skeleton Implementation
Date: 11.06.2018 - 15.06.2018

Description:
It should be shown that it is indeed possible to prove the state of a foreign on a native chain and make assertions of that state. Separately it should be shown that the validators from foreign chain can be added to the native chain. Blocks submitted and validated on the foreign chain validated on the native chain using the signature of the foreign validator set.

Output:
  * Smart contract for state proof verification
  * Tests for state proofs
  * Smart contract for block validation
  * Tests for block validation

### Sprint 3 - Validation of Passed State
Date: 18.06.2018 - 22.06.2018

Description:
The two separate problems of validation and proofs should be integrated and a minimum smart-contract that allows the immediate validation of a submitted block be developed.

Output:
  * Single contract which allows state proof and block validation to be performed simultaneously

### Sprint 4 - User Flow Development
Date: 25.06.2018 - 29.06.2018

Description:
Given the original user stories the smart contract should now contain minimum functions necessary to interact with the project. This should naturally be an extension of the previous week to smooth out the integration and interaction flows of the stack.

Output:
  * Smart contract should now have protection for edge-cases
  * Addition of user permissioniong
  * Automation of block generation

### Sprint 5 - Tooling and Documentation
Date: 02.07.2018 - 06.07.2018

Description:
Develop the tooling and documentation necessary for users to clone the repository and run the base functions immediately. We should write enough API documentation to allow developers to immediately be able to interface their own interoperability contracts to the Relay contract.

Output:
  * CLI
  * Tutorial
  * Development Documentation

### Sprint 6 - Testing and QA
Date: 09.07.2018 - 13.07.2018

Description:
Enhance testing to show attack resilience and any known vulnerabilities.

Output:
  * Complete test coverage
  * Code Review
  * Documentation of all possible vulnerabilities

