# Performing On-Chain Sync

This is a proof of concept of the syncing of plasma chain blocks to the onchain `IonLink` contract. This shows the practical ability to issue payments which are committed to a local plasma chain to be state-validated by onchain contracts. From this proof of concept, we can extend application to prove the ability to track state across different chains and to create a service to relay between these two chains.

As of current I have modified the codebase to utilise `IonLink` smart contract to verify the plasma state chain of transactions. This should demonstrate the capturing of transactions on the plasma chain and subsequently saving this state to the Ethereum blockchain for validation.

## Deploying the contracts

The only contract being interacted with is `IonLink`.

Initiate a `ganache` instance to test on if necessary.

```bash
$ ganache-cli gas=0xFFFFFFFFF -p <any_port_number>
```

Compile the contracts with `truffle`.
```bash
$ truffle compile
...
Compiling ./contracts/IonLink.sol...
...
```

Then deploy them to your chosen network. Here I have configured my `truffle.js` and added a network called `testnet1` that is simply another `ganache` instance on a custom port. The `reset` flag is used to re-migrate all compiled contracts to make sure it has all the updated contracts if you've modified and recompiled contracts.
```bash
$ truffle migrate --reset --network testnet1
  Using network 'testnet1'.

  Running migration: 1_initial_migration.js
    Deploying Migrations...
    ...
  Saving successful migration to network...
    Deploying IonLink...
    ... 0x6fd6fc54f4681b94a70ff36b858524a57cd8d3628eb3ad5918dd4069542ac218
  Saving artifacts...
    ... 0x2f728029538705e8c596ebb1b50480af02dfedf0acac1b98c6e289724b51c42f
    IonLink: 0x1c62c3336b1e25fd601a5261ce294c2791d5151f
  ...
```

The `IonLink address` here will be used in a later step.

The contracts are now successfully deployed to your network.


## Initiating payments

This implementation requires a clean plasma chain. If your directory has a `chaindata` subdirectory, this needs to be removed. Back-up (where necessary) and delete/move this out of the project directory.

The below command will make a payment between two randomly generated addresses with 'currencies' being denoted as the same as their source address. I.e. If Alice with address `0x1234` sends a currency, this currency will have an id of `0x1234` and acts as Alice's native currency.
All addresses supplied below as arguments should be without the starting `0x`.
```bash
$ python -mion rpc client main --ion-rpc <host>:<port> --ion-account <account_address> --ion-contract <ionlink_contract_address>
```

Each use of the command above will initiate a single payment between two addresses. If the plasma chain does not exist, this will create a new plasma chain, with a genesis block and transactions will be added to the block above this.
Each call will perform a single payment and then commit the block to the plasma chain and sync the plasma block headers to the blockchain via `IonLink`. This means that there will only be a single transaction per plasma block.

Use multiple calls of this command to add more transactions (and thereby blocks) to the plasma chain and subsequently to the `IonLink` chain. Once called, two addresses are generated, for each sender and recipient involved in this round of transaction, and the sender is populated with some native currency. Then a payment is initiated from the sender to the recipient where a payment is signed with the hash of the previous plasma block, then added to the plasma TxPool. This is then committed and transactions are pruned from the current active TxPool, sealed into a block and committed to the plasma chain. Once written to the plasma chain, an `IonLink` sync is attempted by providing the new plasma block root and verifying that the last known block hash onchain is the same as the last known plasma block hash offchain. If this assertion fails, the blockchain will not be updated. This will need to be modified to recursively update the blockchain if we are several blocks behind (i.e. we have made many valid plasma chain commitments but they have yet to be synced to `IonLink`). Otherwise the block root is handed to `IonLink`, it forms a new block hash with the hash of the previous block and new block root and set the new hash the latest block hash.

## Tracking the Plasma chain

To view the contents of the plasma blocks:
```bash
$ python -mion plasma chain -l
```
This will access the plasma chain on disk and list out all blocks up until the latest including all transactions in each block. It will detail the addresses of the participants and their balances.

Latest blocks are output first.

## Tracking the Ion Chain

To view the contents of the tracked plasma chain on `IonLink`:
```bash
$ python -mion rpc client get_tree --ion-rpc <host>:<port> --ion-account <account_address> --ion-contract <ionlink_contract_address>
```
This will query the blockchain and recursively fetch all blocks and report their hashes, roots and previous blocks. This can be used in tandem with the plasma information to verify the correct state of the chains.

Latest blocks are output first.
