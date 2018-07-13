import rlp
from .ethrpc import EthJsonRpc
from ethereum.utils import sha3
from .utils import require
import click

class Block(object):
    def __init__(self, block_json):
        self.parent_hash = hexstring_to_bytes(block_json['parentHash'])
        self.sha3_uncles = hexstring_to_bytes(block_json['sha3Uncles'])
        self.coinbase = hexstring_to_bytes(block_json['miner'])
        self.root = hexstring_to_bytes(block_json['stateRoot'])
        self.tx_hash = hexstring_to_bytes(block_json['transactionsRoot'])
        self.receipt_hash = hexstring_to_bytes(block_json['receiptsRoot'])
        self.bloom = hexstring_to_bytes(block_json['logsBloom'])
        self.difficulty = hexstring_to_bytes(block_json['difficulty'])
        self.number = hexstring_to_bytes(block_json['number'])
        self.gas_limit = hexstring_to_bytes(block_json['gasLimit'])
        self.gas_used = hexstring_to_bytes(block_json['gasUsed'])
        self.timestamp = hexstring_to_bytes(block_json['timestamp'])
        self.extra_data = hexstring_to_bytes(block_json['extraData'])
        self.mix = hexstring_to_bytes(block_json['mixHash'])
        self.nonce = hexstring_to_bytes(block_json['nonce'])
        self.transactions = block_json['transactions']

    @property
    def header(self):
        return [
            self.parent_hash,
            self.sha3_uncles,
            self.coinbase,
            self.root,
            self.tx_hash,
            self.receipt_hash,
            self.bloom,
            self.difficulty,
            self.number,
            self.gas_limit,
            self.gas_used,
            self.timestamp,
            self.extra_data,
            self.mix,
            self.nonce
        ]


# Converts hex string into a bytes of the hex representation from "0x..." to b'\x...'
# If value is less than a byte, returns as int
def hexstring_to_bytes(hex):
    try:
        return bytes.fromhex(hex[2:])
    except ValueError:
        return int(hex, 16)

class RLPEncoder(object):
    def __init__(self, host, port):
        if port == 443:
            tls = True
        else:
            tls = False
        self.rpc = EthJsonRpc(host, port, tls)

    def get_block_by_number(self, n):
        block = self.rpc.eth_getBlockByNumber(n)
        return block

    def get_block_by_hash(self, h):
        block = self.rpc.eth_getBlockByHash(h)
        return block

    def get_block(self, x):
        if isinstance(x, int):
            return self.get_block_by_number(x)
        elif isinstance(x, str):
            return self.get_block_by_hash(x)
        else:
            raise Exception("Must supply number or hash")

    def get_transactions(self, block):
        block = Block(block)
        return block.transactions

    def encode_block(self, block):
        return rlp.encode(Block(block).header)

    def hash_block_header(self, block):
        encoded = self.encode_block(block)
        hash = sha3(encoded).hex()
        require('0x'+hash == block['hash'], "Block hash and hashed header do not match:\n{} expected\n{} acquired".format(block['hash'], '0x'+hash))

        return hash



@click.command(help="Returns an RLP encoded block in hexadecimal format.")
@click.argument('rpc-host', nargs=1, type=str)
@click.argument('rpc-port', nargs=1, type=int)
@click.argument('number', nargs=1, type=int)
def get_encoded_block(rpc_host, rpc_port, number):
    rlp_encoder = RLPEncoder(rpc_host, rpc_port)
    block = rlp_encoder.get_block(number)
    click.echo('0x'+rlp_encoder.encode_block(block).hex())

@click.command(help="Returns block hash in hexadecimal format.")
@click.argument('rpc-host', nargs=1, type=str)
@click.argument('rpc-port', nargs=1, type=int)
@click.argument('number', nargs=1, type=int)
def get_block_hash(rpc_host, rpc_port, number):
    rlp_encoder = RLPEncoder(rpc_host, rpc_port)
    block = rlp_encoder.get_block(number)
    click.echo('0x'+rlp_encoder.hash_block_header(block))


@click.command(help="Returns a list of transaction hashes from a specified block")
@click.argument('rpc-host', nargs=1, type=str)
@click.argument('rpc-port', nargs=1, type=int)
@click.argument('number', nargs=1, type=int)
def get_block_transactions(rpc_host, rpc_port, number):
    rlp_encoder = RLPEncoder(rpc_host, rpc_port)
    block = rlp_encoder.get_block(number)
    click.echo(rlp_encoder.get_transactions(block))

commands = click.Group('commands')
commands.add_command(get_encoded_block, "encodeblock")
commands.add_command(get_block_hash, "blockhash")
commands.add_command(get_block_transactions, "gettx")

if __name__ == "__main__":
    commands.main()