from __future__ import print_function

import click
import json
import time
import threading
import binascii
import random
import string
from ethereum.utils import scan_bin, sha3, decode_int256, zpad, int_to_big_endian, keccak

from ion.args import arg_bytes20, arg_ethrpc
from ion.merkle import merkle_tree, merkle_hash, merkle_path
from ion.utils import u256be
from chris import LithiumRESTAPI

on_transfer_signature = keccak.new(digest_bits=256).update('IonTransfer(address,address,uint256,bytes32,bytes)').hexdigest()

event_signatures = [on_transfer_signature]

def random_string(N):
    """
    Returns a random string to hash as pseudo data

    """
    return ''.join(random.SystemRandom().choice(string.ascii_uppercase + string.digits) for _ in range(N))


def jsonPrint(message, snippet):
    """
    Prints a message with a formatted slice of json which doesn't make your eyes bleed

    """
    print(message)
    print(json.dumps(snippet, indent=4, sort_keys=True))


def pack_txn(block_no, tx):
    """
    Packs all the information about a transaction into a deterministic fixed-sized array of bytes

        from || to

    """
    tx_from, tx_to, tx_value, tx_input = [scan_bin(x + ('0' * (len(x) % 2))) for x in [tx['from'], tx['to'], tx['value'], tx['input']]]
    tx_value = decode_int256(tx_value)

    return ''.join([
        tx_from,
        tx_to
    ])


def pack_log(block_no, tx, log):
    """
    Packs a log entry into one or more entries.

        from || to || address || topics[1] || topics[2]

    """

    return ''.join([
        scan_bin(tx['from']),
        scan_bin(tx['to']),
        scan_bin(log['address']),
        scan_bin(log['topics'][1]),
        scan_bin(log['topics'][2]),
    ])

def pack_items(items):
    """
    Ensures items has minimum of 4 leaves.

    """
    start = len(items)
    if start < 4:
        for val in range(start, 4):
            new_item = random_string(16)
            items.append(sha3(new_item))
    else:
        pass

def processProof(item, item_tree, rpc):
    file = open("./data/merklePath" + str(rpc.port) + ".txt", "w")
    path = merkle_path(item, item_tree)
    for item in path:
        file.write("%s\n" % item)

    file.close()

def processReference(ref, rpc):
    file = open("./data/reference" + str(rpc.port) + ".txt", "w")
    var = ref
    var = str(var)
    file.write(var)
    file.close()

def processLatestBlock(block, rpc):
    file = open("./data/latestBlock" + str(rpc.port) + ".txt", "w")
    file.write(str(block))
    file.close()



def lithium_process_block(rpc, block_height, transfers):
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

            tx_count += 1
            packed_txns = pack_txn(block_height, tx)
            item_value = packed_txns
            receipt = rpc.eth_getTransactionReceipt(tx_hash)
            transfer = False
            if len(receipt['logs']):
                for log_entry in receipt['logs']:
                    if log_entry['topics'][0][2:] in event_signatures:
                        print("Processing IonLock Transfer Event")
                        processReference(log_entry['topics'][2], rpc)
                        log_items = pack_log(block_height, tx, log_entry)
                        item_value = log_items
                        log_count += 1
                        transfer = True

            transfers.append(transfer)
            items.append(item_value)

    return items, tx_count, log_count


def lithium_process_block_group(rpc, block_group):
    """Process a group of blocks, returning the packed events and transactions"""
    print("Processing block group")
    items = []
    transfers = []
    group_tx_count = 0
    group_log_count = 0
    for block_height in block_group:
        block_items, tx_count, log_count = lithium_process_block(rpc, block_height, transfers)
        items += block_items
        group_tx_count += tx_count
        group_log_count += log_count

    return items, group_tx_count, group_log_count, transfers


def lithium_submit(batch, prev_root, rpc, link, account):
    """Submit batch of merkle roots to IonLink"""
    ionlink = rpc.proxy("abi/IonLink.abi", link, account)
    if not len(batch):
        return False

    current_block = batch[0][0]
    for pair in batch:
        if pair[2]:
            current_root = pair[1]
            ionlink.Update([prev_root, current_root])
            processLatestBlock(ionlink.GetLatestBlock(), rpc)
            prev_root = current_root

    return prev_root


def iter_blocks(run_event, rpc, start=1, group=1, backlog=0, interval=1):
    """Iterate through the block numbers"""
    obh = min(start, max(1, rpc.eth_blockNumber() - backlog))
    print("Starting block header: ", start)
    print("Previous block header: ", obh)
    obh -= obh % group
    blocks = []
    is_latest = False

    # Infinite loop event listener...
    while run_event.is_set():
        bh = rpc.eth_blockNumber() + 1
        for i in range(obh, bh):
            if i == (bh - 1):
                is_latest = True
            blocks.append(i)
            if len(blocks) % group == 0:
                yield is_latest, blocks
                blocks = []
                is_latest = False
        obh = bh

def lithiumInstance(run_event, rpc_from, rpc_to, from_account, to_account, lock, link, batch_size):
    ionlock = rpc_from.proxy("abi/IonLock.abi", lock, from_account)
    batch = []
    transfers = []
    prev_root = merkle_hash("merkle-tree-extra")

    print("Starting block iterator")
    print("Latest Block: ", ionlock.LatestBlock)

    for is_latest, block_group in iter_blocks(run_event, rpc_from, ionlock.LatestBlock()):
      items, group_tx_count, group_log_count, transfers = lithium_process_block_group(rpc_from, block_group)
      if len(items):
          pack_items(items)
          print("blocks %d-%d (%d tx, %d events)" % (min(block_group), max(block_group), group_tx_count, group_log_count))
          item_tree, root = merkle_tree(items)
          batch.append( (block_group[0], root, transfers[0]) )

          if transfers[0]==True:
              processProof(items[0], item_tree, rpc_from)

          if is_latest or len(batch) >= batch_size:
              print("Submitting batch of", len(batch), "blocks")
              prev_root = lithium_submit(batch, prev_root, rpc_to, link, to_account)
              batch = []
    return 0


@click.command(help="Ethereum event merkle tree relay daemon")
@click.option('--rpc-from', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Source Ethereum JSON-RPC server")
@click.option('--rpc-to', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8546', help="Destination Ethereum JSON-RPC server")
@click.option('--from-account', callback=arg_bytes20, metavar="0x...20", required=True, help="Sender")
@click.option('--to-account', callback=arg_bytes20, metavar="0x...20", required=True, help="Recipient")
@click.option('--lock', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLock contract address")
@click.option('--link', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLink contract address")
@click.option('--batch-size', type=int, default=32, metavar="N", help="Upload at most N items per transaction")
def threadedrelay(rpc_from, rpc_to, from_account, to_account, lock, link, batch_size):
    # Create new threads
    run_event = threading.Event()
    run_event.set()

    relay_to = threading.Thread(target = (lithiumInstance), args = (run_event, rpc_from, rpc_to, from_account, to_account, lock, link, batch_size))
    relay_to.start()

    try:
        while 1:
            time.sleep(.010)
    except KeyboardInterrupt:
        print("Attempting to close threads.")
        run_event.clear()
        relay_to.join()
        # relay_from.join()
        print( "threads successfully closed")




if __name__ == "__main__":
    import sys
    from os import path
    sys.path.append(path.dirname(path.dirname(path.abspath(__file__))))
    threadedrelay()
