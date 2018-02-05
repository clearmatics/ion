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


# High Level Overview

The reality is that as 'block-chain' gains wider adoption there will be even more commpeting standards
and a wider collection of dispirate financial systems which will need to interoperate with each other.
Emerging technologies have the potential to break the mould by tackling existing problems in more elegant
ways and introducing new paradigms for integration and cross-system compatibility.

Ion seeks to explore the following topics:

 * Trust relationships and the nature of security

 * Cross-system cooperation and interoperability

 * Underpinning requirements for finance applications

 * Overcoming the limitations of existing protocols

 * Points of integration and common use cases

Ultimately one recognition is that compatibility with the rest of the ecosystem - be it the existing one 
or future expected one - is necessary to a greater or lesser extent for the percieved benefits of block-chain 
technology to be realised.


## Summary of Findings

Research started by identifying types of cross-chain transactions and the general challenges encountered
when settling financial transactions, and then progressed to reviewing research to identify patterns and
technologies to create prototypes and specify methods which can meet new requirements as they were
encountered. The research resulted in draft specifications for several locking patterns and implementations
of prototypes for the following components:

 * Merkle tree generation and verification
 * Multi-currency payment netting
 * Solidity / Python integration
 * Minimal Plasma chain
 * Privacy-enhancing methods (such as ring signatures)
 * Ethereum event and transaction relay
 * Cross-chain notifications

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


# Inspiration and Influences

## Theoretical Preface

The core precept of the Ion prototype is that there is a source of truth and finality which must witness
all payments so that it can keep a tally of the balances and enforce the rules. The security model of all
crypto-currencies relies on a source of validity and agreement on what validity is, usually by consensus of
many participants, to prove that history hasn't been tampered with. There are various other methods to
prevent history from being tampered with without it being easily detectable such as the rolling average
difficulty in the PoW mechanism employed by BitCoin which makes any divergence from the status quo an
expensive endeavour because it must overpower the inertia introduced by thousands of ASIC miners, but at 
the same time it means that, regardless of how much energy you invest in Proof of Work, if you try to forge
history your attempts will be mostly ignored unless there is a trust relationship where the consumer of
the information doesn't perform their own verification (such as in the case of a 'light client'), or where
finality is a window where each transaction must be 'set in stone' by the mining of subsequent blocks.   

However, I find that thinking about the security model is easier when the properties are reviewed within 
the context of more general dichotomies, such as temporal, spacial and level of benevolence or participation,
to better understand the relationships between threats, and also as a way of getting a deeper intuitive
understanding of the reasons for more general or specific design choices. The temporal categories are:
Future, Immediate and History; the spacial categories are: Local and Distant; and the adversaries are Neutral,
Malicious and Chaotic. For example, a local chaotic adversary which exists in the future is hard-drive failure.


## Trust at the point of consumption versus scalability

If you look at a program as a tree of decisions then every branch it takes relies on it knowing which
questions to ask and having to blindly depend on any answers it gets. Any component which retrieves
information from multiple sources ends up with the worst collective security properties from each one
and becomes the weakest link unless it independently verifies the information it consumes, but the more
information required for verification increases the amount of data which needs to be transmitted between
every participant. This scenario forms the basis of the axiom that trusted summaries are required to reduce
data transmission.

Whenever we consume information we trust that it has been verified at the point of origin to be correct in
terms of being honestly and accurately relayed without any attempt to deliberately deceive the consumer.
Cryptographic signatures, hashes, block-chains and and consensus algorithms provide an additional layer of
assurance by allowing the collective us, as participants in the system, to verify that we all have the same
information as everybody else.

From the perspective of a software designer and progressive logician I would prefer to have a system which
allows trustworthy properties such as non-repudibility and provenance to be verifiable whenever is necessary,
but  in a way which doesn't impose a specific structure to my trust relationships and allows me to decide
whatever I think a consensus is in order to be able to hedge against specific threats that I perceive.


## RSCoin, abstracted

