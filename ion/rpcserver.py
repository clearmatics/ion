#!/usr/bin/env python
from __future__ import print_function
import gevent
import gevent.monkey
gevent.monkey.patch_all()

import sys
import argparse
from gevent.pywsgi import WSGIServer
from flask import Flask, request, jsonify 
from jsonrpc2 import JsonRpc
from ethjsonrpc import EthJsonRpc

from .args import EthRpcAction, Bytes20Action
from .utils import marshal, unmarshal, require
from .chain import payments_load, blockchain_apply, balances_load, chaindata_latest_get, block_load, chaindata_path
from .payment import payments_graphviz, SignedPayment
from .txpool import TxPool
from .solproxy import solproxy


class IonRpcServer(object):
    def __init__(self, opts=None):
        self._opts = opts or argparse.Namespace()
        self._new_pool()
        self._webapp_init()

    def _webapp_init(self):
        self._webapp = Flask(__name__)
        self._webapp.add_url_rule('/rpc', 'jsonrpc_handler', self._jsonrpc_handler, methods=['POST'])
        self._webrpc = JsonRpc(methods={
            name: getattr(self, name)
            for name in dir(self)
            if name[0] != '_'
        })

    def _new_pool(self, block_hash=None):
        block_hash = block_hash or chaindata_latest_get()
        require( block_hash is not None, "Must specify block hash" )
        self._pool = TxPool(block_hash)

    def ionlink_sync(self):
        ion = getattr(self._opts, 'ion', None)
        if not ion:
            return False

        ion_latest_block = ion.LatestBlock()
        my_block = block_load(chaindata_latest_get())

        require( ion_latest_block == my_block.hash, "Block mismatch" )

        return ion.Update([my_block.root])

    def _jsonrpc_handler(self):
        print("Request is", request.json)
        response = self._webrpc(request.json)
        print("Response is", response)
        return jsonify(response)

    def graph(self):
        g = payments_graphviz(self._pool.payments)
        g.render(chaindata_path('txpool', 'graphviz'))
        return True

    def block_get(self, block_no):
        """Retrieve information about a specific block
        :rtype: .model.Block
        """
        return block_load(block_no).marshal()

    def block_get_latest(self):
        """Lookup the latest block number, perform `block_get` on it"""
        latest = chaindata_latest_get()
        require( latest is not None, "No latest block" )
        return self.block_get(latest)

    def block_hash(self):
        return marshal(self._pool.target)

    def balance(self, currency, holder, block_hash=None):
        currency = unmarshal(currency)
        holder = unmarshal(holder)
        if block_hash is None:
            block_hash = chaindata_latest_get()
        # XXX: this is slow, loads whole balances dict!
        balances = balances_load(block_hash)
        require( currency in balances, "Unknown currency" )
        return balances[currency].get(holder, 0)

    def block_commit(self):
        """Whatever is in the payment list, commit it to a block"""
        # then call `blockchain_apply` on the payments
        prev_hash = chaindata_latest_get()
        balances = balances_load(prev_hash)
        pruned_payments = self._pool.prune(balances)
        block = blockchain_apply(prev_hash, pruned_payments)
        self.ionlink_sync()
        self._new_pool(block.hash)
        return block.marshal()

    def payment_submit(self, sp):
        """Submit a payment, for inclusion in the next block"""
        return self._pool.add(SignedPayment.unmarshal(sp))

    def payment_pending(self, ref):
        return self._pool.pending(unmarshal(ref))

    def payment_cancel(self, ref):
        # TODO: require signature from owner to cancel
        return self._pool.cancel(unmarshal(ref))

    def payment_confirmed(self, block_no, ref):
        """Was a given payment included in the block"""
        ref = unmarshal(ref)
        signed_payments = payments_load(block_no)
        return any([sp.p.r == ref for sp in signed_payments])


# --------------------------------------------------------------------
# Program entry


def rpcserver_options(args=None):
    parser = argparse.ArgumentParser(description="Ion: JSON-RPC Server")

    # Connect to uplink
    # TODO: have ionlink options, --ion-rpc, --ion-account, --ion-contract etc.
    parser.add_argument('--ion-rpc', dest='ion_rpc', action=EthRpcAction, default='127.0.0.1:8545',
                        help='Ethereum JSON-RPC HTTP endpoint')
    parser.add_argument('--ion-account', dest='ion_account', action=Bytes20Action,
                        help='Ethereum account address, 0x...20')
    parser.add_argument('--ion-contract', dest='ion_contract', action=Bytes20Action,
                        help='IonLink contract address, 0x...20')

    parser.add_argument('-p', '--port', dest='port', type=int,
                        help='HTTP server port', default=5000)
    parser.add_argument('-l', '--listen', dest='listen', type=str, default='',
                        help="Listen address, default: any")

    opts = parser.parse_args(args or sys.argv[1:])

    if not opts.ion_contract or not opts.ion_account:
        print("IonLink disabled")
        opts.ionlink = None
    else:
        if not opts.rpc:
            opts.ion_rpc = EthJsonRpc('127.0.0.1', 8545)
        # TODO: load ABI from package resources
        opts.ion = solproxy(opts.ion_rpc, "abi/IonLink.abi", opts.ion_contract, opts.ion_account)

    return opts


def main():
    opts = rpcserver_options()
    server = IonRpcServer(opts)
    gevent_server = WSGIServer((opts.listen, opts.port), server._webapp)
    try:
        gevent_server.serve_forever()
    except KeyboardInterrupt:
        pass
    return 0


if __name__ == '__main__':
    sys.exit(main())
