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
COMMANDS.add_command(coordinator, 'coordinator')


if __name__ == "__main__":
    COMMANDS.main()
