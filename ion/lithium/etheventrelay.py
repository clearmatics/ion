# Copyright (c) 2016-2018 Clearmatics Technologies Ltd
# SPDX-License-Identifier: LGPL-3.0+
"""
etheventrelay: manages the update of Ionlink and Merkle proof for corresponding chain
"""
from __future__ import print_function

import threading
import random
import string
import click
from ethereum.utils import scan_bin, sha3, keccak

from ion.args import arg_bytes20, arg_ethrpc
from ion.merkle import merkle_tree, merkle_hash

from ion.lithium.api import LithiumRestApi

TRANSFER_SIGNATURE = keccak.new(digest_bits=256) \
    .update('IonTransfer(address,address,uint256,bytes32,bytes)') \
    .hexdigest()

EVENT_SIGNATURES = [TRANSFER_SIGNATURE]


def random_string(amount):
    """
    Returns a random string to hash as pseudo data
    """
    return ''.join(random.SystemRandom() \
        .choice(string.ascii_uppercase + string.digits) for _ in range(amount))


def pack_txn(txn):
    """
    Packs all the information about a transaction into a deterministic fixed-sized array of bytes
        from || to
    """
    tx_from, tx_to, tx_value, tx_input = [scan_bin(x + ('0' * (len(x) % 2))) \
        for x in [txn['from'], txn['to'], txn['value'], txn['input']]]

    return ''.join([
        tx_from,
        tx_to
    ])


def pack_log(txn, log):
    """
    Packs a log entry into one or more entries.
        from || to || address || topics[1] || topics[2]
    """

    return ''.join([
        scan_bin(txn['from']),
        scan_bin(txn['to']),
        scan_bin(log['address']),
        scan_bin(log['topics'][1]),
        scan_bin(log['topics'][2]),
    ])


def pack_items(items):
    """
    Ensures items has minimum of 4 leaves
    """
    start = len(items)
    if start < 4:
        for _ in range(start, 4):
            new_item = random_string(16)
            items.append(sha3(new_item))
    else:
        pass