RSCoin describes an incompatibilist approach which relies upon a system with a single trusted authority
which ultimately signs-off on everything as being the pinnacle standard for what is correct or not, but
in a distributed system consisting of thousands of nodes it comes with the presupposition that it should
define for us what trust is rather than simply that its structure and properties can be trusted. For RSCoin
to be truly distributed and trustable each information source should provide its own mechanisms to verify
its own reliability by co-proving that others are meeting the same or higher standards in a way which
enables the human-actors who decide what trust is to make their own choices when forming a consensus.

### If everybody was the 'Central Bank', then who would you trust? That should be the real question...

But, what is the purpose of a bank? Essentially it is to value assets, keep them in deposit, and provide
a proportionate fungible quantity of it or an item of exchangeable value on-demand. If I were to get a
mortgage for the value of my house the bank would be liable to provide something of equivalent value, like
honouring a cheque or card payment up to the limit of their liabilities or my overdraft.

However, how many levels of abstraction are theDichotomyre between this conception of 'banked' and the model proposed
by the RSCoin paper which adds an unnecessary 'centralised' requirement, it seems more like a central
validator node and single point of failure without the benefits of BitCoin segwit, where the web-of-trust
is imposed for our collective benefit without it being easy to opt-out regardless of whatever the
underlying technology can achieve. 


## Non-repudibility is the enemy of privacy

And money is just a collective hallucination... but regardless of whatever it is there is a quantity
of it, either as a physical possession or a balance held with some entity which exists as a positive
or negative potential with an intuitive probability of being met. 

Non-repudibility with a limit is used to determine the creditworthiness of somebody, every application
they make and significant event is correlated to one or more sources of truth about a specific account,
the credit rating agencies are implicitly provided with this information as per your agreement with
the lender. At the same time every dollar, pound or euro they provide in return is treated as equal
to each and every other unit of that currency in a fungible sense.

The case for involuntary non-repudibility is that evidence of your participation is available to others
to verify that your claims are true regardless of whether or not you lie. Each source of information
promotes the notion that it adheres to high standards and a set of rules which give it the distinctions
of being trustworthy and reputable - this requires a degree of impartiality and accepting corrections
to maintain accurate up to date information. Does keeping information disclosure voluntary provide any
less benefit to the parties involved in lending or other activities where trust is necessary compared
to enforced disclosure of all transactions and agreements between all parties.

But why should a system enforce one or another, or add specific properties which could be optional
depending on the use case. With a plurality of different systems with various properties unfortunately
the lowest common denominator often defines the highest level of security attainable.


## How have these insight guided Ion?

Instead of a hub and spoke architecture there are many hubs and spokes, where the architecture evolves
as a chaotic ecosystem of overlapping implementations, and amongst them they must find a harmonious
way to interoperate without unnecessary or artificial boundaries. Central Banking in a cryptocurrency
ecosystem extends further than where the coin is and who owns it, when the shared measure of value is
a currency like GBP, USD or EUR it is possible for policy to influence any parties which use it to
compare the value of one asset to another; as long as the 'Central Bank' policy has influence wherever
it needs to then why does it matter what the medium of exchange is?

The model that I think should be encouraged is to focus on the mechanism of exchange, the security
properties and trust relationships etc. and let the market determine what the best mediums are at
the point of exchange, where trust is a developing intuitive gauge rather than a package which we
are provided with, possibly even with a cost associated with it?

After distilling these ideas into their underlying mechanisms I've come to realise that all that matters
is a gaurantee that a provable history won't change. For example it isn't even necessary for the whole
blockchain to be stored if you can trust that any money you receive is backed by a redeemable
asset upon proof of your interest in it - like the note in your hand, as long as at the point of exchange
the recipient can be assured of it meeting whatever their preferences are for storage of value.

A good rule of thumb would be: never depend on any single entity which could cause the entire system
to fail, but without the implications of the SkyNet or HAL scenarios etc.


# Desirable properties:

 * There is no window where it is possible to double-spend, e.g. waiting for 3 confirmations on BitCoin
 * Each transaction accepted into a block is finalised, it can never be revoked.
 * It is the duty of the reciever to verify that no message is acted upon twice
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
 * Payments within a block modify the balance accordingly



# Cross-blockchain interoperability
 
Distributed blockchains of blockchains and cross-chain operations:
 
 * https://lightning.network/
 * https://interledger.org/ ( https://interledger.org/interledger.pdf )
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
