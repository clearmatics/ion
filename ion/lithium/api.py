## Copyright (c) 2016-2018 Clearmatics Technologies Ltd
## SPDX-License-Identifier: LGPL-3.0+

#!/usr/bin/env python
"""
API

Provides a set of endpoints from which users can derive the key information regarding proofs
which is required when withdrawing funds from IonLock
"""

import os
from flask import Flask, request, jsonify
from flask_restful import Resource, Api

# from flask import Flask, url_for
from ion.merkle import merkle_tree, merkle_path, merkle_proof, merkle_hash


app = Flask(__name__)

@app.route('/')
def index():
    return 'index'

@app.route('/api/leaves', methods=['GET', 'POST'])
def api_leaves():
    """Return all the leaves for the merkle tree"""

    if request.method == 'POST':
        json = request.get_json()
        blockid = json[u'blockid']
        if blockid is not None:
            nleaves = app.lithium.checkpoints[blockid]
            byte_leaves = app.lithium.leaves[0:nleaves]
            hex_leaves = [x.encode('hex') for x in byte_leaves]

    elif request.method == 'GET':
        byte_leaves = app.lithium.leaves
        hex_leaves = [x.encode('hex') for x in byte_leaves]

    dict = {u'leaves': hex_leaves}

    return jsonify(dict)

@app.route('/api/root', methods=['GET'])
def api_root():
    """Return the root for the merkle tree"""
    byte_leaves = app.lithium.leaves
    tree, root = merkle_tree(byte_leaves)
    dict = {u'root': root}
    return jsonify(dict)

@app.route('/api/checkpoints', methods=['GET'])
def api_checkpoint():
    """Returns the checkpoints"""
    return jsonify(app.lithium.checkpoints)

@app.route('/api/blockid', methods=['POST'])
def api_blockid():
    """If passed a valid leaf returns corresponding ionlink blockId"""
    if request.method == 'POST':
        json = request.get_json()
        leaf = json[u'leaf']

    hex_leaves = [x.encode('hex') for x in app.lithium.leaves]
    byte_checkpoints = app.lithium.checkpoints

    if leaf is not None:
        leaf_index = hex_leaves.index(leaf)
        blockid = None
        for block in byte_checkpoints:
            if leaf_index >= byte_checkpoints[block]:
                pass
            else:
                blockid = block
                break

        dict = {u'blockid': str(blockid)}
        return jsonify(dict)
    else:
        return "No valid leaf received."

@app.route('/api/proof', methods=['POST'])
def api_proof():
    """If passed a valid leaf returns merkle proof"""
    if request.method == 'POST':
        json = request.get_json()
        leaf = json[u'leaf']
        blockid = json[u'blockid']
    else:
        return "Please POST leaf data."


    if leaf is not None and blockid is not None:
        nleaves = app.lithium.checkpoints[blockid]
        tree, root = merkle_tree(app.lithium.leaves[:nleaves])

        hex_leaf = leaf.decode('hex')

        path = merkle_path(hex_leaf, tree)

        string_path = [str(x) for x in path]
        dict = {u'proof': string_path}
        return jsonify(dict)
    else:
        return "No valid leaf received."

@app.route('/api/verify', methods=['POST'])
def api_verify_proof():
    """Verifies a given leaf is part of the merkle tree"""
    if request.method == 'POST':
        json = dict(request.get_json())
        leaf = json[u'leaf']
        proof = json[u'proof']
        blockid = json[u'blockid']
    else:
        return "Please POST leaf and path data."

    if leaf is not None and proof is not None and blockid is not None:
        nleaves = app.lithium.checkpoints[blockid]
        leaves = app.lithium.leaves[0:nleaves]
        tree, root = merkle_tree(leaves)

        hex_leaf = leaf.decode('hex')
        proof = merkle_proof(hex_leaf, proof, root)

        return jsonify({"verified":proof})
    else:
        return "No valid leaf or path provided."

class ServerError(Exception):
    status_code = 400

    def __init__(self, message, status_code=None, payload=None):
        Exception.__init__(self)
        self.message = message
        if status_code is not None:
            self.status_code = status_code
        self.payload = payload

    def to_dict(self):
        rv = dict(self.payload or ())
        rv['message'] = self.message
        return rv

@app.errorhandler(ServerError)
def handle_any_error(error):
    response = jsonify(error.to_dict())
    response.status_code = error.status_code
    return response
