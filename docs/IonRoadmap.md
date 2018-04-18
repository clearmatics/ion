# Ion Roadmap

## Interoperability Background Reading

* [Hashed Timelock Contracts](https://en.bitcoin.it/wiki/Hashed_Timelock_Contracts)
* [State Channels](http://www.jeffcoleman.ca/state-channels/)
* [Plasma](https://plasma.io/)
* [Plasma Cash](https://ethresear.ch/t/plasma-cash-plasma-with-much-less-per-user-data-checking/1298)
* [Polkadot](https://polkadot.network/)
* [Cosmos](https://cosmos.network/)

## Assumptions

* Permissioned Blockchain network
* Cryptographic off chain tools are available - such as hashing algorithms or signatures
* Gas usage is not limited
* No minimums on throughput

## Design of the Solution

* There are 2 distinct approaches that allow value/state to be synchronously transferred between two separate blockchains
  - 2nd layer applications: developing and deploying smart contracts
  - Protocol level: potentially requiring modification to consensus mechanism and the DCN

## Implementation

### 2nd layer application

Description of the Ion Stage 1 can be found [here](https://github.com/clearmatics/ion/wiki/Ion-Stage-1).

Some of its characteristics are:
  * Funds deposited to an escrow contract
  * Merkle tree is created from details of deposit, root of which is given to the opposite chain
  * Withdraw dependent on merkle proof using the data of the deposit
  * A service is responsible for updating the root used to verify the proof for withdrawing funds form the escrow
  * Mechanism design and use of signatures can improve the robustness of the main mechanics for syncing state

### Protocol level development

* Sketch of this approach was described on the Phase II Technical Architecture
* Some of the main points on this solution would be:
  - Consensus is able to read and update state through EVM - meaning consensus algorithm can interact with smart contracts
  - An underlying chain, V-Chain, is used to sync state between multiple blockchains
  - Nodes that belong to V-Chain and to counterparty chains are described as validators (since they have the ability to sync the state locally and publish a block with the current state)

## Roadmap Plan

Although these indicate a certain rigid structure of the development cycle, the main goal is to highlight the focus areas during each phase.

The optimal situation is to have continuous deployment, that can be reviewed and improved by the stakeholders, following Agile methodologies.

### First release (April - May 2018)

#### Description
The main goal of this release (named Ion Stage 1) is to have some practical experience on approaching the interoperability problem.

This release is intended to be a simple version of a 2nd layer interoperability solution.

#### Goals
  - Develop first simple prototype to demonstrate minimum viable cross chain payments
  - Write and publish report with the specific mechanics of the solution and recommendations for future work
  - Engage in educational discussion of a specific problem within Clearmatics
  - Open source release of Ion stage 1

#### Tasks
The release will include:
  - [ERC223](https://github.com/ethereum/EIPs/issues/223) Token: similar to the [asset token](https://github.com/clearmatics/asset-token)
  - IonLock: contract responsible for locking funds until an appropriate root as been submitted to IonLink
  - IonLink: contract that manages roots used to prove valid withdraws from IonLock
  - Lithium: off chain service that continuously listens to the chains waiting for events from IonLock deposits
  - Beryllium: off chain service responsible for maintaining multiple TxO objects and storing the roots created
  - UI for demonstration of the concept

* Open Source and Community
  - Release of simple implementation as work in progress to github
  - Thread on [ethresear.ch](http://ethresear.ch), or any suitable alternate forums, regard the Ion Stage 1

#### Resources
* 3 Developers / Technical researchers


### Second Release (May - July 2018)

#### Description
The second release should try to solve some of the key issues with Ion Stage 1. The scope of the research should be broadened and solutions to the key issues implemented.

#### Goals
  - Reconcile the business requirements with the technical solutions
  - Develop in the open, i.e. push changes to open github repo as WIP
  - Focus the discussion within the community towards our business requirements
  - Produce technical report on findings, with recommendations for continued work

#### Tasks
  * Research of other approaches
  * Business requirements specification
  * Improvement of Ion Stage 1 (into Stage 2)
    - Addition of signatures and other cryptographic proofs to improve robustness
    - Adapt the solution to business requirements where it is needed
  * Definition of goals for next approach to be implemented through a PoC

#### Resources
  - 3 Developers / Technical researchers
  - 1 Business analyst  

### Third checkpoint (July - September 2018)

#### Description
The maturity of the product should be at a stage where we understand what a viable generic solution could be.

#### Goals
  - First approach with a known solution to clients
  - Show how can interoperability solve a specific business problem
  - Marketing material release such as videos and blog posts about presentations to clients

#### Tasks
 * Validation of approach for interoperability with clients
  - New business requirements
  - Knowledge exchange between clients and Clearmatics on the subject
 * Implementation of PoC agreed previously
 * Development of UI demo for clients
  - Validation of business requirements
  - Demo the mechanism appropriate for the business requirements
  - Hightlight the business problem that interoperability solves
 * Discuss and define what integration tests should be done on a release of this kind

#### Resources
  - 3 Developers / Technical researchers
  - 1 UI (demo builder)
  - 1 Business analyst


### Fourth checkpoint (September - November 2018)
#### Description
Following the previous increment the generic viable solution should have been modified to a specific use case with potential clients.

#### Goals
  - Start thinking about production (what new problems it brings for our specific solution)
  - Make a decision to follow (and funnel into) in the next months

#### Tasks
  * Build basic QA framework (or test suite) against which development will have to happen in the future
  * Define what is the Clearmatics solution on the interoperability problem
    - This would be the result of the research and business requirements being nailed down
  * Define the main bottlenecks with the solution chosen and how those can be approached
    - Bottlenecks might be created by technology or business requirements
  * Define architecture of the implementation

  * Open Source and comunity
    - Query the comunity for opinions on preferences according to their business requirements
    - Closed open WIP and focus on chosen solutions
#### Resources
  - 3 Developers / Technical researchers
  - 1 UI (demo builder)
  - 1 QA Tester
  - 1 Business analyst (full time, requirements need to be determined before next phase)

### Fifth checkpoint (November 2018 - February 2019)
#### Description

#### Goals
  - Follow good practices during development (hopefully they are in place by this time)
  - Manage output of production (post its on walls to follow progress)

#### Tasks
  * Implementation of interoperability solution
  * Implementation of test suite

  * Open Source and community
    - Nothing specific (too much involvement will be more noise than help)

#### Resources
  - 3 Developers / Technical researchers
  - 1 UI (demo builder)
  - 1 QA Tester
  - 1 Business analyst

### Sixth checkpoint (February - April 2019)
#### Description
A

#### Goals
  - Ship to clients
  - Follow ticket system and do support

#### Tasks
  * Packaging of solution
  * Presentation of solution to friendly clients

  * Open Source and community
    - Nothing specific (too much involvement will be more noise than help)

#### Resources
  - 2 Developers / Technical researchers
  - 1 UI (demo builder)
  - 1 QA Tester
  - 1 Business analyst

### Seventh checkpoint (April - ...)

#### Description

#### Goals
  - Go back to starting point
  - Revisit old solutions
  - Research new state of the art

#### Tasks
  * Release of final solution

  * Open Source and comunity
    - Marketing material in the form of blog posts and anouncements

#### Resources
  - 1 Developers / Technical researchers
  - 1 UI (demo builder)
  - 1 Business analyst
