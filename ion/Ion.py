from __future__ import print_function

import click

from .args import arg_ethrpc, arg_bytes20
from .ethrpc import BadStatusCodeError, BadJsonError, BadResponseError, ConnectionError

primitive = (int, long, float, str, bool)

def rpc_call_with_exceptions(function, *args):

    try:
        result = function(*args)
        if isinstance( result, primitive ):
            return result
        return True
    except BadStatusCodeError as e:
        print("Error with status code ", e.message)
    except BadJsonError as e:
        print("BadJson Error: ", e.message)
    except BadResponseError as e:
        print("BadResponseError: ", e.message)
    except ConnectionError as e:
        print("Connection Error: ", e.message)

    return False


@click.command(help="Mint Token. Mints tokens to target account")
@click.option('--rpc', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Source Ethereum JSON-RPC server")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Target account of funds.")
@click.option('--tkn', callback=arg_bytes20, metavar="0x...20", required=True, help="Token contract address")
@click.option('--value', type=int, required=True, metavar="N", help="Value")
def mint(rpc, account, tkn, value):
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