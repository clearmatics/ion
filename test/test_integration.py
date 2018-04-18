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
import binascii

from click.testing import CliRunner
from ethereum.utils import scan_bin, sha3, decode_int256, zpad, int_to_big_endian

from ion.args import arg_bytes20, arg_ethrpc
from ion.merkle import merkle_tree, merkle_hash, merkle_path, merkle_proof
from ion.utils import u256be

from ion.lithium.etheventrelay import iter_blocks, lithium_process_block_group, lithium_submit

# Definition of the fundamental variables required
chainA  = "127.0.0.1:8545"
chainB  = "127.0.0.1:8546"
owner   = arg_bytes20(None, None, "0x90f8bf6a479f320ead074411a4b0e7944ea8c9c1")
lock    = arg_bytes20(None, None, "0xe982e462b094850f12af94d21d470e21be9d0e9c")
link    = arg_bytes20(None, None, "0xc89ce4735882c9f0f0fe26686c53074e09b0d550")
tokAddr = arg_bytes20(None, None, "0x9561c133dd8580860b6b7e504bc5aa500f0f06a7")
send    = arg_bytes20(None, None, "0xffcf8fdee72ac11b5c542428b35eef5769c409f0")
recv    = arg_bytes20(None, None, "0x22d491bde2303f2f43325b2108d26f1eaba1e32b")
dummy   = arg_bytes20(None, None, "0x1df62f291b2e969fb0849d99d9ce41e2f137006e")

class IntegrationTest(unittest.TestCase):

    def test_integration(self):
        print("\nTest: Integration")

        totalSupply_a   = 1000
        totalSupply_b   = 1000
        value_a         = 10
        value_b         = 10
        rawRef_a        = 'Hello world!'
        rawRef_b        = 'Hello world!'

        # Assert that both testrpc A and B are live
        sockA = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sockA.settimeout(2)
        self.assertFalse(sockA.connect_ex(('127.0.0.1',8545)), "Please run testrpc on 127.0.0.1:8545")
        sockB = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sockB.settimeout(2)
        self.assertFalse(sockB.connect_ex(('127.0.0.1',8546)), "Please run testrpc on 127.0.0.1:8546")

        # First iteration of block
        rpc_a       = arg_ethrpc(None, None, '127.0.0.1:8545')
        rpc_b       = arg_ethrpc(None, None, '127.0.0.1:8546')
        ionlock_a   = rpc_a.proxy("abi/IonLock.abi", lock, owner)
        ionlock_b   = rpc_b.proxy("abi/IonLock.abi", lock, owner)
        ionlink_a   = rpc_a.proxy("abi/IonLink.abi", link, owner)
        ionlink_b   = rpc_b.proxy("abi/IonLink.abi", link, owner)
        token_a     = rpc_a.proxy("abi/Token.abi", tokAddr, send)
        token_b     = rpc_b.proxy("abi/Token.abi", tokAddr, recv)
        batch_a     = []
        batch_b     = []
        transfers   = []
        prev_root   = merkle_hash("merkle-tree-extra")



################################################################################
        print("Starting block iterator: rpc_a")
        print("latest Block rpc_a: ", ionlock_a.LatestBlock())
        # Note this is the un-wound infinite loop of iter_blocks
        start = rpc_a.eth_blockNumber()


        # Now deploy the tokens on the networks
        token_a.mint(totalSupply_a)
        token_a.metadataTransfer(lock, value_a, rawRef_a)

        obh = min(start + 1, max(1, rpc_a.eth_blockNumber() - 0))
        print("Starting block header  rpc_a: ", start)
        print("Previous block header  rpc_a: ", obh)
        self.assertEqual(start, 8)
        obh -= obh % 1
        blocks = []
        is_latest = False


        # Pack the blocks into batches
        bh = rpc_a.eth_blockNumber() + 1
        print("Current block header rpc_a: ", bh)
        for i in range(obh, bh):
            # XXX TODO: I think this is why the latest block info is not always in sync with geth
            if i == (bh - 1):
                is_latest = True
            blocks.append(i)
        obh = bh

        # Show the blocks to work on and then rename just to keep convention from etheventrelay.py
        block_group = blocks
        self.assertEqual(len(blocks), 2)

        # Process the block on chain A
        items, group_tx_count, group_log_count, transfers = lithium_process_block_group(rpc_a, block_group)
        item_tree_a, root = merkle_tree(items)

        for i in range(0, len(block_group)):
            batch_a.append( (block_group[i], root, transfers[i]) )

        prev_root = lithium_submit(batch_a, prev_root, rpc_b, link, owner)
        latestBlockB = ionlink_b.GetLatestBlock()
        path_a       = merkle_path(items[1], item_tree_a)

        leafHashB    = sha3(items[1])
        reference_a = sha3(rawRef_a)

        ionlock_b   = rpc_b.proxy("abi/IonLock.abi", lock, send)
        # print(ionlink_b.Verify(latestBlockB, leafHashB, path_a))

################################################################################
        print("Starting block iterator: rpc_b")
        print("latest Block rpc_b: ", ionlock_b.LatestBlock())
        # Note this is the un-wound infinite loop of iter_blocks
        start = rpc_b.eth_blockNumber()


        # Now deploy the tokens on the networks
        token_b.mint(totalSupply_b)
        token_b.metadataTransfer(lock, value_b, rawRef_b)

        obh = min(start + 1, max(1, rpc_b.eth_blockNumber() - 0))
        print("Starting block header  rpc_b: ", start)
        print("Previous block header  rpc_b: ", obh)
        self.assertEqual(start, 9)
        obh -= obh % 1
        blocks = []
        is_latest = False


        # Pack the blocks into batches
        bh = rpc_b.eth_blockNumber() + 1
        print("Current block header rpc_b: ", bh)
        for i in range(obh, bh):
            # XXX TODO: I think this is why the latest block info is not always in sync with geth
            if i == (bh - 1):
                is_latest = True
            blocks.append(i)
        obh = bh

        # Show the blocks to work on and then rename just to keep convention from etheventrelay.py
        block_group = blocks
        self.assertEqual(len(blocks), 2)


        # Process the block on chain B
        items, group_tx_count, group_log_count, transfers = lithium_process_block_group(rpc_b, block_group)

        item_tree_b, root = merkle_tree(items)

        for i in range(0, len(block_group)):
            batch_b.append( (block_group[i], root, transfers[i]) )

        prev_root = lithium_submit(batch_b, prev_root, rpc_a, link, owner)
        latestBlockA = ionlink_a.GetLatestBlock()
        path_b       = merkle_path(items[1], item_tree_b)

        leafHashA    = sha3(items[1])
        reference_b = sha3(rawRef_b)

        ionlock_a   = rpc_a.proxy("abi/IonLock.abi", lock, recv)

################################################################################

        self.assertEqual(token_a.balanceOf(lock), value_a)
        self.assertEqual(token_b.balanceOf(lock), value_b)
        self.assertEqual(token_a.balanceOf(recv), 0)
        self.assertEqual(token_b.balanceOf(send), 0)

        ionlock_b.Withdraw(value_a, reference_a, latestBlockB, path_a)
        ionlock_a.Withdraw(value_b, reference_b, latestBlockA, path_b)
        self.assertEqual(token_a.balanceOf(lock), 0)
        self.assertEqual(token_b.balanceOf(lock), 0)
        self.assertEqual(token_a.balanceOf(recv), value_a)
        self.assertEqual(token_b.balanceOf(send), value_b)

if __name__ == '__main__':
    import sys
    from os import path
    sys.path.append(path.dirname(path.dirname(path.abspath(__file__))))
    unittest.main()
