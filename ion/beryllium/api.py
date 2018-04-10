#!/usr/bin/env python

import os
from flask import Flask, request, jsonify

from ion.rpc.IonClient import IonClient


class PlasmaIonRESTAPI(object):
    def __init__(self, rpc_endpoint, host=None, port=None):
        self.rpc_endpoint = rpc_endpoint
        self.host = host
        self.port = port

    def serve_endpoints(self):
        if self.host is None:
            self.host = '127.0.0.1'
        if self.port is None:
            self.port = 5000

        app = Flask("REST API")

        @app.route('/api/add_payment', methods=['POST'])
        def add_payment():
            print(request)
            if request.method == 'POST':
                sender_public = request.args.get('sender_public')
                recipient_public = request.args.get('recipient_public')
                value = int(request.args.get('value'))
                currency = request.args.get('currency')
                data = {'sender': sender_public, 'recipient': recipient_public, 'value': value, 'currency': currency}

                sender = IonClient(self.rpc_endpoint, os.urandom(32))
                recipient = IonClient(self.rpc_endpoint, os.urandom(32))
                sender.mint(1000)
                sender.pay(recipient, sender.public, value)

                return jsonify(data)


        @app.route('/api/commit_block')
        def commit_block():
            self.rpc_endpoint.block_commit()
            return jsonify("Complete")

        @app.route('/api/sync_blocks')
        def sync_blocks():
            if self.rpc_endpoint.ionlink_sync():
                return "Successful Sync"
            else:
                return "Unsuccessful Sync"


        app.run(host=self.host, port=self.port)

