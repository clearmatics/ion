# Ion Interoperability Protocol

The Ion Interoperability Protocol provides mechanisms to perform atomic swaps and currency transfers across
multiple turing-complete blockchains.

Ion consists of 3 core smart contracts:
* IonLock: Escrow contract where funds are deposited to and withdrawn from
* IonLink: Maintains state of counter-blockchain and verifies withdrawals with merkle proofs
* ERC223 Token: A placeholder ERC223 Token to perform exchanges with.

A tool called Lithium is an event relay used to facilitate to communication between the chains. Lithium forwards `IonLock` deposit events to the opposite chain's `IonLink` as a state update to inform of a party's escrowing of funds.

## Cross-chain payment

Cross-chain payment flow of two different tokens on their respective blockchains with Alice and Bob as the parties involved is as follows:

Chain A: Alice's chain
Chain B: Bob's chain

1. Alice Deposits to chain A IonLock
2. Wait for Lithium (Event Relay) to update chain B IonLink
3. Bob Deposits to chain B IonLock
4. Wait for Lithium (Event Relay) to update chain A IonLink
5. Alice withdraws from chain B IonLock with proof of her deposit on to chain A IonLock
6. Bob withdraws from chain A IonLock with proof of his deposit on to chain B IonLock
7. Both parties successfully withdraw and atomic swap is complete.

Note that withdrawing is blocked for both chains until funds are deposited into the escrow of the opposite chain.

### Caveats

Currently notable flaws in the design:
* All funds deposited must be withdrawn at once
* The number of tokens deposited must also equal the funds attempting to be withdrawn i.e. 1:1 exchange
* Payment references are currently redundant as proofs submitted to verify a withdrawal are only used to prove that the party has deposited on the other chain and is not used to distinguish the funds to be withdrawn as noted in the first point.

## Install and Test

Install all the dependencies which need Node v9.0.0, NPM, and Python 2.7. Furthermore it is recommended to use a isolated Python environment with a tool such as `virtualenv`.

```
$ make build
```

### Testing

Prior to running contract tests please launch an Ethereum client. A simple way to do this is through the `ganache-cli` or alternatively use the `npm run testrpca`.

Test:
```
$ make test
```

This will run both the Javascript tests for the smart contracts and the Python tests for the Lithium RPC relay.

Additionally 

## Setup

To perform cross-chain payments, the contracts must be deployed on each chain.

Deploy two testrpc networks (if necessary):
```
$ npm run testrpca
$ npm run testrpcb
```

Compile and deploy the contracts on to the relevant networks:
```
$ npm run compile
$ npm run deploya
$ npm run deployb
```

Run the event relay to transmit state between the chains:
```
$ python -mion etheventrelay --rpc-from <IP_TESTRPC_B:PORT> --rpc-to <IP_TESTRPC_A:PORT> --from-account <FROM_ACCOUNT_Y> --to-account <TO_ACCOUNT_X> --lock <IONLOCK_ADDRESS_TESTRPC_B> --link <IONLINK_ADDRESS_ADDRESS_TESTRPC_B> --api-port <PORT>
$ python -mion etheventrelay --rpc-from <IP_TESTRPC_A:PORT> --rpc-to <IP_TESTRPC_B:PORT> --from-account <FROM_ACCOUNT_X> --to-account <TO_ACCOUNT_Y> --lock <IONLOCK_ADDRESS_TESTRPC_A> --link <IONLINK_ADDRESS_ADDRESS_TESTRPC_A> --api-port <PORT>
```

## Usage

### Mint Tokens
```
$ python -mion ion mint --rpc <ip:port> --account <beneficiary_address> --tkn <token_contract_address> --value <amount_of_token>
```

### Deposit
```
$ python -mion ion deposit --rpc <ip:port> --account <beneficiary_address> --lock <ionlock_contract_address> --tkn <token_contract_address> --value <amount_of_token> --data <arbitrary_data_payment_reference>
```
