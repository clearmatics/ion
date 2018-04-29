## Copyright (c) 2016-2018 Clearmatics Technologies Ltd
## SPDX-License-Identifier: LGPL-3.0+

#!/usr/bin/env python
'''
API tests

This tests the API which gives users access to the data required to make a withdrawal.
'''

import unittest
from flask import Flask

from ion.lithium.etheventrelay import Lithium
from ion.lithium.api import LithiumRestApi

test_leaves = ['\x90\xf8\xbfjG\x9f2\x0e\xad\x07D\x11\xa4\xb0\xe7\x94N\xa8\xc9\xc1\x00',
        '\x17f\xd7\x07\xbb\xee\x15\xf79(\x9c\xec\xa4D\x82\xf8\xdf\x9f\xe0\xcd:b\xc0\xb8\x08Am)\x90u\xdb\xe7',
        '/l\xf3\xcd\xfe\xf8=;7\xaf\x0cN\xe7\xbe\xe6\x9e\xd3l\x11r\x87m\xcbq\xb8\x0cs[8\x02?\xdd',
        'D\xe0\xc6\xec>uWm\x99)\xdam\tRv\t\x91;\xaf\xf9\xe5\x8d\xa9(\xd5\xcas0N\xa4\x98\xc4',
        '\x90\xf8\xbfjG\x9f2\x0e\xad\x07D\x11\xa4\xb0\xe7\x94N\xa8\xc9\xc1\xc8\x9c\xe4sX\x82\xc9\xf0\xf0\xfe&hlS\x07N\t\xb0\xd5P'
        ]

CHECKPOINTS_URL = 'http://localhost:5000/api/checkpoints?leaf=ffcf8fdee72ac11b5c542428b35eef5769c409f0c89ce4735882c9f0f0fe26686c53074e09b0d550d833215cbcc3f914bd1c9ece3ee7bf8b14f841bb00000000000000000000000000000000000000000000000000000000000003e8e6a7765f46f721f4be3c5369983d84a6a86ed7c17bcea6b39876d1920c6532fa'
BASE_URL = 'http://1270.0.0.1:5000/'
BLOCKID_URL = 'http://localhost:5000/api/blockid?leaf=ffcf8fdee72ac11b5c542428b35eef5769c409f0c89ce4735882c9f0f0fe26686c53074e09b0d550d833215cbcc3f914bd1c9ece3ee7bf8b14f841bb00000000000000000000000000000000000000000000000000000000000003e8e6a7765f46f721f4be3c5369983d84a6a86ed7c17bcea6b39876d1920c6532fa'
PROOF_URL = 'http://localhost:5000/api/proof?leaf=ffcf8fdee72ac11b5c542428b35eef5769c409f0c89ce4735882c9f0f0fe26686c53074e09b0d550d833215cbcc3f914bd1c9ece3ee7bf8b14f841bb00000000000000000000000000000000000000000000000000000000000003e8e6a7765f46f721f4be3c5369983d84a6a86ed7c17bcea6b39876d1920c6532fa'

class MockLithium():
    def leaves():
        return test_leaves


class ApiTest(unittest.TestCase):
    def setUp(self):
        lithium = MockLithium()
        api = LithiumRestApi(lithium, '127.0.0.1', 5000)
        # api = LithiumRestApi(self, lithium, '127.0.0.1', 5000)
        self.app = api.app.test_client()
        self.app.testing = True


    def test_leaves2(self):
        print("Test: GET Leaves2")
        response = self.app.get(BASE_URL)
        print(response)
        self.assertEqual(response.status_code, 200)



if __name__ == '__main__':
    import sys
    from os import path
    sys.path.append(path.dirname(path.dirname(path.abspath(__file__))))
    unittest.main()
