---
title: "Ion" Interoperability Protocol
date: February 2018
toc: yes
titlepage: true
---


# "Ion" Interoperability Protocol 

Ion is a collection of components, it provides ways to perform atomic cross-chain transactions and is
designed to work towards integrating the Ethereum block-chain with real-time gross settlement systems.

Ion serves as a research platform which aims to provide easy to use building blocks which can then be
combined or used alone to tackle different scenarios and more complex types of transactions. Insights
from research will be used to guide the project into a state which can fulfill a small but viable
feature-set.

Research started by identifying types of cross-chain transactions and the general challenges encountered
when settling financial transactions, and then progressed to reviewing research to identify patterns and
technologies to create prototypes and specify methods which can meet new requirements as they were
encountered. The research resulted in draft specifications for several locking patterns and implementations
of prototypes for the following components:

 * Merkle tree generation and verification
 * Multi-currency payment netting
 * Solidity / Python / Go integration
 * Minimal Plasma chain
 * Cross-chain interoperability
 * Privacy-enhancing methods (such as ring signatures)
 * Transaction relay for Ethereum events 
 * Cross-chain notifications


## Summary of Findings

After performing the research and creating several prototypes I have come to some interesting conclusions
and a better understanding of the trade-offs involved:

 * Passing arguments to and storing data on the block-chain is very expensive and should be avoided as
   much as possible.
 * 'Pre-image reveal' locks are a good option for trustless atomic swaps
 * More complex locks can be created, but introduce significant overheads
 * Trusted side-chains (like Plasma or Cosmos) can be very flexible


## Discussion Points

 * Which scenarios will Ion be used in (e.g. use-cases)
 * How does Ion fit with other components
 * Is there an advantage to doing RTGS within Ion
 * Start with a trusted third-party, or prefer strong trustlessness.

# Desirable properties:

 * There is no window where it is possible to double-spend, e.g. waiting for 3 confirmations
 * Each transaction accepted into a block is finalised, it can never be revoked.
 * It is the duty of the receiver to verify that no message is acted upon twice
 * 
 
 * Altering the content of a merkle root would allow malicous events to be triggered and funds to be withdrawn
   because the on-chain contracts and other systems rely on the merkle root being the single source of authority.

Threat model:

 * One or more witnesses are malicious
 * All witnesses are malicious


Double-spending is prevented as long as the following constraints are enforced:

 * The transaction signature is only valid for one block
 * A transaction reference may only occur once on a block
 * Only the owner of the `from` address can update a payment by providing a new signature
 * Blocks are applied sequentially
 * The verifier nodes keep a list of all current balances
 * Payments within a block never exceed their available balance



# Cross-blockchain interoperability
 
Distributed blockchains of blockchains:
 
 * https://cosmos.network
 * https://github.com/theloopkr/loopchain ( http://docs.icon.foundation/ICON-Whitepaper-EN-Draft-4.8.pdf )

## 'Centrally banked' crypto-currencies

 * https://arxiv.org/pdf/1505.06895.pdf
 * https://iamjustatad.files.wordpress.com/2016/11/rscoin_thesis.pdf
 * https://iamjustatad.wordpress.com/2016/11/17/extending-rscoin-part-2-rscoin-architecture-and-consensus-mechanism/
 * https://github.com/gdanezis/rscoin

## Plasma

