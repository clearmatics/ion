from __future__ import print_function

import click
import requests
import simplejson

from ion.merkle import merkle_hash
from ethereum.utils import keccak
from .ethrpc import EthJsonRpc, BadStatusCodeError, BadJsonError, BadResponseError, ConnectionError
from .args import arg_ethrpc, arg_bytes20, arg_bytes, make_uint_n, make_bytes_n
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
        if result is not None:
            print("New balance =", result)

    return 0

@click.command(help="Balance check. Check token balance of target account")
@click.option('--rpc', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Source Ethereum JSON-RPC server")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Account of balance check.")
@click.option('--tkn', callback=arg_bytes20, metavar="0x...20", required=True, help="Token contract address")
def balance(rpc, account, tkn):
    token = rpc.proxy("abi/Token.abi", tkn, account)

    result = rpc_call_with_exceptions(token.balanceOf, account)
    if result is not None:
        print("Balance =", result)


    return 0

@click.command(help="IonLock Deposit. Deposits funds to IonLock contract.")
@click.option('--rpc', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Source Ethereum JSON-RPC server")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Source account of funds.")
@click.option('--lock', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLock contract address")
@click.option('--tkn', callback=arg_bytes20, metavar="0x...20", required=True, help="Token contract address")
@click.option('--value', type=int, required=True, metavar="N", help="Value")
@click.option('--ref', type=str, required=True, metavar="abcd", help="Some payment reference")
def ionlock_deposit(rpc, account, lock, tkn, value, ref):
    token = rpc.proxy("abi/Token.abi", tkn, account)

    result = rpc_call_with_exceptions(token.metadataTransfer, lock, value, ref)
    if result:
        print("Token transferred.")

        result = rpc_call_with_exceptions(token.balanceOf, account)
        if result is not None:
            print("New balance =", result)

    return 0

@click.command(help="IonLock Withdraw. Withdraws funds from IonLock contract.")
@click.option('--rpc-from', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Ethereum JSON-RPC server to obtain proof from")
@click.option('--rpc-to', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Ethereum JSON-RPC server to verify against")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Beneficiary of funds.")
@click.option('--lock', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLock contract address")
@click.option('--tkn', callback=arg_bytes20, metavar="0x...20", required=True, help="Token contract address of chain to receive from")
@click.option('--value', type=int, required=True, metavar="N", help="Value")
@click.option('--ref', type=str, required=True, metavar="0x...20", help="Payment reference")
def ionlock_withdraw(rpc_from, rpc_to, account, lock, tkn, value, ref):
    ionlock = rpc_to.proxy("abi/IonLock.abi", lock, account)
    token = rpc_to.proxy("abi/Token.abi", tkn, account)

    joined_data = account.encode('hex') + tkn.encode('hex') + lock.encode('hex') + "{0:0{1}x}".format(value,64) + keccak.new(digest_bits=256).update(str(ref)).hexdigest()
    api_url = 'http://127.0.0.1:' + str(rpc_from.port + 10)
    r = requests.post(api_url + "/api/blockid", json={'leaf': joined_data})

    try:
        blockid = r.json()['blockid']
        r = requests.post(api_url + "/api/proof", json={'leaf': joined_data, 'blockid': blockid})

        path = r.json()['proof']
        path = [int(x) for x in path]
        hashed_ref = keccak.new(digest_bits=256).update(str(ref)).hexdigest()

        result = rpc_call_with_exceptions(ionlock.Withdraw, value, hashed_ref.decode('hex'), int(blockid), path)

        result = rpc_call_with_exceptions(token.balanceOf, account)
        if result is not None:
            print("New balance =", result)


    except simplejson.errors.JSONDecodeError as e:
        print(e.message)

    return 0

@click.command(help="IonLink Verify. Checks the supplied proof with IonLink.")
@click.option('--rpc-from', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Ethereum JSON-RPC server to obtain proof from")
@click.option('--rpc-to', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Ethereum JSON-RPC server to verify against")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Beneficiary of funds.")
@click.option('--link', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLock contract address")
@click.option('--lock', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLock contract address of chain to receive from")
@click.option('--tkn', callback=arg_bytes20, metavar="0x...20", required=True, help="Token contract address of chain to receive from")
@click.option('--value', type=int, required=True, metavar="N", help="Value")
@click.option('--ref', type=str, required=True, metavar="abcd", help="Payment reference")
def ionlink_verify(rpc_from, rpc_to, account, link, lock, tkn, value, ref):
    ionlink = rpc_to.proxy("abi/IonLink.abi", link, account)

    joined_data = account.encode('hex') + tkn.encode('hex') + lock.encode('hex') + "{0:0{1}x}".format(value,64) + keccak.new(digest_bits=256).update(str(ref)).hexdigest()
    hashed_data = merkle_hash(int(joined_data, 16))
    api_url = 'http://127.0.0.1:' + str(rpc_from.port + 10)
    r = requests.post(api_url + "/api/blockid", json={'leaf': joined_data})

    try:
        blockid = r.json()['blockid']
        r = requests.post(api_url + "/api/proof", json={'leaf': joined_data, 'blockid': blockid})

        path = r.json()['proof']
        path = [int(x) for x in path]

        r = requests.post(api_url + "/api/verify", json={'leaf': joined_data, 'proof': path, 'blockid': blockid})
        print("Lithium proof:")
        print(r.text)

        print("IonLink Proof at block id", blockid)
        result = rpc_call_with_exceptions(ionlink.Verify, int(blockid), hashed_data, path)
        print(result)

    except simplejson.errors.JSONDecodeError as e:
        print(e.message)

    return 0



@click.command(help="Merkle proof. Acquires the merkle path to a leaf in Lithium merkle tree for submission during withdraw.")
@click.option('--rpc', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Source Ethereum JSON-RPC server")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Sender of funds.")
@click.option('--lock', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLock contract address of chain to receive from")
@click.option('--tkn', callback=arg_bytes20, metavar="0x...20", required=True, help="Token contract address of chain to receive from")
@click.option('--value', type=int, required=True, metavar="N", help="Value")
@click.option('--ref', type=str, required=True, metavar="abcd", help="Payment reference")
def merkle_proof_path(rpc, account, lock, tkn, value, ref):

    joined_data = account.encode('hex') + tkn.encode('hex') + lock.encode('hex') + "{0:0{1}x}".format(value,64) + keccak.new(digest_bits=256).update(str(ref)).hexdigest()
    api_url = 'http://127.0.0.1:' + str(rpc.port + 10)
    r = requests.post(api_url + "/api/blockid", json={'leaf': joined_data})

    try:
        blockid = r.json()['blockid']
        r = requests.post(api_url + "/api/proof", json={'leaf': joined_data, 'blockid':blockid})

        print("Received proof:")
        [print("Path ", r.json()['proof'].index(x), " : ",  x) for x in r.json()['proof']]

        print("Latest IonLink block",blockid)

    except simplejson.errors.JSONDecodeError as e:
        print(e.message)

    # Requests errors have not been formally caught as all errors will pertain to the success of making
    # connections and requests to the API which will raise the relevant errors

    return 0

@click.command(help="Merkle Verify. Verifies proof with Lithium merkle tree.")
@click.argument('proof', nargs=-1)
@click.option('--rpc', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Source Ethereum JSON-RPC server")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Sender of funds.")
@click.option('--lock', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLock contract address of chain to receive from")
@click.option('--tkn', callback=arg_bytes20, metavar="0x...20", required=True, help="Token contract address of chain to receive from")
@click.option('--value', type=int, required=True, metavar="N", help="Value")
@click.option('--ref', type=str, required=True, metavar="abcd", help="Payment reference")
def merkle_verify(proof, rpc, account, lock, tkn, value, ref):

    joined_data = account.encode('hex') + tkn.encode('hex') + lock.encode('hex') + "{0:0{1}x}".format(value,64) + keccak.new(digest_bits=256).update(str(ref)).hexdigest()
    proof = [int(x) for x in proof]
    api_url = 'http://127.0.0.1:' + str(rpc.port + 10)
    r = requests.post(api_url + "/api/blockid", json={'leaf': joined_data})

    try:
        blockid = r.json()['blockid']
        r = requests.post(api_url + "/api/verify", json={'leaf': joined_data, 'proof': proof, 'blockid':blockid})
        print("Received proof:")
        print(r.text)

    except simplejson.errors.JSONDecodeError as e:
        print(e.message)

    # Requests errors have not been formally caught as all errors will pertain to the success of making
    # connections and requests to the API which will raise the relevant errors

    return 0



commands = click.Group('commands', help="Ion Interface")
commands.add_command(mint, "mint")
commands.add_command(balance, "balance")
commands.add_command(ionlock_deposit, "deposit")
commands.add_command(ionlock_withdraw, "withdraw")
commands.add_command(ionlink_verify, "ionlink_verify")
commands.add_command(merkle_proof_path, "proof")
commands.add_command(merkle_verify, "verify")