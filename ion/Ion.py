# Copyright (c) 2016-2018 Clearmatics Technologies Ltd
# SPDX-License-Identifier: LGPL-3.0+
"""
Ion: command line tool to allow users to intereact with lithium
"""
from __future__ import print_function

from binascii import hexlify, unhexlify

import click
import requests
from sha3 import keccak_256

from .merkle import merkle_hash
from .ethrpc import BadStatusCodeError, BadJsonError, BadResponseError, ConnectionError
from .args import arg_ethrpc, arg_bytes20 #, arg_lithium_api
from .utils import json_dumps

PRIMITIVE = (int, float, str, bool, bytes)

def rpc_call_with_exceptions(function, *args):
    """
    Wraps RPC function calls with expected errors
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
    Mints specified number of token to address
    :param rpc: ip:port of RPC endpoint
    :param account: address of recipient of tokens
    :param tkn: token contract address
    :param value: amount of token to mint
    :return: 0, address balance is printed to console
    """
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
    """
    Returns token balance information for an address
    :param rpc: ip:port of RPC endpoint
    :param account: address to check balance of
    :param tkn: token address of token contract currency of the balance check
    :return: 0, result is printed to console as int
    """
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
    """
    Deposits from an address to IonLock contract with reference
    :param rpc: ip:port of RPC endpoint
    :param account: address of source of funds for deposit
    :param lock: IonLock contract address
    :param tkn: token contract address
    :param value: amount to deposit
    :param ref: str, an arbitrary reference for the payment
    :return: 0, result is printed to the console
    """
    token = rpc.proxy("abi/Token.abi", tkn, account)

    result = rpc_call_with_exceptions(token.metadataTransfer, lock, value, ref)
    if result:
        print("Token transferred.")

        result = rpc_call_with_exceptions(token.balanceOf, account)
        if result is not None:
            print("New balance =", result)

    return 0

