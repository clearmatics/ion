#!/usr/bin/env python
from __future__ import print_function

import gevent.monkey

gevent.monkey.patch_all()

import click
from flask import Flask, request, jsonify
from jsonrpc2 import JsonRpc

from ..ethrpc import EthJsonRpc
from ..args import arg_ethrpc
from ..utils import marshal, unmarshal, require
from ..plasma.chain import payments_load, blockchain_apply, balances_load, chaindata_latest_get, block_load, chaindata_path, block_genesis, find_block, find_next_block, BlockNotFoundException
from ..plasma.payment import payments_graphviz, SignedPayment
from ..plasma.txpool import TxPool
from .api import PlasmaIonRESTAPI


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
        if block_hash is None:
            print("No plasma blocks exist yet. Creating genesis block.")
            block_hash = block_genesis()
            self.ionlink_sync()
        self._pool = TxPool(block_hash)

    def ionlink_sync(self):
        ion = self._ionlink
        if not ion:
            return False
        ion_latest_block = ion.LatestBlock()

        # plasma_latest_block = block_load(chaindata_latest_get())
        # working_block = block_load(chaindata_get_next_block_from(ion_latest_block))
        # print(working_block)
        # print(plasma_latest_block)

        # latest_path = chaindata_path(ion_latest_block, 'block')

        synced_blocks = 0
        if ion_latest_block != 0:
            try:
                print("Finding last ion block in plasma chain")
                ion_block = find_block(format(ion_latest_block, "02x"))
                print("LAST ION BLOCK ", ion_block.hash.encode('hex'))
                print("Finding next block in plasma chain")
                block_to_sync = find_next_block(ion_block.hash)
                while block_to_sync:
                    print("Syncing block root...... 0x",block_to_sync.root.encode('hex'))
                    ion.Update(block_to_sync.root)
                    synced_blocks += 1
                    block_to_sync = find_next_block(block_to_sync.hash)
            except BlockNotFoundException:
                print("No blocks left to sync")

        else:
            print("Syncing Genesis Block")
            latest_plasma_block = block_load(chaindata_latest_get())
            ion.Update(latest_plasma_block.root)
            synced_blocks += 1

        return synced_blocks > 0

    def ionlink_fetch_tree(self):
        ion = self._ionlink
        if not ion:
            return False
        ion_latest_block = ion.LatestBlock()

        blocks = {'latest': format(ion_latest_block, "02x")}

        while ion_latest_block != 0:
            ion_block_root = ion.GetRoot(ion_latest_block)
            ion_block_prev = ion.GetPrevious(ion_latest_block)

            blocks[format(ion_latest_block, "02x")] = {'root': format(ion_block_root, "02x"), 'prev': format(ion_block_prev, "02x")}

            ion_latest_block = ion_block_prev

        return blocks

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
        # self.ionlink_sync()
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


class BlockMismatchException(Exception):
    pass

# --------------------------------------------------------------------
# Program entry


@click.command()
@click.option('--ion-rpc', default='127.0.0.1:8545', help='Ethereum JSON-RPC HTTP endpoint', callback=arg_ethrpc)
@click.option('--ion-account', help='Ethereum account address')
@click.option('--ion-contract', help='IonLink contract address')
@click.option('--host', help='REST host address, default: any', default='')
@click.option('--port', help='HTTP server port', default=5000, type=int)
def server(ion_rpc, ion_account, ion_contract, host, port):
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

    restAPI = PlasmaIonRESTAPI(server, host, port)
    try:
        restAPI.serve_endpoints()
    except KeyboardInterrupt:
        pass
    return 0


if __name__ == '__main__':
    server(auto_envvar_prefix='ION')
