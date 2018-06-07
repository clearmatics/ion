# Ion Interoperability Protocol

The Ion Interoperability Protocol provides mechanisms to perform atomic swaps and currency transfers
across multiple turing-complete blockchains.

Ion consists of 3 core smart contracts:
* IonLock: Escrow contract where funds are deposited to and withdrawn from
* IonLink: Maintains state of counter-blockchain and verifies withdrawals with merkle proofs
* ERC223 Token: A placeholder ERC223 Token to perform exchanges with.

A tool called Lithium is an event relay used to facilitate to communication between the chains. Lithium forwards `IonLock` deposit events to the opposite chain's `IonLink` as a state update to inform of a party's escrowing of funds.

Check out the [Wiki](https://github.com/clearmatics/ion/wiki) for more detailed explanations.

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

Install all the dependencies which need Node v8+, NPM, and Python 2.7. Furthermore it is recommended to use a isolated Python environment with a tool such as `virtualenv`, Vagrant or Docker.

```
$ make
```

### Testing

Prior to running contract tests please launch an Ethereum client. A simple way to do this is through the `ganache-cli` or alternatively use the `npm run testrpca`.

Test:
```
$ make test
```

This will run both the Javascript tests for the smart contracts and the Python tests for the Lithium RPC relay.

Additionally contributors to this project should use linting tools when making commits, for both solidity and python code.

```
$ make python-lint
$ make solidity-lint
```

## Setup

To perform cross-chain payments, the contracts must be deployed on each chain.

Deploy two testrpc networks, if necessary, in separate terminals:
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

## Tutorial

The following tutorial describes how to perform token transfer between two accounts on separate blockchains.

There is an example script that runs through a basic flow from start to finish that is designed to run on the testrpcs.
```
$ ./example.sh
```

This tutorial leverages Ganache and Truffle but could easily be performed on other test networks.

To perform cross-chain payments, the contracts must be deployed on each chain, which for the sake of simplicity the account and contract addresses are assumed to be the same on both chains.

ALICE=0x22d491bde2303f2f43325b2108d26f1eaba1e32b

BOB=0xffcf8fdee72ac11b5c542428b35eef5769c409f0

TOKEN=0x254dffcd3277c0b1660f6d42efbb754edababc2b

IONLOCK=0xc89ce4735882c9f0f0fe26686c53074e09b0d550

IONLINK=0xcfeb869f69431e42cdb54a4f4f105c19c080a601

IP=127.0.0.1

PORT_A=8545

PORT_B=8546

API_PORT_A=8555

API_PORT_B=8556

It is recommended to use an isolated python environment.

### Step 1: Deploy event listeners

Launch `lithium` listener A:
```
$ python -mion lithium --rpc-from $IP:$PORT_A --rpc-to $IP:$PORT_B --from-account $ALICE --to-account $BOB --lock $IONLOCK --link $IONLINK --api-port $API_PORT_A
```
Launch `lithium` listener B:
```
$ python -mion lithium --rpc-from $IP:$PORT_B --rpc-to $IP:$PORT_A --from-account $BOB --to-account $ALICE --lock $IONLOCK --link $IONLINK --api-port $API_PORT_B
```

### Step 2: Mint token on each chain

Mint for Alice on chain A:
```
$ python -mion ion mint --rpc $IP_A:$PORT_A --account $ACC_A --tkn $TOKEN_ADDR --value 5000
$ Token minted.
$ New balance = 5000
```

Mint for Bob on chain B:
```
$ python -mion ion mint --rpc $IP_B:$PORT_B --account $ACC_B --tkn $TOKEN_ADDR --value 5000
$ Token minted.
$ New balance = 5000
```

### Step 3: Escrow funds and get proof on chain A

Alice deposits to IonLock on chain A:
```
$ python -mion ion deposit --rpc $IP_A:$PORT_A --account $ACC_A --lock $LOCK_ADDR --tkn $TOKEN_ADDR --value 5000 --ref stuff
$ Token transferred.
$ New balance = 0
```

Alice finds her merkle proof on chain A:
```
$ python -mion ion proof --lithium-port $((PORT_A + 10)) --account $ACC_A --lock $LOCK_ADDR --tkn $TOKEN_ADDR --value 5000 --ref stuff
$ Received proof:
$ Path  0  :  96504857948636356700030147503635580074187355628971816059136194586624797022097
$ Path  1  :  94482339386605321136956967184442353585778610538212146199456190006347461027622
$ Path  2  :  4063950032426277920165979059513600522532612014504803720221874727295772434160
$ Latest IonLink block 72772631658565070356215801224320765885121569368220205553212543964032472153198
```

### Step 4: Escrow funds and get proof on chain B

Bob deposits to IonLock on chain B:
```
$ python -mion ion deposit --rpc $IP_B:$PORT_B --account $ACC_B --lock $LOCK_ADDR --tkn $TOKEN_ADDR --value 5000 --ref stuff
$ Token transferred.
$ New balance = 0
```

Bob finds his merkle proof on chain B:
```
$ python -mion ion proof --lithium-port $((PORT_B + 10)) --account $ACC_B --lock $LOCK_ADDR --tkn $TOKEN_ADDR --value 5000 --ref stuff
$ Received proof:
$ Path  0  :  59798365828871537698849691593400364996559135249658580970523805101316187754033
$ Path  1  :  91398783457376278236011129913922372139721274533348447063742181262540672449047
$ Path  2  :  23390520989103446330618879673836571332049395218389607622791772153046182206533
$ Latest IonLink block 20043025639256222802481390718671994518152666652712633686609639039181086747014
```

### Step 5: Withdraw funds using proof on chain B

Alice withdraws giving her proof and reference:
```
$ python -mion ion withdraw --lithium-port $((PORT_B + 10)) --rpc $IP_A:$PORT_A --account $ACC_B --lock $LOCK_ADDR --tkn $TOKEN_ADDR --value 5000 --ref stuff
$ New balance = 5000
```

### Step 6: Withdraw funds using proof on chain A

Bob withdraws giving his proof and reference:
```
$ python -mion ion withdraw --lithium-port $((PORT_A + 10)) --rpc $IP_B:$PORT_B --account $ACC_A --lock $LOCK_ADDR --tkn $TOKEN_ADDR --value 5000 --ref stuff
$ New balance = 5000
```