@click.command(help="IonLock Withdraw. Withdraws funds from IonLock contract.")
@click.option('--lithium-port', type=int, metavar="8555", required=True, help="Lithium API port")
@click.option('--rpc', callback=arg_ethrpc, metavar="ip:port", default='127.0.0.1:8545', help="Ethereum JSON-RPC server to withdraw funds from")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Beneficiary of funds.")
@click.option('--lock', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLock contract address")
@click.option('--tkn', callback=arg_bytes20, metavar="0x...20", required=True, help="Token contract address of chain to receive from")
@click.option('--value', type=int, required=True, metavar="N", help="Value")
@click.option('--ref', type=str, required=True, metavar="0x...20", help="Payment reference")
def ionlock_withdraw(lithium_port, rpc, account, lock, tkn, value, ref):
    """
    Withdraws to an address from IonLock contract by supplying information from the deposit made on the opposite chain by the account attempting withdrawal
    :param lithium_port: port of lithium api
    :param rpc: ip:port of RPC endpoint to the chain that the withdrawal is being made from
    :param account: address of recipient of withdrawal funds and also address that deposit on opposite chain was made under
    :param lock: IonLock contract address (Currently only works if contract addresses on both chains are identical)
    :param tkn: Token contract address (Currently only works if contract addresses on both chains are identical)
    :param value: Amount to be withdrawn (Currently only works if value deposited on other chain is same as amount to be withdrawn)
    :param ref: str, the payment reference used in the deposit on the opposite chain
    :return: 0, results are printed to the console
    """
    ionlock = rpc.proxy("abi/IonLock.abi", lock, account)
    token = rpc.proxy("abi/Token.abi", tkn, account)

    joined_data = hexlify(account) + hexlify(tkn) + hexlify(lock) + "{0:0{1}x}".format(value,64).encode('utf-8') + hexlify(keccak_256(ref.encode('utf-8')).digest())
    api_url = 'http://127.0.0.1:' + str(lithium_port)
    r = requests.post(api_url + "/api/blockid", json={'leaf': joined_data.decode('ascii')})

    blockid = r.json()['blockid']
    r = requests.post(api_url + "/api/proof", json={'leaf': joined_data.decode('ascii'), 'blockid': blockid})

    path = r.json()['proof']
    path = [int(x) for x in path]
    hashed_ref = hexlify(keccak_256(ref.encode('utf-8')).digest())

    result = rpc_call_with_exceptions(ionlock.Withdraw, value, unhexlify(hashed_ref), int(blockid), path)

    result = rpc_call_with_exceptions(token.balanceOf, account)
    if result is not None:
        print("New balance =", result)

    return 0

@click.command(help="IonLink Verify. Checks the supplied proof with IonLink.")
@click.option('--lithium-port', type=int, metavar="8555", required=True, help="Lithium API port")
@click.option('--rpc', callback=arg_ethrpc, metavar="ip:port", required=True, help="Ethereum JSON-RPC server to verify against")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Beneficiary of funds.")
@click.option('--link', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLock contract address")
@click.option('--lock', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLock contract address of chain to receive from")
@click.option('--tkn', callback=arg_bytes20, metavar="0x...20", required=True, help="Token contract address of chain to receive from")
@click.option('--value', type=int, required=True, metavar="N", help="Value")
@click.option('--ref', type=str, required=True, metavar="abcd", help="Payment reference")
def ionlink_verify(lithium_port, rpc, account, link, lock, tkn, value, ref):
    """
    Verifies a supplied merkle leaf of a deposit with a path against the merkle root held by IonLink
    Deposit chain: chain that deposit was made on
    Verification chain: chain that has been updated with merkle roots that is being verified against
    :param lithium_port: port of lithium api
    :param rpc: ip:port of RPC endpoint to Verification chain
    :param account: address that deposit was made by
    :param link: IonLink contact address on Verification chain
    :param lock: IonLock contract address on Deposit chain
    :param tkn: Token contract address on Deposit chain
    :param value: amount deposited
    :param ref: str, reference used in the deposit
    :return: 0, results are printed to console
    """
    ionlink = rpc.proxy("abi/IonLink.abi", link, account)

    joined_data = hexlify(account) + hexlify(tkn) + hexlify(lock) + "{0:0{1}x}".format(value,64).encode('utf-8') + hexlify(keccak_256(ref.encode('utf-8')).digest())
    hashed_data = merkle_hash(int(joined_data, 16))
    api_url = 'http://127.0.0.1:' + str(lithium_port)
    r = requests.post(api_url + "/api/blockid", json={'leaf': joined_data.decode('ascii')})

    blockid = r.json()['blockid']
    r = requests.post(api_url + "/api/proof", json={'leaf': joined_data.decode('ascii'), 'blockid': blockid})

    path = r.json()['proof']
    path = [int(x) for x in path]

    r = requests.post(api_url + "/api/verify", json={'leaf': joined_data.decode('ascii'), 'proof': path, 'blockid': blockid})
    print("Lithium proof:")
    print(r.text)

    print("IonLink Proof at block id", blockid)
    result = rpc_call_with_exceptions(ionlink.Verify, int(blockid), hashed_data, path)
    print(result)

    return 0



@click.command(help="Merkle proof. Acquires the merkle path to a leaf in Lithium merkle tree for submission during withdraw.")
@click.option('--lithium-port', type=int, metavar="8555", required=True, help="Lithium API port")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Sender of funds.")
@click.option('--lock', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLock contract address of chain to receive from")
@click.option('--tkn', callback=arg_bytes20, metavar="0x...20", required=True, help="Token contract address of chain to receive from")
@click.option('--value', type=int, required=True, metavar="N", help="Value")
@click.option('--ref', type=str, required=True, metavar="abcd", help="Payment reference")
def merkle_proof_path(lithium_port, account, lock, tkn, value, ref):
    """
    Generates a merkle path to a leaf corresponding to a deposit made
    :param lithium_port: port of lithium api
    :param account: address that deposit was made by
    :param lock: IonLock contract address where deposit was made to
    :param tkn: Token contract address of token deposited
    :param value: amount of token deposited
    :param ref: str, reference used for the deposit
    :return: 0, merkle path is printed to the console
    """
    joined_data = hexlify(account) + hexlify(tkn) + hexlify(lock) + "{0:0{1}x}".format(value,64).encode('utf-8') + hexlify(keccak_256(ref.encode('utf-8')).digest())
    print("Joined data", joined_data)

    api_url = 'http://127.0.0.1:' + str(lithium_port)
    r = requests.post(api_url + "/api/blockid", json={'leaf': joined_data.decode('ascii')})

    blockid = r.json()['blockid']
    r = requests.post(api_url + "/api/proof", json={'leaf': joined_data.decode('ascii'), 'blockid':blockid})

    print("Received proof:")
    [print("Path ", r.json()['proof'].index(x), " : ",  x) for x in r.json()['proof']]

    print("Latest IonLink block",blockid)

    return 0

@click.command(help="Merkle Verify. Verifies proof with Lithium merkle tree.")
@click.argument('proof', nargs=-1)
@click.option('--lithium-port', type=int, metavar="8555", required=True, help="Lithium API port")
@click.option('--account', callback=arg_bytes20, metavar="0x...20", required=True, help="Sender of funds.")
@click.option('--lock', callback=arg_bytes20, metavar="0x...20", required=True, help="IonLock contract address of chain to receive from")
@click.option('--tkn', callback=arg_bytes20, metavar="0x...20", required=True, help="Token contract address of chain to receive from")
@click.option('--value', type=int, required=True, metavar="N", help="Value")
@click.option('--ref', type=str, required=True, metavar="abcd", help="Payment reference")
def merkle_verify(proof, lithium_port, account, lock, tkn, value, ref):
    """
    Verifies a supplied merkle path with a leaf corresponding to a deposit made to the merkle tree held by Lithium
    :param proof: space-separated list of hashes as decimal of the path to the leaf
    :param lithium_port: port of lithium api
    :param account: address that deposit was made by
    :param lock: IonLock contract address where deposit was made to
    :param tkn: Token contract address of token deposited
    :param value: amount of token deposited
    :param ref: str, reference used for the deposit
    :return: 0, result is printed to the console
    """
    joined_data = hexlify(account) + hexlify(tkn) + hexlify(lock) + "{0:0{1}x}".format(value,64).encode('utf-8') + hexlify(keccak_256(ref.encode('utf-8')).digest())
    proof = [int(x) for x in proof]
    api_url = 'http://127.0.0.1:' + str(lithium_port)
    r = requests.post(api_url + "/api/blockid", json={'leaf': joined_data.decode('ascii')})

    blockid = r.json()['blockid']
    r = requests.post(api_url + "/api/verify", json={'leaf': joined_data.decode('ascii'), 'proof': proof, 'blockid':blockid})
    print("Received proof:")
    print(r.text)

    return 0



commands = click.Group('commands', help="Ion Interface")
commands.add_command(mint, "mint")
commands.add_command(balance, "balance")
commands.add_command(ionlock_deposit, "deposit")
commands.add_command(ionlock_withdraw, "withdraw")
commands.add_command(ionlink_verify, "ionlink_verify")
commands.add_command(merkle_proof_path, "proof")
commands.add_command(merkle_verify, "verify")
