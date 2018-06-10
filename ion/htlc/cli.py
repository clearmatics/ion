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


@click.command(help="Deposit into Hash-Time-Lock contract")
@click.option('--rpc', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Ethereum JSON-RPC server")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Account to transfer from.")
@click.option('--contract', callback=arg_bytes20, metavar="0x...20", required=True, help="HTLC contract address")
@click.option('--receiver', callback=arg_bytes20, metavar="0x...20", required=True, help="Receiver address")
@click.option('--secret', callback=arg_bytes32, metavar="0x...32", default=get_random_secret_32, help="Secret to be supplied upon withdraw")
@click.option('--expires', metavar="seconds|unixtime", callback=arg_expiry, type=int, default=get_default_expiry, help="Expiry time, as duration (seconds), or UNIX epoch")
def deposit(rpc, account, contract, receiver, secret, expires):
	now = int(time.time())
	print("Expires in", expires - now, "seconds")
	api = make_htlc_proxy(rpc, contract, account)
	image = sha256(secret).digest()		# the hash pre-image is the 'secret'
	api.Deposit( receiver, image, expires )


@click.command(help="Withdraw from Hash-Time-Lock contract")
@click.option('--rpc', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Ethereum JSON-RPC server")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Account to withdraw to")
@click.option('--contract', callback=arg_bytes20, metavar="0x...20", required=True, help="HTLC contract address")
@click.option('--secret', callback=arg_bytes32, metavar="0x...32", required=True, help="Exchange ID")
def withdraw(rpc, account, contract, secret):
	api = make_htlc_proxy(rpc, contract, account)
	image = sha256(secret).digest()		# the hash pre-image is the 'secret'
	api.Withdraw( image, secret )


@click.command(help="Refund a Hash-Time-Lock contract")
@click.option('--rpc', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Ethereum JSON-RPC server")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Account to withdraw to")
@click.option('--contract', callback=arg_bytes20, metavar="0x...20", required=True, help="HTLC contract address")
@click.option('--image', callback=arg_bytes32, metavar="0x...32", required=True, help="Exchange hash image")
def refund(rpc, account, contract, image):
	api = make_htlc_proxy(rpc, contract, account)
	api.Refund( image )


commands = click.Group("htlc", help="Hash-Time-Lock Contract Interface")
commands.add_command(deposit, "deposit")
commands.add_command(withdraw, "withdraw")
commands.add_command(refund, "refund")


if __name__ == "__main__":
	commands.main()
