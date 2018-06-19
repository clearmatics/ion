## Copyright (c) 2016-2018 Clearmatics Technologies Ltd
## SPDX-License-Identifier: LGPL-3.0+

#!/usr/bin/env python
'''
API tests

This tests the API which gives users access to the data required to make a withdrawal.
'''

import unittest
import requests
import json

from ion.lithium.api import app


CHECKPOINTS = \
    {
        u'81d9d8277b8f741b859de5455b9b56ff240d2ecf19101df3da9b76b137e5a7e6': 4,
        u'6ce75c011eac6f587c54493784ce2139b70e38b5b04fedab2bf5a84b500d0d92': 7
    }

LEAVES = \
    {
        u'leaves': [
            u'ffcf8fdee72ac11b5c542428b35eef5769c409f0c89ce4735882c9f0f0fe26686c53074e09b0d550d833215cbcc3f914bd1c9ece3ee7bf8b14f841bb03e8e6a7765f46f721f4be3c5369983d84a6a86ed7c17bcea6b39876d1920c6532fa',
            u'45b6d9232f9a2d8808fef6ee5339482aed37a8588a8668cb88fde5ffaab67ba1',
            u'51f8fcbbea5fb362345c94c4cd6809db199965fc22d7aa913a8ca73387fb5a30',
            u'73e34f24b39488c9d00971d23969be48654df4ad8e275c698ba4e88f4206b612',
            u'90f8bf6a479f320ead074411a4b0e7944ea8c9c1c89ce4735882c9f0f0fe26686c53074e09b0d550',
            u'90f8bf6a479f320ead074411a4b0e7944ea8c9c1c89ce4735882c9f0f0fe26686c53074e09b0d550',
            u'ffcf8fdee72ac11b5c542428b35eef5769c409f0c89ce4735882c9f0f0fe26686c53074e09b0d550d833215cbcc3f914bd1c9ece3ee7bf8b14f841bb03e8e6a7765f46f721f4be3c5369983d84a6a86ed7c17bcea6b39876d1920c6532fa']
    }

EXPECTED_BLOCKID = \
    {
        u'blockid': u'81d9d8277b8f741b859de5455b9b56ff240d2ecf19101df3da9b76b137e5a7e6'
    }

PROOF = \
        {
            u'proof': [
                u'97923772266235395715382770652917280357200452858347804453224054177914019790625',
                u'98712613182025294006941209131503858290191202453957227771980554773023716000070'
            ]
        }

class MockLithium():
    leaves = [b'\xff\xcf\x8f\xde\xe7*\xc1\x1b\\T$(\xb3^\xefWi\xc4\t\xf0\xc8\x9c\xe4sX\x82\xc9\xf0\xf0\xfe&hlS\x07N\t\xb0\xd5P\xd83!\\\xbc\xc3\xf9\x14\xbd\x1c\x9e\xce>\xe7\xbf\x8b\x14\xf8A\xbb\x03\xe8\xe6\xa7v_F\xf7!\xf4\xbe<Si\x98=\x84\xa6\xa8n\xd7\xc1{\xce\xa6\xb3\x98v\xd1\x92\x0ce2\xfa',
              b'E\xb6\xd9#/\x9a-\x88\x08\xfe\xf6\xeeS9H*\xed7\xa8X\x8a\x86h\xcb\x88\xfd\xe5\xff\xaa\xb6{\xa1',
              b'Q\xf8\xfc\xbb\xea_\xb3b4\\\x94\xc4\xcdh\t\xdb\x19\x99e\xfc"\xd7\xaa\x91:\x8c\xa73\x87\xfbZ0',
              b"s\xe3O$\xb3\x94\x88\xc9\xd0\tq\xd29i\xbeHeM\xf4\xad\x8e'\\i\x8b\xa4\xe8\x8fB\x06\xb6\x12",
              b'\x90\xf8\xbfjG\x9f2\x0e\xad\x07D\x11\xa4\xb0\xe7\x94N\xa8\xc9\xc1\xc8\x9c\xe4sX\x82\xc9\xf0\xf0\xfe&hlS\x07N\t\xb0\xd5P',
              b'\x90\xf8\xbfjG\x9f2\x0e\xad\x07D\x11\xa4\xb0\xe7\x94N\xa8\xc9\xc1\xc8\x9c\xe4sX\x82\xc9\xf0\xf0\xfe&hlS\x07N\t\xb0\xd5P',
              b'\xff\xcf\x8f\xde\xe7*\xc1\x1b\\T$(\xb3^\xefWi\xc4\t\xf0\xc8\x9c\xe4sX\x82\xc9\xf0\xf0\xfe&hlS\x07N\t\xb0\xd5P\xd83!\\\xbc\xc3\xf9\x14\xbd\x1c\x9e\xce>\xe7\xbf\x8b\x14\xf8A\xbb\x03\xe8\xe6\xa7v_F\xf7!\xf4\xbe<Si\x98=\x84\xa6\xa8n\xd7\xc1{\xce\xa6\xb3\x98v\xd1\x92\x0ce2\xfa']
    checkpoints = {
        u'81d9d8277b8f741b859de5455b9b56ff240d2ecf19101df3da9b76b137e5a7e6': 4,
        u'6ce75c011eac6f587c54493784ce2139b70e38b5b04fedab2bf5a84b500d0d92': 7
    }

