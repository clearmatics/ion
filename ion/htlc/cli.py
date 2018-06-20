# Copyright (c) 2018 Harry Roberts. All Rights Reserved.
# SPDX-License-Identifier: LGPL-3.0+

from __future__ import print_function
import time
from hashlib import sha256

import click

from ..args import arg_ethrpc, arg_bytes20, arg_bytes32, arg_expiry, arg_uint256

from .common import get_random_secret_32, get_default_expiry, make_htlc_proxy


#######################################################################
#
# Command-line interface to the HTLC contract
#
#   $ ion htlc contract [options] sub-command [sub-options]
#
# e.g.
#
#   $ ion htlc contract --account X --contract Y deposit --receiver Z ...
#


@click.command()
@click.pass_obj
@click.option('--receiver', callback=arg_bytes20, metavar="0x...20", required=True, help="Receiver address")
@click.option('--secret', callback=arg_bytes32, metavar="0x...32", default=get_random_secret_32, help="Secret to be supplied upon withdraw")
@click.option('--amount', callback=arg_uint256, metavar='wei', help='Amount of WEI to deposit')
@click.option('--expires', metavar="seconds|unixtime", callback=arg_expiry, type=int, default=get_default_expiry, help="Expiry time, as duration (seconds), or UNIX epoch")
def contract_deposit(contract, receiver, secret, amount, expires):
    now = int(time.time())
    print("Expires in", expires - now, "seconds")

    # TODO: verify balance for account is above or equal to `amount`

    image = sha256(secret).digest()     # the hash pre-image is the 'secret'
    contract.Deposit(receiver, image, expires, value=amount)


@click.command()
@click.pass_obj
@click.option('--secret', callback=arg_bytes32, metavar="0x...32", required=True, help="Exchange ID")
def contract_withdraw(contract, secret):
    image = sha256(secret).digest()     # the hash pre-image is the 'secret'
    contract.Withdraw(image, secret)


@click.command()
@click.pass_obj
@click.option('--image', callback=arg_bytes32, metavar="0x...32", required=True, help="Exchange hash image")
def contract_refund(contract, image):
    contract.Refund(image)


@click.group('contract', help='Command-line interface to Ethereum HTLC contract')
@click.pass_context
@click.option('--rpc', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Ethereum JSON-RPC server")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Account to transfer from.")
@click.option('--contract', callback=arg_bytes20, metavar="0x...20", required=True, help="HTLC contract address")
def contract_multicommand(ctx, rpc, account, contract):
    # Contract will get passed to sub-commands as first object when using `@click.pass_obj`
    ctx.obj = make_htlc_proxy(rpc, contract, account)


contract_multicommand.add_command(contract_deposit, "deposit")
contract_multicommand.add_command(contract_withdraw, "withdraw")
contract_multicommand.add_command(contract_refund, "refund")


#######################################################################
#
# HTLC coordinator server
#

@click.command()
@click.option('--contract', callback=arg_bytes20, metavar="0x...20", required=True, help="HTLC contract address")
@click.option('--rpc', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Ethereum JSON-RPC server")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=False, help="Account to transfer from.")
def coordinator(contract, rpc, account):
    from .coordinator import main as coordinator_main
    return coordinator_main(contract, rpc)


#######################################################################
#
# Multi-command entry-point
#

COMMANDS = click.Group("htlc", help="Hash-Time-Lock Atomic Swap")
COMMANDS.add_command(contract_multicommand, 'contract')
COMMANDS.add_command(coordinator, 'coordinator')


if __name__ == "__main__":
    COMMANDS.main()
