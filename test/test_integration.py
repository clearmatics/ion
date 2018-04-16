#!/usr/bin/env python
'''
integration tests

This set of tests should test the locking of funds on chain A and the subsequent withdrawal
on chain B
'''

import click
import unittest
import socket
import subprocess

from click.testing import CliRunner
from ethereum.utils import scan_bin, sha3, decode_int256, zpad, int_to_big_endian

from ion.args import arg_bytes20, arg_ethrpc
from ion.merkle import merkle_tree, merkle_hash
from ion.utils import u256be

from ion.lithium.etheventrelay import iter_blocks,lithium_process_block_group

# Definition of the fundamental variables required
chainA  = "127.0.0.1:8545"
chainB  = "127.0.0.1:8546"
account = arg_bytes20(None, None, "0x90f8bf6a479f320ead074411a4b0e7944ea8c9c1")
lock    = arg_bytes20(None, None, "0xe982e462b094850f12af94d21d470e21be9d0e9c")
link    = arg_bytes20(None, None, "0xc89ce4735882c9f0f0fe26686c53074e09b0d550")
tokAddr = arg_bytes20(None, None, "0x9561c133dd8580860b6b7e504bc5aa500f0f06a7")
send    = arg_bytes20(None, None, "0xffcf8fdee72ac11b5c542428b35eef5769c409f0")
recv    = arg_bytes20(None, None, "0x22d491bde2303f2f43325b2108d26f1eaba1e32b")

class IntegrationTest(unittest.TestCase):

    def test_integration(self):
        print("\nTest: Integration")

        # Assert that both testrpc A and B are live
        sockA = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sockA.settimeout(2)
        self.assertFalse(sockA.connect_ex(('127.0.0.1',8545)), "Please run testrpc on 127.0.0.1:8545")
        sockB = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sockB.settimeout(2)
        self.assertFalse(sockB.connect_ex(('127.0.0.1',8546)), "Please run testrpc on 127.0.0.1:8546")

        # First iteration of block
        rpc_to   = arg_ethrpc(None, None, '127.0.0.1:8545')
        rpc_from = arg_ethrpc(None, None, '127.0.0.1:8546')
        ionlock  = rpc_from.proxy("abi/IonLock.abi", lock, account)

        print("Starting block iterator")
        print("latest Block: ", ionlock.LatestBlock())

        # Now deploy the tokens on the networks
        token  = rpc_from.proxy("abi/Token.abi", tokAddr, account)
        ## TODO: ensure that the starting balance of the accounts are as we expect...

        # Note this is the un-wound infinite loop of iter_blocks
        start = ionlock.LatestBlock()
        obh = min(start, max(1, rpc_from.eth_blockNumber() - 0))
        print("Starting block header: ", start)
        print("Previous block header: ", obh)
        self.assertEqual(start, 8)
        obh -= obh % 1
        blocks = []
        is_latest = False

        # Pack the blocks into batches
        bh = rpc_from.eth_blockNumber()
        for i in range(obh, bh):
            # XXX TODO: I think this is why the latest block info is not always in sync with geth
            if i == (bh - 1):
                is_latest = True
            blocks.append(i)
        obh = bh

        # Show the blocks to work on and then rename just to keep convention from etheventrelay.py
        print(blocks)
        block_group = blocks

        items, group_tx_count, group_log_count, transfers = lithium_process_block_group(rpc_from, block_group)
        print(items)
        print(group_tx_count)
        print(group_log_count)
        print(transfers)




        # is_latest, block_group in iter_blocks(rpc_from, ionlock.LatestBlock())

if __name__ == '__main__':
    import sys
    from os import path
    sys.path.append(path.dirname(path.dirname(path.abspath(__file__))))
    unittest.main()
