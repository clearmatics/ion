#!/usr/bin/env python
from __future__ import print_function
import time
import sys
import argparse
from ethereum.utils import scan_bin, sha3, decode_int256, zpad, int_to_big_endian

from .utils import u256be
from .merkle import merkle_tree
from .solproxy import solproxy

from .args import Bytes20Action, EthRpcAction


def pack_txn(block_no, tx):
    """
    Packs all the information about a transaction into a deterministic fixed-sized array of bytes

        block_no || from || to || value || sha3(input)

    Where `value` is expanded to a 256bit big endian integer.
    """
    tx_from, tx_to, tx_value, tx_input = [scan_bin(x + ('0' * (len(x) % 2)))
                                          for x in [tx['from'], tx['to'], tx['value'], tx['input']]]
    tx_value = decode_int256(tx_value)
    return ''.join([
        u256be(block_no),
        tx_from,
        tx_to,
        zpad(int_to_big_endian(tx_value), 32),
        sha3(tx_input)
    ])


def pack_log(block_no, log):
    """
    Packs a log entry into one or more entries.

        block_no || address || topic || sha3(data)

    Where topic is the SHA3 of the event signature, e.g. `OnDeposit(bytes32)`
    """
    if not len(log['topics']):
        return []
    return [''.join([
        u256be(block_no),
        scan_bin(log['address']),
        scan_bin(log['topics'][0]),
        sha3(scan_bin(log['data']))
    ])]


def iter_blocks(c, start=1, group=1, backlog=0, interval=1):
    """Iterate through the block numbers"""
    obh = min(start, max(1, c.eth_blockNumber() - backlog))
    obh -= obh % group
    blocks = []
    is_latest = False
    while True:
        bh = c.eth_blockNumber()
        for i in range(obh, bh):
            if i == (bh - 1):
                is_latest = True
            blocks.append(i)
            if len(blocks) % group == 0:
                yield is_latest, blocks
                blocks = []
                is_latest = False
        obh = bh
        try:
            time.sleep(interval)
        except KeyboardInterrupt:
            raise StopIteration


def lithium_options():
    """Parse commandline options"""
    parser = argparse.ArgumentParser(description='Lithium merkle tree builder')
    parser.add_argument('--from', dest='rpc_from', action=EthRpcAction, required=True,
                        help='Source Ethereum RPC address, ip:port')
    parser.add_argument('--to', dest='rpc_to', action=EthRpcAction, required=True,
                        help='Destination Ethereum RPC address, ip:port')
    parser.add_argument('--account', dest='account', required=True, action=Bytes20Action,
                        help='Ethereum account address, 0x....')
    parser.add_argument('--contract', dest='contract', required=True, action=Bytes20Action,
                        help='Sodium contract address, 0x....')
    parser.add_argument('--batch-size', dest='batchsize', default=32, type=int,
                        help='Maximum number of merkle roots to submit to Sodium at once')
    opt = parser.parse_args()
    opt.sodium = solproxy(opt.rpc_to, "abi/Sodium.abi", opt.contract, opt.account)

    return opt


def lithium_process_block(rpc, block_height):
    """Returns all items within the block"""
    block = rpc.eth_getBlockByNumber(block_height, False)
    items = []
    log_count = 0
    tx_count = 0
    if len(block['transactions']):
        for tx_hash in block['transactions']:
            tx = rpc.eth_getTransactionByHash(tx_hash)
            if tx['to'] is None:
                continue
            items.append(pack_txn(block_height, tx))
            tx_count += 1
            receipt = rpc.eth_getTransactionReceipt(tx_hash)
            if len(receipt['logs']):
                for log_entry in receipt['logs']:
                    log_items = pack_log(block_height, log_entry)
                    items += log_items
                    log_count += len(log_items)
    return items, tx_count, log_count


def lithium_process_block_group(rpc, block_group):
    """Process a group of blocks, returning the packed events and transactions"""
    items = []
    group_tx_count = 0
    group_log_count = 0
    for block_height in block_group:
        block_items, tx_count, log_count = lithium_process_block(rpc, block_height)
        items += block_items
        group_tx_count += tx_count
        group_log_count += log_count
    return items, group_tx_count, group_log_count


def lithium_submit(opt, batch):
    """Submit batch of merkle roots to Sodium"""
    if not len(batch):
        return False
    start_block = batch[0][0]
    roots = [pair[1] for pair in batch]
    opt.sodium.Update(start_block, roots)
    return True


def lithium_loop(opt):
    Na = opt.sodium
    batch = []
    for is_latest, block_group in iter_blocks(opt.rpc_from, Na.NextBlock(), Na.GroupSize()):
        items, group_tx_count, group_log_count = lithium_process_block_group(opt.rpc_from, block_group)
        if len(items):
            print("blocks %d-%d (%d tx, %d events)" % (min(block_group), max(block_group), group_tx_count, group_log_count))

            _, root = merkle_tree(items)
            batch.append( (block_group[0], root) )

            if is_latest or len(batch) >= opt.batchsize:
                print("submitting batch of", len(batch), "blocks")
                lithium_submit(opt, batch)
                batch = []
    return 0


def main():
    opt = lithium_options()
    return lithium_loop(opt)


if __name__ == "__main__":
    sys.exit(main())
