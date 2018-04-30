# Ion Interoperability Protocol

The Ion Interoperability Protocol provides mechanisms to perform atomic swaps and currency transfers across
multiple turing-complete blockchains.

Ion consists of 3 core smart contracts:
* IonLock
* IonLink
* ERC223 Token

A tool called Lithium is used to facilitate to communication between the chains.

## Cross chain payment

Executing cross payment of two tokens

1. Alice Deposit to IonLock
2. Wait for Lithium / Event Relay to update
3. Bob Withdraw from IonLock

This process needs to be excuted on both chains, step 3. is blocked for both chains until funds are deposited into the escrow of the opposite chain.

## Install

Install all the dependencies which need Node v9.0.0, NPM, and Python 2.7.

```
npm install
pip install -r requirements.txt
```

## Setup

To perform cross-chain payments, the contracts must be deployed on each chain.

Deploy two testrpc networks (if necessary):
```bash
npm run testrpca
npm run testrpcb
```

Compile and deploy the contracts on to the relevant networks:
```bash
npm run compile
npm run deploya
npm run deployb
```

Run the event relay to transmit state between the chains:
```bash
python -mion etheventrelay --rpc-from <IP_TESTRPC_A:PORT> --rpc-to <IP_TESTRPC_B:PORT> --from-account <FROM_ACCOUNT_X> --to-account <TO_ACCOUNT_Y> --lock <IONLOCK_ADDRESS_TESTRPC_A> --link <IONLINK_ADDRESS_ADDRESS_TESTRPC_A> --batch-size <BATCH_SIZE>
python -mion etheventrelay --rpc-from <IP_TESTRPC_B:PORT> --rpc-to <IP_TESTRPC_A:PORT> --from-account <FROM_ACCOUNT_Y> --to-account <TO_ACCOUNT_X> --lock <IONLOCK_ADDRESS_TESTRPC_B> --link <IONLINK_ADDRESS_ADDRESS_TESTRPC_B> --batch-size <BATCH_SIZE>
```

## Usage

### Mint Tokens
```bash
python -mion ion mint --rpc <ip:port> --account <beneficiary_address> --tkn <token_contract_address> --value <amount_of_token>
```

### Deposit
```bash
python -mion ion deposit --rpc <ip:port> --account <beneficiary_address> --lock <ionlock_contract_address> --tkn <token_contract_address> --value <amount_of_token> --data <arbitrary_data_payment_reference>
```
