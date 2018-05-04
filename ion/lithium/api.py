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

# from flask import Flask, url_for
from ion.merkle import merkle_tree, merkle_path, merkle_proof


app = Flask(__name__)

@app.route('/')

@app.route('/api/leaves', methods=['GET', 'POST'])
def api_leaves():
    """
        GET:    Return: All leaves held by Lithium
        POST:   Arguments: JSON with a blockid in the form {'blockid': <some_blockid>}
                Return: All relevant leaves under the specified blockid
    """

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
    """
        GET:    Return: The root for the merkle tree of all leaves
    """
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
    """
        POST:   Arguments: JSON with a leaf in the form {'leaf': <some_leaf>}
                Return: If passed a valid leaf returns corresponding blockid otherwise returns error string.
    """
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
    """
        POST:   Arguments: JSON with leaf, blockid in the form {'leaf': <some_leaf>, 'blockid': <some_blockid>}
                Return: If passed valid information returns merkle proof to supplied leaf of relevant blockid
    """
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
    """
        POST:   Arguments: JSON with leaf, proof, blockid in the form {'leaf': <some_leaf>, 'proof': [<array_of_path>] 'blockid': <some_blockid>}
                Return: If passed valid information returns a boolean verifying in the supplied leaf is part of the merkle tree
    """
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
