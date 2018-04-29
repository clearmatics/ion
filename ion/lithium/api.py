## Copyright (c) 2016-2018 Clearmatics Technologies Ltd
## SPDX-License-Identifier: LGPL-3.0+

#!/usr/bin/env python
"""
API

Provides a set of endpoints from which users can derive the key information regarding proofs
which is required when withdrawing funds from IonLock
"""

from flask import Flask, request, jsonify
from ion.merkle import merkle_tree, merkle_path

class LithiumRestApi(object):
    """
    Class containing the API for use with lithium
    """
    def __init__(self, lithium, host=None, port=None):
        # self.app = app
        self.app = Flask('REST API')
        self.host = host
        self.port = port
        self.lithium = lithium

    def serve_endpoints(self):
        if self.host is None:
            self.host = '127.0.0.1'
        if self.port is None:
            self.port = 5000

        # self.app = Flask("REST API")

        @self.app.route('/api/leaves', methods=['GET'])
        def api_leaves():
            """Return all the leaves for the merkle tree"""
            byte_leaves = self.lithium.leaves
            # print(byte_leaves)
            hex_leaves = [x.encode('hex') for x in byte_leaves]
            dict = {u'leaves': hex_leaves}
            return jsonify(dict)

        @self.app.route('/api/checkpoints', methods=['GET'])
        def api_checkpoint():
            """Returns the checkpoints"""
            byte_checkpoints = self.lithium.checkpoints

            index_map = [column[0] for column in byte_checkpoints]
            blockid = [column[1] for column in byte_checkpoints]

            hex_blockid = [format(x, 'x') for x in blockid]

            dict = {u'index': index_map, u'blockId': hex_blockid}
            return jsonify(dict)

        @self.app.route('/api/blockid', methods=['POST'])
        def api_blockid():
            """If passed a valid leaf returns corresponding ionlink blockId"""
            if request.method == 'POST':
                leaf = request.args.get('leaf')

            byte_leaves = self.lithium.leaves
            byte_checkpoints = self.lithium.checkpoints

            hex_leaves = [x.encode('hex') for x in byte_leaves]
            idx = hex_leaves[0].index(leaf)
            output = None
            for el in byte_checkpoints:
                if idx >= el[0]:
                    pass
                else:
                    output = el[1]
                    break

            dict = {u'blockId': format(output, 'x')}
            return jsonify(dict)

        @self.app.route('/api/proof', methods=['POST'])
        def api_proof():
            """If passed a valid leaf returns merkle proof"""
            if request.method == 'POST':
                leaf = request.args.get('leaf')

            byte_leaves = self.lithium.leaves
            tree, _ = merkle_tree(self.lithium.leaves)
            hex_leaf = leaf.decode('hex')

            path = merkle_path(hex_leaf, tree)
            hex_path = [format(x, 'x') for x in path]

            dict = {u'leaves': hex_path}
            return jsonify(dict)

        self.app.run(host=self.host, port=self.port)
