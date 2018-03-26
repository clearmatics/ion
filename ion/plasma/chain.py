#!/usr/bin/env python
"""
# Minimal Plasma Chain

This is a minimal 'plasma chain', it's a fast payment network
which can be used to create and exchange tokens across multiple
block-chains.


## Processing Payments

Each payment contains the following information:

 from - address: source
 to - address: destination
 currency - address: currency
 value - int256: signed integer
 ref - bytes32: unique payment reference
 deps - bytes32[]: optional array of dependent payment hashes
 sig - ecdsa: recoverable signature

The signed message (and transaction ID) is the following:

    H(prev-blockhash || from || to || currency || value || ref || deps)

Because each signed payment is only valid for one-block it must
be re-signed and re-broadcast until the payment is accepted by
the network.


### Dependencies

For every payment it is possible to specify one or more dependencies, these
are hashes of other payments within the same block. If any of the dependencies
don't exist in the same block the payment will be ommitted.

The dependency hash is:

    H(to || currency || value || ref)

The from address and the other payments dependencies are excluded from the
dependency hash so that only the forward facing direction of the graph is
verified.

The forward-only nature of dependencies allows for opportunistic trades and
atomic exchanges, for example multiple counterparties can attempt to fill
an exchange order but only one will succeed. Complex graphs of trades can be
constructed using this technique.


### Minting Tokens

Anybody is allowed to create as many tokens of their own currency
as they want by sending a special 'mint' transaction from itself to
itself using its own address as the currency.


## Smart-contract Integration

If an on-chain contract uses a trusted merkle root relay it can verify if
a payment has been made to it. After sending a payment to a contract address
an on-chain transaction must activate the contract and supply it with the
merkle proof.

The merkle proof ensures:

 - The `value` of `currency` (the contracts own address) has been
   permanently removed from the senders balance.

 - There is opportunity to use a unique reference to prevent
   double-spend by the contract.

The contract only needs the `value` and `ref` parameters and the
merkle proof (block height height, path), it has access to `msg.sender`
and its own address which is everything necessary to validate a payment
exists in the merkle root of a trusted side-chain.

The blockchain is monitored for all plasma compatible events from contracts
which are then recorded just as regular payments would be.


## Block Format

Blocks contain three parameters:

    hash - bytes32 - H(id, root, prev-hash)
    id - integer - Block sequence number
    root - bytes32 - Merkle tree root


## Merkle Proofs

Each block is a merkle tree which contains three types of items:

 - payments
 - signatures
 - balances


### Record types

Payment entry format is 124 bytes, with an additional 32 bytes per dependency:

  - from - address
  - to - address
  - currency - address
  - value - uint256
  - ref - bytes32
  - dep - bytes32[]

Signature format is 65 bytes:

  - v - byte
  - r - uint256
  - s - uint256

Balance format is 72 bytes:

  - currency - address
  - account - address
  - balance - uint256


# References

 * http://cryptonite.info/files/mbc-scheme-rev3.pdf
"""


from __future__ import print_function
import os
from base64 import b32encode

import msgpack
import click

from ..args import arg_bytes32
from ..utils import u256be, require, marshal
from ..merkle import merkle_tree

from .model import Block
from .payment import payments_apply, random_payments, payments_graphviz, SignedPayment, Payment


# --------------------------------------------------------------------
# Data persistence load/save


data_load = lambda x: msgpack.unpack(open(x))
data_save = lambda d, x: msgpack.pack(d, open(x, 'wb'))

chaindata_filename = lambda n, x: "%s.%s" % (b32encode(n[:15]).strip('=') if len(n) in [20,32] else n, x,)
chaindata_path = lambda n, x: "chaindata/" + chaindata_filename(n, x)
chaindata_load = lambda n, x: data_load(chaindata_path(n, x))
chaindata_save = lambda n, x, d: data_save(d, chaindata_path(n, x))
chaindata_exists = lambda n, x: os.path.exists(chaindata_path(n, x))

block_exists = lambda n: chaindata_exists(n, 'block')
block_save = lambda n, d: chaindata_save(n, 'block', d.marshal())
block_load = lambda n: Block.unmarshal(chaindata_load(n, 'block'))
diff_save = lambda n, d: chaindata_save(n, 'diff', d)
payments_save = lambda n, d: chaindata_save(n, 'payments', map(marshal, d))
#payments_load = lambda n, d: chaindata_save(n, 'payments', map(marshal, d))
payments_load = lambda n: chaindata_load(n, 'payments')
balances_save = lambda n, d: chaindata_save(n, 'balances', d)
balances_load = lambda n: chaindata_load(n, 'balances')


def chaindata_latest_set(n):
    """Maintain symlink to latest block"""
    latest_path = chaindata_path("", 'latest')
    if os.path.exists(latest_path):
        os.unlink(latest_path)
    os.symlink(chaindata_filename(n, 'block'), latest_path)


def chaindata_latest_get():
    """Which is the most recent block"""
    latest_path = chaindata_path("", 'latest')
    if os.path.exists(latest_path):
        return Block.unmarshal(data_load(latest_path)).hash
    return None


# --------------------------------------------------------------------
# Balance difference and aggregate balance functions


def diff_apply(diff, balances):
    for c, d in diff.items():
        for k, v in d.items():
            if c not in balances:
                balances[c] = dict()
            if k not in balances[c]:
                balances[c][k] = 0
            balances[c][k] += v
    return balances