class Lithium(object):
    """
    Lithium process the blocks for the event relat to identify the IonLock transactions which occur
    on the rpc_from chains, which are then added to the IonLink of the rpc_to chain.
    """
    def __init__(self):
        self.checkpoints = []
        self.leaves = []
        self._run_event = None
        self._relay_to = None

    def lithium_process_block(self, rpc, block_height, transfers):
        """Returns all items within the block"""
        block = rpc.eth_getBlockByNumber(block_height, False)
        items = []
        log_count = 0
        tx_count = 0
        if block['transactions']:
            for tx_hash in block['transactions']:
                transaction = rpc.eth_getTransactionByHash(tx_hash)
                if transaction['to'] is None:
                    continue

                tx_count += 1
                packed_txns = pack_txn(transaction)
                item_value = packed_txns
                receipt = rpc.eth_getTransactionReceipt(tx_hash)
                transfer = False
                if receipt['logs']:
                    for log_entry in receipt['logs']:
                        if log_entry['topics'][0][2:] in EVENT_SIGNATURES:
                            print("Processing IonLock Transfer Event")
                            log_items = pack_log(transaction, log_entry)
                            item_value = log_items
                            log_count += 1
                            transfer = True

                transfers.append(transfer)
                items.append(item_value)

        return items, tx_count, log_count


    def lithium_process_block_group(self, rpc, block_group):
        """Process a group of blocks, returning the packed events and transactions"""
        print("Processing block group")
        items = []
        transfers = []
        group_tx_count = 0
        group_log_count = 0
        for block_height in block_group:
            block_items, tx_count, log_count = self.lithium_process_block(rpc, block_height, transfers)
            items += block_items
            group_tx_count += tx_count
            group_log_count += log_count

        return items, group_tx_count, group_log_count, transfers


    def iter_blocks(self, run_event, rpc, start=1, group=1, backlog=0, interval=1):
        """Iterate through the block numbers"""
        old_head = min(start, max(1, rpc.eth_blockNumber() - backlog))
        print("Starting block header: ", start)
        print("Previous block header: ", old_head)
        old_head -= old_head % group
        blocks = []
        is_latest = False

        # Infinite loop event listener...
        while run_event.is_set():
            head = rpc.eth_blockNumber() + 1
            for i in range(old_head, head):
                if i == (head - 1):
                    is_latest = True
                blocks.append(i)
                if len(blocks) % group == 0:
                    yield is_latest, blocks
                    blocks = []
                    is_latest = False
            old_head = head


    def lithium_submit(self, batch, prev_root, rpc, link, account, checkpoints, nleaves):
        """Submit batch of merkle roots to IonLink"""
        ionlink = rpc.proxy("abi/IonLink.abi", link, account)
        if not batch:
            return False

        current_block = batch[0][0]
        print(len(batch))

        for pair in batch:
            if pair[2]:
                current_root = pair[1]
                ionlink.Update([prev_root, current_root])
                ionlink_latest = ionlink.GetLatestBlock()
                checkpoints.append((nleaves, ionlink_latest))
                prev_root = current_root

        return prev_root


    def lithium_instance(self, run_event, rpc_from, rpc_to, from_account, to_account, lock,
                         link, batch_size):
        ionlock = rpc_from.proxy("abi/IonLock.abi", lock, from_account)
        batch = []
        transfers = []

        prev_root = merkle_hash("merkle-tree-extra")

        print("Starting block iterator")
        print("Latest Block: ", ionlock.LatestBlock)

        for is_latest, block_group in self.iter_blocks(run_event, rpc_from, ionlock.LatestBlock()):
            items, group_tx_count, group_log_count, transfers = self.lithium_process_block_group(rpc_from, block_group)
            if items:
                for value in items:
                    self.leaves.append(value)

                print(self.leaves)
                pack_items(self.leaves)
                print("blocks %d-%d (%d tx, %d events)" % (min(block_group), max(block_group), group_tx_count, group_log_count))
                _, root = merkle_tree(self.leaves)
                batch.append((block_group[0], root, transfers[0]))

                if is_latest or len(batch) >= batch_size:
                    print("Submitting batch of", len(batch), "blocks")
                    prev_root = self.lithium_submit(batch, prev_root, rpc_to, link, to_account, self.checkpoints, len(self.leaves))
                    batch = []
        return 0

    def run(self, rpc_from, rpc_to, from_account, to_account, lock, link, batch_size):
        """ Launches the etheventrelay on a thread"""
        self._run_event = threading.Event()
        self._run_event.set()

        self._relay_to = threading.Thread(target=(self.lithium_instance), \
            args=(self._run_event, rpc_from, rpc_to, from_account, to_account, lock, link, batch_size))
        self._relay_to.start()

    def stop(self):
        """ Stops the etheventrelay thread """
        self._run_event.clear()
        self._relay_to.join()


@click.command(help="Ethereum event merkle tree relay daemon")
@click.option('--rpc-from', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', \
              help="Source Ethereum JSON-RPC server")
@click.option('--rpc-to', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8546', \
              help="Destination Ethereum JSON-RPC server")
@click.option('--from-account', callback=arg_bytes20, metavar="0x...20", required=True, help="Sender")
@click.option('--to-account', callback=arg_bytes20, metavar="0x...20", required=True, help="Recipient")
@click.option('--lock', callback=arg_bytes20, metavar="0x...20", required=True,
              help="IonLock contract address")
@click.option('--link', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLink contract address")
@click.option('--batch-size', type=int, default=32, metavar="N", help="Upload at most N items per transaction")
def etheventrelay(rpc_from, rpc_to, from_account, to_account, lock, link, batch_size):
    lithium = Lithium()
    api = LithiumRestApi(lithium)
    lithium.run(rpc_from, rpc_to, from_account, to_account, lock, link, batch_size)
    api.serve_endpoints()
    lithium.stop()


if __name__ == "__main__":
    import sys
    from os import path
    sys.path.append(path.dirname(path.dirname(path.abspath(__file__))))
