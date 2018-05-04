## Copyright (c) 2016-2018 Clearmatics Technologies Ltd
## SPDX-License-Identifier: LGPL-3.0+

#!/usr/bin/env python
'''
lithium tests

This tests the core components of the Lithium stack.
The Lithium module sits between a blockchain via RPC and Plasma, it takes data from the blockchain, hashes them together
and submits a merklised root of all the items to a validating plasma chain. Since RPC calls and merklisation are being tested
as separate components of the stack, the tests included here only pertain to Lithium-Specific operations of data
gathering and marshalling
'''

import unittest
from ethereum.utils import scan_bin, sha3, decode_int256, zpad, int_to_big_endian
from ion.lithium.lithium import Lithium, pack_txn, pack_log

test_tx_hash = u'0x999999'
test_sender_addr = u'0x123456'
test_recipient_addr = u'0x678910'
test_input = u'0x11111111'
test_value = u'0xfffd'

class MockRPC():
    def port():
        return 8545

    def eth_blockNumber(self):
        return 1

    def eth_getBlockByNumber(self, block_number=1, tx_objects=True):
        return {u'transactions': [test_tx_hash]}

    def eth_getTransactionByHash(self, block_number=1):
        return {u'from': test_sender_addr, u'value': test_value, u'to': test_recipient_addr,
                u'input': test_input}

    def eth_getTransactionReceipt(self, block_number=1):
        json = {u'logs':
                    [{
                        u'type': u'mined',
                        u'blockHash': u'0xdbb21e6784da1a2631a331be8bed3f94a738951c090e011cb74bbd80efc3c08d',
                        u'transactionHash': u'0x16414feeb07b3fc7f5eb343276ad016cbfbe3c98d1343ed58758d024ba5ae824',
                        u'data': u'0x000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000573706f6f6e000000000000000000000000000000000000000000000000000000',
                        u'topics': [
                                    u'0x2ff272db005d5490d7b8cc4833b3ae7018bc1ab8160c0253564c1c0724c8962d', u'0x00000000000000000000000090f8bf6a479f320ead074411a4b0e7944ea8c9c1', u'0x000000000000000000000000e982e462b094850f12af94d21d470e21be9d0e9c'
                                    ],
                        u'blockNumber': u'0x1c',
                        u'address': u'0x9561c133dd8580860b6b7e504bc5aa500f0f06a7',
                        u'logIndex': u'0x02',
                        u'transactionIndex': u'0x00'
                    }]
                }
        return json


class LithiumTest(unittest.TestCase):

    def test_pack_txn(self):
        print("\n==== Test: Pack Transaction ====")
        rpc = MockRPC()
        txn = rpc.eth_getTransactionByHash()

        packed_txn = pack_txn(txn).encode('hex')
        expected_result = '' + (test_sender_addr[2:]) + (test_recipient_addr[2:])

        self.assertEqual(packed_txn, expected_result)
        print("Test: Pack Transaction Success")


    def test_pack_log(self):
        print("\n==== Test: Pack Transaction Logs ====")
        rpc = MockRPC()
        txn = rpc.eth_getTransactionByHash()
        receipt = rpc.eth_getTransactionReceipt()

        self.assertTrue(len(receipt['logs']) > 0)

        for log in receipt['logs']:
            packed_log = pack_log(txn, log)

            address = scan_bin(log['address']).encode('hex')
            topic1 = scan_bin(log['topics'][1]).encode('hex')
            topic2 = scan_bin(log['topics'][2]).encode('hex')

            expected_result = '' + (test_sender_addr[2:]) + (test_recipient_addr[2:]) + address + topic1 + topic2

            self.assertEqual(packed_log.encode('hex'), expected_result)

        print("Test: Pack Transaction Logs Success")


    def test_process_block(self):
        print("\n==== Test: Process Single Block ====")
        lithium = Lithium()
        rpc = MockRPC()
        transfers = []
        items, tx_count, log_count = lithium.process_block(rpc, rpc.eth_blockNumber(), transfers)

        txn = rpc.eth_getTransactionByHash()
        receipt = rpc.eth_getTransactionReceipt()

        self.assertEqual(len(receipt['logs']), 1)

        log = receipt['logs'][0]
        packed_log = pack_log(txn, log)

        address = scan_bin(log['address']).encode('hex')
        topic1 = scan_bin(log['topics'][1]).encode('hex')
        topic2 = scan_bin(log['topics'][2]).encode('hex')

        expected_txn_result = '' + (test_sender_addr[2:]) + (test_recipient_addr[2:])  + address + topic1 + topic2
        self.assertEqual(packed_log.encode('hex'), expected_txn_result)

        self.assertEqual(len(items), 1)
        self.assertEqual(items[0].encode('hex'), expected_txn_result)
        self.assertEqual(tx_count, 1)
        self.assertEqual(log_count, 1)

        print("Test: Process Single Block Success")


    def test_process_block_group(self):
        print("\n==== Test: Process Block Group ====")
        lithium = Lithium()
        rpc = MockRPC()
        items, group_tx_count, group_log_count, transfers = lithium.process_block_group(rpc, [1])

        txn = rpc.eth_getTransactionByHash()
        receipt = rpc.eth_getTransactionReceipt()

        self.assertEqual(len(receipt['logs']), 1)

        log = receipt['logs'][0]
        packed_log = pack_log(txn, log)

        address = scan_bin(log['address']).encode('hex')
        topic1 = scan_bin(log['topics'][1]).encode('hex')
        topic2 = scan_bin(log['topics'][2]).encode('hex')

        expected_txn_result = '' + (test_sender_addr[2:]) + (test_recipient_addr[2:])  + address + topic1 + topic2
        self.assertEqual(packed_log.encode('hex'), expected_txn_result)

        self.assertEqual(len(items), 1)
        self.assertEqual(items[0].encode('hex'), expected_txn_result)
        self.assertEqual(group_tx_count, 1)
        self.assertEqual(group_log_count, 1)

        print("Test: Process Block Group Success")



if __name__ == '__main__':
    import sys
    from os import path
    sys.path.append(path.dirname(path.dirname(path.abspath(__file__))))

    unittest.main()
