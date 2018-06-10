# Copyright (c) 2018 Harry Roberts. All Rights Reserved.
# SPDX-License-Identifier: LGPL-3.0+
from __future__ import print_function
import time
import os
from hashlib import sha256

import click

from ..args import arg_ethrpc, arg_bytes20, arg_bytes32

ONE_MINUTE = 60
ONE_HOUR = ONE_MINUTE * 60
ONE_DAY = ONE_HOUR * 24
ONE_YEAR = ONE_DAY * 365

DEFAULT_EXPIRY_DURATION = 10 * ONE_MINUTE
DURATION_OR_EPOCH_SPLIT = ONE_YEAR


def make_htlc_proxy(rpc, contract, account):
    # TODO: embed 'abi/HTLC.abi' file in package resources?
    return rpc.proxy('abi/HTLC.abi', contract, account)


def get_default_expiry():
    return int(time.time()) + DEFAULT_EXPIRY_DURATION


def arg_expiry(ctx, param, value):
    """
    Accepts either a duration, or an absolute UNIX epoch time
    Returns absolute UNIX epoch time
    """
    value = int(value)
    if value < DURATION_OR_EPOCH_SPLIT:
        return int(time.time()) + value
    return value


def get_random_secret_32():
    return '0x' + os.urandom(32).encode('hex')


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
@click.option('--expires', metavar="seconds|unixtime", callback=arg_expiry, type=int, default=get_default_expiry, help="Expiry time, as duration (seconds), or UNIX epoch")
def contract_deposit(contract, receiver, secret, expires):
    now = int(time.time())
    print("Expires in", expires - now, "seconds")

    image = sha256(secret).digest()     # the hash pre-image is the 'secret'
    contract.Deposit(receiver, image, expires)


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
# Multi-command entry-point
#

COMMANDS = click.Group("htlc", help="Hash-Time-Lock Atomic Swap")
COMMANDS.add_command(contract_multicommand, 'contract')


if __name__ == "__main__":
    COMMANDS.main()