def diff_dump(diff, balances):
    """Provides proof of balance in packed format:
         currency || address || balance"""
    return [''.join([c, k, u256be(balances.get(c, dict()).get(k, 0) + v)])
            for c, d in diff.items()
            for k, v in d.items()]


def diff_is_balanced(diff, created):
    """Check if the sum value for all payments for each currency sum to zero"""
    return all([ sum(d.values()) - created[c] >= 0
                 for c, d in diff.items() ])


def diff_is_funded(diff, balances):
    """If, after applying balance diff, are all balances above zero?"""
    for c, d in diff.items():
        for k, v in d.items():
            balance = balances.get(c, dict()).get(k, 0) + v
            if balance < 0:
                raise RuntimeError("Not funded: %r %r %r" % (
                                   c.encode('hex'), k.encode('hex'), v))
    return True


# --------------------------------------------------------------------
# Block functions


def block_genesis():
    if not os.path.exists('chaindata'):
        os.mkdir('chaindata')

    block = Block('\0' * 32, os.urandom(32))
    block_hash = block.hash

    block_save(block_hash, block)
    payments_save(block_hash , list())
    balances_save(block_hash , dict())
    chaindata_latest_set(block_hash)
    return block_hash


def block_seal(prev_hash, signed_payments, balances):
    """
    Seal block by applying signed payments to balances
    :type signed_payments: list[SignedPayment]
    """
    require( len(prev_hash) == 32 )

    # Verify that payments are valid at the current state
    payments = [sp.open(prev_hash) for sp in signed_payments]
    balances_diff, created = payments_apply(payments)

    # TODO: check for duplicate payment references...

    require( diff_is_balanced(balances_diff, created), "Misbalance" )
    require( diff_is_funded(balances_diff, balances), "Not funded" )

    # Store payments, signatures and balances in merkle tree
    dumped_sigs = sorted([sp.s.dump() for sp in signed_payments])
    dumped_payments = sorted([p.dump() for p in payments])
    dumped_balances = sorted(diff_dump(balances_diff, balances))
    tree, root = merkle_tree(dumped_sigs + dumped_payments + dumped_balances)

    return Block(prev_hash, u256be(root)), balances_diff


def blockchain_apply(prev_hash, signed_payments):
    """
    :type prev_hash: str | bytes
    :type signed_payments: list[SignedPayment])
    :returns: Block
    """
    require( len(prev_hash) == 32, "Invalid previous block hash" )
    require( len(signed_payments) > 0, "Payments required" )

    prev_balances = balances_load(prev_hash)
    block, balances_diff = block_seal(prev_hash, signed_payments, prev_balances)

    block_hash = block.hash
    if block_exists(block_hash):
        raise RuntimeError("Cannot overwrite block", block_hash)

    block_save(block_hash, block)

    #chaindata_save(block_no, 'diff', dict_dump(balances_diff))
    payments_save(block_hash, signed_payments)
    graph = payments_graphviz([sp.p for sp in signed_payments])
    graph.render( chaindata_path(block_hash, 'graphviz') )

    after_balances = diff_apply(balances_diff, prev_balances)
    balances_save(block_hash, after_balances)
    chaindata_latest_set(block_hash)
    return block


# --------------------------------------------------------------------
# Program entry

import json

# this just cleans up to make it look nice when printing
def clean_balances_dict(balances_dict):
    clean_keys = dict(('0x'+k.encode('hex'), \
            clean_balances_dict(balances_dict[k]) if type(balances_dict[k]) is dict else balances_dict[k]) \
            for k in balances_dict.keys())
    return clean_keys

def print_block(block_hash):
        latest_path = chaindata_path(block_hash, 'block')
        if os.path.exists(latest_path):
            block = Block.unmarshal(data_load(latest_path))
            print('====== BLOCK HASH: ',block.hash.encode('hex'), ' ======')
            signed_payment_arr = [SignedPayment.unmarshal(sign_pay) for sign_pay in payments_load(block.hash)]
            print("======== Latest block ========\n", block)
            print("======== Payments ========")
            [print(s_pay) for s_pay in signed_payment_arr]
            print("======== Balances ========\n", json.dumps(clean_balances_dict(balances_load(block.hash)),indent=2))
            return block


@click.command()
@click.option('--block', '-b', metavar="HASH", required=False, callback=arg_bytes32, help="Most recent block hash")
@click.option('--genesis', '-g', is_flag=True, help="Create new genesis block")
@click.option('--random', '-r', metavar='NUM', type=int, default=0, help="Create a block of random payments")
@click.argument('payments', nargs=-1, type=click.Path(exists=True))
@click.option('--latest', '-l', is_flag=True, help="Print latest state")
def main(block, genesis, random, payments, latest):
    if not chaindata_latest_get():
        if not genesis:
            raise ValueError("Must create genesis block")
        block_genesis()

    prev_hash = block
    if not prev_hash:
        prev_hash = chaindata_latest_get()

    # Create a block of random payments
    if random:
        signed_payments = random_payments(prev_hash, random)
        blockchain_apply(prev_hash, signed_payments)
        return 1

    # Otherwise, process payment files to create blocks
    for payment_file in payments:
        signed_payments = data_load(payment_file)
        block = blockchain_apply(prev_hash, signed_payments)
        prev_hash = block.hash


    # doart3 print latest state
    if latest:
        last_block_hash = chaindata_latest_get()
        block = print_block(last_block_hash)
        while block:
            print('\n\n')
            block = print_block(block.prev)

    return 1


if __name__ == "__main__":
    main()
