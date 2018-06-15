import click
from .encoder import RLPEncoder

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