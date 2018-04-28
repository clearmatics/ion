#!/usr/bin/env python
import os
import json
import binascii

from flask import Flask, request, jsonify
from ion.merkle import merkle_tree, merkle_path


class LithiumRestApi(object):
    def __init__(self, lithium, host=None, port=None):
        self.host = host
        self.port = port
        self.lithium = lithium

    def serve_endpoints(self):
        if self.host is None:
            self.host = '127.0.0.1'
        if self.port is None:
            self.port = 5000

        app = Flask("REST API")

        @app.route('/api/leaves', methods=['GET'])
        def api_leaves():
            """Return all the leaves for the merkle tree"""
            byte_leaves = self.lithium.leaves
            hex_leaves = map(lambda x:x.encode('hex'), byte_leaves)
            dict = {u'leaves': hex_leaves}
            return jsonify(dict)

        @app.route('/api/checkpoints', methods=['GET'])
        def api_checkpoint():
            """Returns the checkpoints"""
            byte_checkpoints = self.lithium.checkpoints

            index = [column[0] for column in byte_checkpoints]
            blockid = [column[1] for column in byte_checkpoints]
            index_map = map(lambda x:x, index)
            hex_blockid = map(lambda x:format(x, 'x'), blockid)
            dict = {u'index': index_map, u'blockId': hex_blockid }
            return jsonify(dict)

        @app.route('/api/blockid', methods=['POST'])
        def api_blockid():
            """If passed a valid leaf returns corresponding ionlink blockId"""
            if request.method == 'POST':
                leaf = request.args.get('leaf')

            byte_leaves = self.lithium.leaves
            byte_checkpoints = self.lithium.checkpoints

            hex_leaves = map(lambda x:x.encode('hex'), byte_leaves)
            idx = hex_leaves[0].index(leaf)
            output = None
            for el in byte_checkpoints:
                if idx>=el[0]:
                    pass
                else:
                    output = el[1]
                    break

            dict = {u'blockId': format(output, 'x') }
            return jsonify(dict)


        @app.route('/api/proof', methods=['POST'])
        def api_proof():
            """If passed a valid leaf returns merkle proof"""
            if request.method == 'POST':
                leaf = request.args.get('leaf')

            byte_leaves = self.lithium.leaves
            tree, _ = merkle_tree(self.lithium.leaves)
            hex_leaf = leaf.decode('hex')

            path = merkle_path(hex_leaf, tree)
            hex_path = map(lambda x:format(x, "x"), path)
            dict = {u'leaves': hex_path }
            return jsonify(dict)

        app.run(host=self.host, port=self.port)
