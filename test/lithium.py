## Copyright (c) 2016-2018 Clearmatics Technologies Ltd
## SPDX-License-Identifier: LGPL-3.0+

#!/usr/bin/env python
'''
Lithium tests

This tests the core components of the Lithium stack.
The Lithium module sits between a blockchain via RPC and Plasma, it takes data from the blockchain, hashes them together
and submits a merklised root of all the items to a validating plasma chain. Since RPC calls and merklisation are being tested
as separate components of the stack, the tests included here only pertain to Lithium-Specific operations of data
gathering and marshalling
'''

import unittest
from ethereum.utils import scan_bin, sha3, decode_int256, zpad, int_to_big_endian

test_tx_hash = u'0x999999'
test_sender_addr = u'0x123456'
test_recipient_addr = u'0x678910'
test_tx_input = u'0x11111111'

class MockRPC():
    def eth_blockNumber(self):
        return 1

    def eth_getBlockByNumber(self, block_number=1, tx_objects=True):
        return {u'transactions': [test_tx_hash]}

    def eth_getTransactionByHash(self, block_number=1):
        return {u'from': test_sender_addr, u'value': u'0x0', u'to': test_recipient_addr,
                u'input': test_tx_input}

    def eth_getTransactionReceipt(self, block_number=1):
        return {u'logs':
            [{u'data': u'0x0000000000000000000000000000000000000000000000000000000000003039', u'topics':
                  [u'0x4e08b679441f63ded9f3c75c680b7a081ccc985d0b60003f946bccc35c510ee7'],
              u'address': u'0x1313e0b8ee307b0258cc180312602e30e0b857cd'}]}


class LithiumTest(unittest.TestCase):

    def test_pack_txn(self):
        print("\nTest: Pack Transaction")
        rpc = MockRPC()
        tx = rpc.eth_getTransactionByHash()

        u256be_block_no_hex = u256be(rpc.eth_blockNumber()).encode('hex')
        pad_be_tx_value = zpad( int_to_big_endian( decode_int256( scan_bin( tx['value'] + ('0' * (len(tx['value']) % 2))))), 32).encode('hex')
        tx_input_hash = sha3( scan_bin( tx['input'] + ('0' * (len(tx['input']) % 2)))).encode('hex')

        packed_txn = pack_txn(rpc.eth_blockNumber(), tx).encode('hex')
        expected_result = '' + u256be_block_no_hex + (test_sender_addr[2:]) + (test_recipient_addr[2:]) + pad_be_tx_value + tx_input_hash

        self.assertTrue(packed_txn == expected_result)
        print("Test: Pack Transaction Success")


    def test_pack_log(self):
        print("\nTest: Pack Transaction Logs")
        rpc = MockRPC()
        receipt = rpc.eth_getTransactionReceipt()
        self.assertTrue(len(receipt['logs']) > 0)
        for log in receipt['logs']:
            packed_log = pack_log(rpc.eth_blockNumber(), log)

            u256be_block_no_hex = u256be(rpc.eth_blockNumber()).encode('hex')
            address = scan_bin(log['address']).encode('hex')
            topics = scan_bin(log['topics'][0]).encode('hex')
            data = sha3(scan_bin(log['data'])).encode('hex')

            expected_result = '' + u256be_block_no_hex + address + topics + data
            self.assertTrue(len(packed_log) == 1)
            self.assertTrue(packed_log[0].encode('hex') == expected_result)
        print("Test: Pack Transaction Logs Success")


    def test_process_block(self):
        print("\nTest: Process Single Block")
        rpc = MockRPC()
        items, tx_count, log_count = lithium_process_block(rpc, rpc.eth_blockNumber())

        tx = rpc.eth_getTransactionByHash()
        u256be_block_no_hex = u256be(rpc.eth_blockNumber()).encode('hex')
        pad_be_tx_value = zpad( int_to_big_endian( decode_int256( scan_bin( tx['value'] + ('0' * (len(tx['value']) % 2))))), 32).encode('hex')
        tx_input_hash = sha3( scan_bin( tx['input'] + ('0' * (len(tx['input']) % 2)))).encode('hex')

        expected_txn_result = '' + u256be_block_no_hex + (test_sender_addr[2:]) + (test_recipient_addr[2:]) + pad_be_tx_value + tx_input_hash

        self.assertTrue(len(items) == 2)
        self.assertTrue(items[0].encode('hex') == expected_txn_result)
        self.assertTrue(tx_count == 1)
        self.assertTrue(log_count == 1)

        print("Test: Process Single Block Success")


    def test_process_block_group(self):
        print("\nTest: Process Block Group")
        rpc = MockRPC()
        items, group_tx_count, group_log_count = lithium_process_block_group(rpc, [1])

        tx = rpc.eth_getTransactionByHash()
        u256be_block_no_hex = u256be(rpc.eth_blockNumber()).encode('hex')
        pad_be_tx_value = zpad( int_to_big_endian( decode_int256( scan_bin( tx['value'] + ('0' * (len(tx['value']) % 2))))), 32).encode('hex')
        tx_input_hash = sha3( scan_bin( tx['input'] + ('0' * (len(tx['input']) % 2)))).encode('hex')

        expected_txn_result = '' + u256be_block_no_hex + (test_sender_addr[2:]) + (test_recipient_addr[2:]) + pad_be_tx_value + tx_input_hash

        self.assertTrue(len(items) == 2)
        self.assertTrue(items[0].encode('hex') == expected_txn_result)
        self.assertTrue(group_tx_count == 1)
        self.assertTrue(group_log_count == 1)

        print("Test: Process Block Group Success")



if __name__ == '__main__':
    import sys
    from os import path
    sys.path.append(path.dirname(path.dirname(path.abspath(__file__))))
    from ion.etheventrelay import lithium_process_block_group, lithium_process_block, pack_txn, pack_log
    from ion.utils import u256be
    unittest.main()