> Plasma is a way to create a blockchain (as a set of smart contracts)
within a blockchain, by simply committing the root hash of the "meta"
child chain (Plasma) on the main chain (Ethereum).
>
> The child chain acts kind of like state channels on steroids, and
only the final state transition hash is committed to the main chain.
If someone commits an invalid plasma state transition hash/block,
the protocol incentivizes committing "fraud proofs", where parties
who are interacting in the child chain can show that someone behaved
maliciously, and penalize them, reversing the faulty transition.
>
>The child plasma chains can themselves also have child plasma chains,
and those chains can have child chains, etc... this "tree" of chains
can be reduced through MapReduce, such that billions of transactions/state
transitions can be trustlessly/verifiably encoded as a simple hash
committment on the root chain.
>
> -- [mattdf](https://www.reddit.com/r/ethereum/comments/6sqca5/plasma_scalable_autonomous_smart_contracts/dleskw5/)


### Implementations

 * [Plasma MVP](https://github.com/omisego/plasma-mvp)
 * [BankEx, Plasma ETH Exchange](https://github.com/BankEx/PlasmaETHexchange)
 * [btcrelay](https://github.com/ethereum/btcrelay)
 * [etgate](https://github.com/mossid/etgate)
 

### Further Resources

 * [Minimal Viable Plasma](https://ethresear.ch/t/minimal-viable-plasma/426)
 * [grid+ - efficiently bridging blockchains](https://blog.gridplus.io/efficiently-bridging-evm-blockchains-8421504e9ced)
 * [Developer Deep Dive—ETGate, Gateway between Ethereum & Tendermint](https://blog.cosmos.network/etgate-md-6a31f049a62f)


# Locking Patterns

There are degrees of trustlessness:

 * Trustless Distributed - requires no intermediary
 * Trustless Centralised - requires an untrusted intermediary
 * Trusted Centralised - requires a trusted intermediary

There are at least three types of transactions:

 * Contingency lock - A perform X if B does Y
 * Atomic Swap - A and B swap X and Y
 * Transfer - Ownership of X is transferred from A to B


Two methods of withdrawing from Plasma chain:

 - provide proof of UXTO, timeout period for proof of forgery
   - requires 1 merkle proof to update balance
   - requires constant monitoring to provide proof of forgery
   - DoS attack - spending many UXTOs / providing proof of forgery, costs both parties an equal amount of money

 - proof of liveness, actions can only be performed if it exists in the most recently published block
   - must be trustless so relayer can't insert spends or roll-back state machines


## Pre-image Reveal to Unlock  

This category of lock is used by the [Lightning Network](https://lightning.network/), it is fully decentralised, trustless and requires no third-party.

References:

 * https://en.bitcoin.it/wiki/Hashed_Timelock_Contracts
 * https://en.bitcoin.it/wiki/Atomic_cross-chain_trading
 * https://rusty.ozlabs.org/?p=462
 * https://interledger.org/rfcs/0022-hashed-timelock-agreements/#background-on-hashed-timelock-contracts-htlcs
 * https://interledger.org/rfcs/0022-hashed-timelock-agreements/
 * https://github.com/chatch/hashed-timelock-contract-ethereum


### Hash pre-image locks

![Hash-Time lock](../docs/Hashlock.pdf)

### Elliptic Curve Locks

As with the hash pre-image reveal, except the secret is the discrete logarithm of a curve point.

![Key-Time lock, using discrete logarithm proof](../docs/Keylock.pdf)


## Third-party escrow 

This lock type requires you to lock your funds in a way which a third party must verify you have completed the precondition, if you have it will provide a signed token that can be given to the smart-contract to release the funds. The smart contract can be programmed in a way which provides the third party no advantage or ability to steal funds even if it is malicious.

This is like putting the funds into a single-use multi-signature wallet, which has a timeout after which the funds will be withdrawable by the original depositor.


## Third-party merkle proof

This lock type requires both parties to use a trusted but unbiased third-party who publishes metadata about transactions and events on different chains. Smart-contracts must be be programmed to rely on this third-party when verifying if an action can be performed.

The advantage of a merkle-proof versus a third-party signature is that the third-party doesn't need to be actively involved in every transaction. The disadvantage is that every action requiring a merkle proof may be very expensive on Ethereum due to the Gas overhead of storing additional words of data (20,000 Gas per 256 bit word).

For example, smart-contract B is programmed to send its funds to an address if (on another chain) another user has sent funds to an address, research was made into the failure conditions and possible attack scenarios, which include:

 * Funds on chain A are never sent
 * Funds on chain B are never deposited
 * Either person A or B is malicious
 * Inability to submit transactions, unlucky timeouts

There are two methods which use third-party merkle proofs:

 * Optimistic deposit, with timeout
 * Lock-step with cancellation

### Optimistic method

The first method is optimistic, the seller locks their funds for trade with a specific buyer, when the buyer sends the seller a specific transaction on another chain this provides them with evidence they can provide to the lock contract to release the funds to them.

This requires three on-chain transactions and one merkle proof:

 * Person A: Deposit & lock
 * Person B: Send Funds 
 * Person B: Proof for unlock


### Lock-step method

The lock-step method accounts for all failure conditions, but the penalty is in the overhead and number of on-chain transactions required.

In the best scenario it requires six on-chain transactions and two merkle proofs; effort was made to reduce the first step to a pre-image reveal, however an edge-case was discovered which could allow one party to cheat unless strict timeouts/cooling-off periods were enforced.

![Lock-step exchange method with cancellation](../docs/FluorideExchange.png)


# Real-Time Gross Settlement

![Netting with State Channels](../docs/Netting-With-State-Channels.pdf)

![Summed Payments](docs/summed.gv.pdf)

![Settled Payments](docs/settled.gv.pdf)

![Netting, Routed Payments](docs/routed.gv.pdf)

![Netting, Settled Routed Payments](docs/routesettled.gv.pdf)
