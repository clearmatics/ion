import click

from ..args import arg_ethrpc, arg_bytes20

from .erc223 import main as erc223_main

@click.group("onchain", help="On-chain interfaces")
@click.option('--rpc', '-r', metavar="ip:port", callback=arg_ethrpc, help="Etherum JSON-RPC server", default='127.0.0.1:8545')
@click.option('--contract', '-c', metavar="0x...20", callback=arg_bytes20, help="Contract address", required=True)
@click.option('--account', '-a', metavar="0x...20", callback=arg_bytes20, help="Account to use for transactions", required=True)
@click.pass_context
def commands(ctx, rpc, contract, account):
    """
    :type ctx: click.Context
    """
    ctx.meta['rpc'] = rpc
    ctx.meta['contract'] = contract
    ctx.meta['account'] = account


commands.add_command(erc223_main, "erc223")


if __name__ == "__main__":
    commands()
