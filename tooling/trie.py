from patricia import trie
import rlp
from .ethrpc import EthJsonRpc
from ethereum.utils import sha3
from ethereum.transactions import Transaction

import click


class Trie(object):
    def __init__(self, host, port, tls=False):
        self.rpc = EthJsonRpc(host, port, tls)

    def construct_tx_trie_from_block(self, n):
        try:
            n = int(n)
            block = self.rpc.eth_getBlockByNumber(n)
        except ValueError:
            block = self.rpc.eth_getBlockByNumber()

        transactions = block['transactions']

        for tx in transactions:
            print(tx)
            print(rlp.encode(tx['transactionIndex']))
            tx_data = [tx['nonce'], tx['gasPrice'], tx['gas'], tx['to'], tx['value'], tx['input']]
            print(rlp.encode(tx_data))

    def get_tx_root(self, n):
        try:
            n = int(n)
            block = self.rpc.eth_getBlockByNumber(n)
        except ValueError:
            block = self.rpc.eth_getBlockByNumber()

        return block['transactionsRoot']

def instantiate_trie(rpc_host, rpc_port):
    if rpc_port == 443:
        t = Trie(rpc_host, rpc_port, True)
    else:
        t = Trie(rpc_host, rpc_port)

    return t

@click.command(help="Returns a list of transaction hashes from a specified block")
@click.argument('rpc-host', nargs=1, type=str)
@click.argument('rpc-port', nargs=1, type=int)
@click.option('--number', nargs=1, default="LATEST", type=str)
def construct_tx_trie(rpc_host, rpc_port, number):
    t = instantiate_trie(rpc_host, rpc_port)

    t.construct_tx_trie_from_block(number)

@click.command(help="Returns TxTrie root of a specified block")
@click.argument('rpc-host', nargs=1, type=str)
@click.argument('rpc-port', nargs=1, type=int)
@click.option('--number', nargs=1, default="LATEST", type=str)
def get_tx_root(rpc_host, rpc_port, number):
    t = instantiate_trie(rpc_host, rpc_port)

    root = t.get_tx_root(number)
    click.echo(root)

commands = click.Group('commands')
commands.add_command(construct_tx_trie, "tx")
commands.add_command(get_tx_root, "txroot")

if __name__ == "__main__":
    commands.main()
