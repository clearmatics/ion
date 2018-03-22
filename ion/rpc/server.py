#!/usr/bin/env python
from __future__ import print_function

import gevent.monkey

gevent.monkey.patch_all()

import click
from gevent.pywsgi import WSGIServer
from flask import Flask, request, jsonify 
from jsonrpc2 import JsonRpc

from ..ethrpc import EthJsonRpc
from ..args import arg_ethrpc
from ..utils import marshal, unmarshal, require
from ..plasma.chain import payments_load, blockchain_apply, balances_load, chaindata_latest_get, block_load, chaindata_path
from ..plasma.payment import payments_graphviz, SignedPayment
from ..plasma.txpool import TxPool


class IonRpcServer(object):
    def __init__(self, ionlink=None):
        self._ionlink = ionlink
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
        ion = self._ionlink
        if not ion:
            return False
        ion_latest_block = ion.GetLatestBlockHash()
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


@click.command()
@click.option('--ion-rpc', default='127.0.0.1:8545', help='Ethereum JSON-RPC HTTP endpoint', callback=arg_ethrpc)
@click.option('--ion-account', help='Ethereum account address')
@click.option('--ion-contract', help='IonLink contract address')
@click.option('--listen', help='Listen address, default: any', default='')
@click.option('--port', help='HTTP server port', default=5000, type=int)
def server(ion_rpc, ion_account, ion_contract, listen, port):
    """
    RPC server for Ion.

    :type ion_rpc: EthJsonRpc
    """
    print(ion_rpc.net_version())
    if not ion_contract or not ion_account:
        print("IonLink disabled")
        ionlink = None
    else:
        if not ion_rpc:
            ion_rpc = EthJsonRpc('127.0.0.1', 8545)
        # TODO: load ABI from package resources
        ionlink = ion_rpc.proxy("abi/IonLink.abi", ion_contract, ion_account)

    server = IonRpcServer(ionlink)
    gevent_server = WSGIServer((listen, port), server._webapp)
    try:
        gevent_server.serve_forever()
    except KeyboardInterrupt:
        pass
    return 0


if __name__ == '__main__':
    server(auto_envvar_prefix='ION')
