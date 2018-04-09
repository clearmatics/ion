import unittest

from plasma.chain import chaindata_latest_get
from plasma.txpool import TxPool


class TestGenesisCreation(unittest.TestCase):

    def genesis_block(self):
        self.assertEqual('foo'.upper(), 'FOO')


class TestTXPool(unittest.TestCase):

    def new_tx_pool(self):
        block_hash = block_hash or chaindata_latest_get()
        require( block_hash is not None, "Must specify block hash" )
        pool = TxPool(block_hash)
        self.assertEqual(pool.payments.length, 0)
        self.assertEqual(pool.target, block_hash)

if __name__ == '__main__':
    unittest.main()
