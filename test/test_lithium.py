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

from ion.lithium.etheventrelay import Lithium, pack_txn, pack_log
# from ion.lithium.etheventrelay import lithium_process_block_group, lithium_process_block, pack_txn, pack_log
from ion.utils import u256be

test_tx_hash = u'0x999999'
test_sender_addr = u'0x123456'
test_recipient_addr = u'0x678910'
test_tx_input = u'0x11111111'

class MockRPC():
    def port():
        return 8545

    def eth_blockNumber(self):
        return 1

    def eth_getBlockByNumber(self, block_number=1, tx_objects=True):
        return {u'transactions': [test_tx_hash]}

    def eth_getTransactionByHash(self, block_number=1):
        return {u'from': test_sender_addr, u'value': u'0x0', u'to': test_recipient_addr,
                u'input': test_tx_input}

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
        print("\nTest: Pack Transaction")
        rpc = MockRPC()
        tx = rpc.eth_getTransactionByHash()

        u256be_block_no_hex = u256be(rpc.eth_blockNumber()).encode('hex')
        pad_be_tx_value = zpad( int_to_big_endian( decode_int256( scan_bin( tx['value'] + ('0' * (len(tx['value']) % 2))))), 32).encode('hex')
        tx_input_hash = sha3( scan_bin( tx['input'] + ('0' * (len(tx['input']) % 2)))).encode('hex')

        packed_txn = pack_txn(rpc.eth_blockNumber(), tx).encode('hex')
        expected_result = '' + (test_sender_addr[2:]) + (test_recipient_addr[2:])

        self.assertTrue(packed_txn == expected_result)
        print("Test: Pack Transaction Success")


    def test_pack_log(self):
        print("\nTest: Pack Transaction Logs")
        rpc = MockRPC()
        tx = rpc.eth_getTransactionByHash()
        receipt = rpc.eth_getTransactionReceipt()
        self.assertTrue(len(receipt['logs']) > 0)
        for log in receipt['logs']:
            packed_log = pack_log(rpc.eth_blockNumber(), tx, log)

            address = scan_bin(log['address']).encode('hex')
            topic1 = scan_bin(log['topics'][1]).encode('hex')
            topic2 = scan_bin(log['topics'][2]).encode('hex')
            data = sha3(scan_bin(log['data'])).encode('hex')

            expected_result = '' + (test_sender_addr[2:]) + (test_recipient_addr[2:]) + address + topic1 + topic2

            self.assertTrue(packed_log.encode('hex') == expected_result)

        print("Test: Pack Transaction Logs Success")


    def test_process_block(self):
        print("\nTest: Process Single Block")
        lithium = Lithium()
        rpc = MockRPC()
        transfers = []
        items, tx_count, log_count = lithium.lithium_process_block(rpc, rpc.eth_blockNumber(), transfers)

        tx = rpc.eth_getTransactionByHash()
        receipt = rpc.eth_getTransactionReceipt()
        self.assertTrue(len(receipt['logs']) > 0)

        u256be_block_no_hex = u256be(rpc.eth_blockNumber()).encode('hex')
        pad_be_tx_value = zpad( int_to_big_endian( decode_int256( scan_bin( tx['value'] + ('0' * (len(tx['value']) % 2))))), 32).encode('hex')
        tx_input_hash = sha3( scan_bin( tx['input'] + ('0' * (len(tx['input']) % 2)))).encode('hex')

        for log in receipt['logs']:
            packed_log = pack_log(rpc.eth_blockNumber(), tx, log)

            address = scan_bin(log['address']).encode('hex')
            topic1 = scan_bin(log['topics'][1]).encode('hex')
            topic2 = scan_bin(log['topics'][2]).encode('hex')
            data = sha3(scan_bin(log['data'])).encode('hex')

        expected_txn_result = '' + (test_sender_addr[2:]) + (test_recipient_addr[2:])  + address + topic1 + topic2

        self.assertTrue(len(items) == 1)
        self.assertTrue(items[0].encode('hex') == expected_txn_result)
        self.assertTrue(tx_count == 1)
        self.assertTrue(log_count == 1)

        print("Test: Process Single Block Success")


    def test_process_block_group(self):
        print("\nTest: Process Block Group")
        lithium = Lithium()
        rpc = MockRPC()
        transfers = []
        items, group_tx_count, group_log_count, transfers = lithium.lithium_process_block_group(rpc, [1])

        tx = rpc.eth_getTransactionByHash()
        receipt = rpc.eth_getTransactionReceipt()
        self.assertTrue(len(receipt['logs']) > 0)

        u256be_block_no_hex = u256be(rpc.eth_blockNumber()).encode('hex')
        pad_be_tx_value = zpad( int_to_big_endian( decode_int256( scan_bin( tx['value'] + ('0' * (len(tx['value']) % 2))))), 32).encode('hex')
        tx_input_hash = sha3( scan_bin( tx['input'] + ('0' * (len(tx['input']) % 2)))).encode('hex')

        for log in receipt['logs']:
            packed_log = pack_log(rpc.eth_blockNumber(), tx, log)

            address = scan_bin(log['address']).encode('hex')
            topic1 = scan_bin(log['topics'][1]).encode('hex')
            topic2 = scan_bin(log['topics'][2]).encode('hex')
            data = sha3(scan_bin(log['data'])).encode('hex')

        expected_txn_result = '' + (test_sender_addr[2:]) + (test_recipient_addr[2:])  + address + topic1 + topic2

        self.assertTrue(len(items) == 1)
        self.assertTrue(items[0].encode('hex') == expected_txn_result)
        self.assertTrue(group_tx_count == 1)
        self.assertTrue(group_log_count == 1)

        print("Test: Process Block Group Success")



if __name__ == '__main__':
    import sys
    from os import path
    sys.path.append(path.dirname(path.dirname(path.abspath(__file__))))
    unittest.main()
