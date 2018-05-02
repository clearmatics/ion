# Copyright (c) 2016-2018 Clearmatics Technologies Ltd
# SPDX-License-Identifier: LGPL-3.0+
"""
Ion: command line tool to allow users to intereact with lithium
"""
from __future__ import print_function

import click

from .args import arg_ethrpc, arg_bytes20
from .ethrpc import BadStatusCodeError, BadJsonError, BadResponseError, ConnectionError

PRIMITIVE = (int, long, float, str, bool)

def rpc_call_with_exceptions(function, *args):
    """
    Test the rpc connections
    """
    try:
        result = function(*args)
        if isinstance(result, PRIMITIVE):
            return result
        return True
    except BadStatusCodeError as err:
        print("Error with status code ", err.message)
    except BadJsonError as err:
        print("BadJson Error: ", err.message)
    except BadResponseError as err:
        print("BadResponseError: ", err.message)
    except ConnectionError as err:
        print("Connection Error: ", err.message)

    return False


@click.command(help="Mint Token. Mints tokens to target account")
@click.option('--rpc', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Source Ethereum JSON-RPC server")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Target account of funds.")
@click.option('--tkn', callback=arg_bytes20, metavar="0x...20", required=True, help="Token contract address")
@click.option('--value', type=int, required=True, metavar="N", help="Value")
def mint(rpc, account, tkn, value):
    """
    Mint: Mints token to the owner address
    """
    token = rpc.proxy("abi/Token.abi", tkn, account)

    result = rpc_call_with_exceptions(token.mint, value)
    if result:
        print("Token minted.")

        result = rpc_call_with_exceptions(token.balanceOf, account)
        if result:
            print("New balance =", result)


    return 0

@click.command(help="IonLock Deposit. Deposits funds to IonLock contract.")
@click.option('--rpc', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Source Ethereum JSON-RPC server")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Source account of funds.")
@click.option('--lock', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLock contract address")
@click.option('--tkn', callback=arg_bytes20, metavar="0x...20", required=True, help="Token contract address")
@click.option('--value', type=int, required=True, metavar="N", help="Value")
@click.option('--data', type=str, required=True, metavar="abcd", help="Some payment reference")
def ionlock_deposit(rpc, account, lock, tkn, value, data):
    """
    deposits token from account to the IonLock contract
    """
    token = rpc.proxy("abi/Token.abi", tkn, account)
    print(data)

    result = rpc_call_with_exceptions(token.metadataTransfer, lock, value, data)
    if result:
        print("Token transferred.")

        result = rpc_call_with_exceptions(token.balanceOf, account)
        if result:
            print("New balance =", result)

    return 0



commands = click.Group('commands', help="Ion Interface")
commands.add_command(mint, "mint")
commands.add_command(ionlock_deposit, "deposit")
