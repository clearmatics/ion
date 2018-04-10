from __future__ import print_function

import click
import time
from ethereum.utils import scan_bin, sha3, decode_int256, zpad, int_to_big_endian, keccak

from ion.args import arg_bytes20, arg_ethrpc
from ion.merkle import merkle_tree
from ion.utils import u256be

on_transfer_signature = keccak.new(digest_bits=256).update('IonTransfer(address,address,uint256,bytes32)').hexdigest()

event_signatures = [on_transfer_signature]

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

    print("DATA")
    print(log['data'])
    print("TOPICS")
    print(log['topics'])
    if not len(log['topics']):
        return []
    return [''.join([
        u256be(block_no),
        scan_bin(log['address']),
        scan_bin(log['topics'][0]),
        sha3(scan_bin(log['data']))
    ])]


def iter_blocks(rpc, start=1, group=1, backlog=0, interval=1):
    """Iterate through the block numbers"""
    obh = min(start, max(1, rpc.eth_blockNumber() - backlog))
    print(start)
    print(obh)
    obh -= obh % group
    blocks = []
    is_latest = False
    while True:
        bh = rpc.eth_blockNumber()
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
            print("Stopped by Keyboard interrupt")
            raise StopIteration

def lithium_process_ionlock_transfer_event(log):
    pass

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
            packed_txns = pack_txn(block_height, tx)
            items.append(packed_txns)
            tx_count += 1
            receipt = rpc.eth_getTransactionReceipt(tx_hash)
            if len(receipt['logs']):
                for log_entry in receipt['logs']:
                    if log_entry['topics'][0][2:] in event_signatures:
                        print("Processing IonLock Transfer Event")
                        log_items = pack_log(block_height, log_entry)
                        items += log_items
                        log_count += len(log_items)
    return items, tx_count, log_count


def lithium_process_block_group(rpc, block_group):
    """Process a group of blocks, returning the packed events and transactions"""
    print("Processing block group")
    items = []
    group_tx_count = 0
    group_log_count = 0
    for block_height in block_group:
        block_items, tx_count, log_count = lithium_process_block(rpc, block_height)
        items += block_items
        group_tx_count += tx_count
        group_log_count += log_count
    return items, group_tx_count, group_log_count


def lithium_submit(batch):
    """Submit batch of merkle roots to Beryllium"""
    if not len(batch):
        return False
    start_block = batch[0][0]
    roots = [pair[1] for pair in batch]
    print("Start block:\n", start_block)
    print("Roots:\n", roots)
    # sodium.Update(start_block, roots)
    print("NEED TO SUBMIT TO BERYLLIUM")
    return True


@click.command(help="Ethereum event merkle tree relay daemon")
@click.option('--rpc-from', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Source Ethereum JSON-RPC server")
@click.option('--rpc-to', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Destination, where contract is")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Pays for Gas")
@click.option('--contract', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLock contract address")
@click.option('--batch-size', type=int, default=32, metavar="N", help="Upload at most N items per transaction")
def etheventrelay(rpc_from, rpc_to, account, contract, batch_size):
    ionlock = rpc_to.proxy("abi/IonLock.abi", contract, account)
    batch = []

    print("Starting block iterator")
    print(ionlock.LatestBlock())    
    for is_latest, block_group in iter_blocks(rpc_from, ionlock.LatestBlock()):
        items, group_tx_count, group_log_count = lithium_process_block_group(rpc_from, block_group)
        if len(items):
            print("blocks %d-%d (%d tx, %d events)" % (min(block_group), max(block_group), group_tx_count, group_log_count))
            print(items)
            item_tree, root = merkle_tree(items)
            batch.append( (block_group[0], root) )

            if is_latest or len(batch) >= batch_size:
                print("submitting batch of", len(batch), "blocks")
                lithium_submit(batch)
                batch = []
    return 0


if __name__ == "__main__":
    import sys
    from os import path
    sys.path.append(path.dirname(path.dirname(path.abspath(__file__))))
    etheventrelay()
