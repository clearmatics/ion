from __future__ import print_function

import click

from ..args import arg_uint256, arg_bytes20


def Token(rpc, contract, account=None):
    with open("abi/Token.abi", "rb") as handle:
        return rpc.proxy(handle, contract, account)


@click.group("erc223", help="ERC-223 compatible token")
@click.pass_context
def main(ctx):
    meta = ctx.meta
    ctx.obj = Token(meta['rpc'], meta['contract'], meta['account'])


@main.command()
@click.argument('value', type=int, callback=arg_uint256)
@click.pass_context
def burn(ctx, value):
    ctx.obj.burn(value)


@main.command()
@click.option('--to', 'dest', type=int, callback=arg_bytes20)
@click.argument('value', type=int, callback=arg_uint256)
@click.pass_context
def transfer(ctx, dest, value):
    ctx.obj.transfer_a9059cbb(dest, value)


@main.command()
@click.argument('value', type=int, callback=arg_uint256)
@click.pass_context
def mint(ctx, value):
    ctx.obj.mint(value)


if __name__ == "__main__":
    main()