class TestFlaskApi(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        pass

    @classmethod
    def tearDownClass(cls):
        pass

    def setUp(self):
        lithium = MockLithium()
        app.lithium = lithium
        self.app = app.test_client()

    def tearDown(self):
        pass

    def test_leaves(self):
        print("\nTest: /api/leaves Internals")
        response = self.app.get('/api/leaves')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.get_json(), LEAVES)

        blockid = '81d9d8277b8f741b859de5455b9b56ff240d2ecf19101df3da9b76b137e5a7e6'
        response = self.app.post('/api/leaves',json={'blockid': blockid})

        self.assertEqual(response.status_code, 200)
        leaves = response.get_json()[u'leaves']
        expected_leaves = LEAVES[u'leaves'][0:4]
        self.assertEqual(leaves, expected_leaves)

        print("Test: /api/leaves Internals Success")

    def test_checkpoints(self):
        print("\nTest: /api/checkpoints Internals")
        response = self.app.get('/api/checkpoints')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.get_json(), CHECKPOINTS)
        print("Test: /api/checkpoints Internals Success")

    def test_blockid(self):
        print("\nTest: /api/blockid Internals")
        value = '45b6d9232f9a2d8808fef6ee5339482aed37a8588a8668cb88fde5ffaab67ba1'
        response = self.app.post('/api/blockid', json={'leaf':value})
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.get_json(), EXPECTED_BLOCKID)
        print("Test: /api/blockid Internals Success")

    def test_proof(self):
        print("\nTest: /api/proof Internals")
        value = '45b6d9232f9a2d8808fef6ee5339482aed37a8588a8668cb88fde5ffaab67ba1'
        blockid = '81d9d8277b8f741b859de5455b9b56ff240d2ecf19101df3da9b76b137e5a7e6'
        response = self.app.post('/api/proof', json={'leaf': value, 'blockid': blockid})
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.get_json(), PROOF)
        print("Test: /api/blockid Internals Success")

    def test_proof_verification(self):
        print("\nTest: /api/verify Internals")
        value = '45b6d9232f9a2d8808fef6ee5339482aed37a8588a8668cb88fde5ffaab67ba1'
        blockid = '81d9d8277b8f741b859de5455b9b56ff240d2ecf19101df3da9b76b137e5a7e6'
        response = self.app.post('/api/proof', json={'leaf': value, 'blockid': blockid})
        proof = response.get_json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.get_json(), PROOF)

        proof = [int(x) for x in proof['proof']]
        response = self.app.post('/api/verify', json={'leaf': value, 'proof': proof, 'blockid': blockid})
        verified = response.get_json()[u'verified']
        assert verified

        print("Test: /api/blockid Internals Success")


if __name__ == '__main__':
    import sys
    from os import path
    sys.path.append(path.dirname(path.dirname(path.abspath(__file__))))
    unittest.main()
