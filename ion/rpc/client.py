import click
import os
import sys

from .IonClient import IonClient
from .server import IonRpcServer
from ..args import arg_ethrpc
from ..ethrpc import EthJsonRpc


# ------------------------------------------------------------------
# Functional Implementation

def create_client(rpc_endpoint, secret):
    client = IonClient(rpc_endpoint, secret)
    return client, client.public

def make_payment(rpc_endpoint, secret):
    print("Making payment")
    client, public = create_client(rpc_endpoint, secret)

    other_client, other_public = create_client(rpc_endpoint, os.urandom(32))

    client.mint(1000)

    client.pay(other_client, public, 500)

    client.commit()

    if rpc_endpoint.ionlink_sync(): print("Payment success")
    else: print("Sync Failed")


def print_ion_history_tree(rpc_endpoint):
    blocks = rpc_endpoint.ionlink_fetch_tree()

    latest = blocks['latest']

    while long(latest, 16) != 0:
        block = blocks[latest]

        print("==================== Block Hash ====================")
        print(latest)

        print("\n==================== Block Root ====================")
        print(block['root'])

        print("\n================== Block Previous ==================")
        print(block['prev'])

        latest = block['prev']

        if long(latest, 16) != 0:
            print("                                                              ||")
            print("                                                              ||")
            print("                                                              ||")
            print("                                                              ||")
            print("                                                              \\/")


# ------------------------------------------------------------------
# Tests

def test_newclient(rpc_endpoint, secret=None):
    client = IonClient(rpc_endpoint, secret or os.urandom(32))
    return client, client.public


def test_pay(rpc_endpoint): 
    client_A, currency_A = test_newclient(rpc_endpoint)
    client_B, currency_B = test_newclient(rpc_endpoint)

    client_A.mint(1000)
    client_B.mint(1000)

    # A Pays B 1000A
    client_A.pay(client_B.public, currency_A, 1000)

    # B pays A 500B
    client_B.pay(client_A.public, currency_B, 500)

    client_A.commit()

    assert client_A.balance(currency_A) == 0
    assert client_A.balance(currency_B) == 500

    assert client_B.balance(currency_A) == 1000
    assert client_B.balance(currency_B) == 500

    print "Test pay passed"


def test_swap(rpc_endpoint):
    client_A, currency_A = test_newclient(rpc_endpoint)
    client_B, currency_B = test_newclient(rpc_endpoint)
    client_C, currency_C = test_newclient(rpc_endpoint)
    client_E, currency_E = test_newclient(rpc_endpoint)

    # Client A gives B, C and D 500 units of test currency
    client_A.mint(1000)
    client_A.pay(client_B, currency_A, 500)
    client_A.pay(client_C, currency_A, 500)

    # E then mints 1000 units of their currency
    client_E.mint(1000)
    client_A.commit()

    assert client_A.balance(currency_A) == 0
    assert client_B.balance(currency_A) == 500
    assert client_C.balance(currency_A) == 500
    assert client_E.balance(currency_E) == 1000

    # Reference is secret key for a per-trade account
    e_for_a_secret = os.urandom(32)
    trade_E_for_A, currency_E_for_A = test_newclient(rpc_endpoint, e_for_a_secret)

    # E offers to swap 1000E for 500A
    # It sends it to the special 'E for A' account which has a known secret
    # Anybody can use the E_for_A account to complete transactions
    e_needs = os.urandom(32)
    client_E.pay(trade_E_for_A, currency_E, 1000, e_for_a_secret, deps=[
        (client_E, currency_A, 500, e_needs)
    ])
    trade_E_for_A.pay(client_E, currency_A, 500, e_needs, deps=[
        (trade_E_for_A, currency_E, 1000, e_for_a_secret),
    ])

    # B offers to swap 500A for 1000E
    b_wins = os.urandom(32)
    b_needs = os.urandom(32)
    client_B.pay(trade_E_for_A, currency_A, 500, b_needs, deps=[
        (client_B, currency_E, 1000, b_wins),
    ])
    trade_E_for_A.pay(client_B, currency_E, 1000, b_wins, deps=[
        (trade_E_for_A, currency_A, 500, b_needs),
    ])

    client_B.commit()
    print "Test swap passed"



def test_swap2(rpc_endpoint):
    client_A, currency_A = test_newclient(rpc_endpoint)
    client_B, currency_B = test_newclient(rpc_endpoint)

    client_A.mint(1000)
    client_B.mint(1000)

    client_A.commit()

    assert client_A.balance(currency_A) == 1000
    assert client_B.balance(currency_B) == 1000

    b_sends = os.urandom(32)
    a_sends = os.urandom(32)

    client_B.pay(client_A, currency_B, 1000, b_sends, deps=[
        (client_B, currency_A, 1000, a_sends),
    ])
    client_A.pay(client_B, currency_A, 1000, a_sends, deps=[
        (client_A, currency_B, 1000, b_sends),
    ])

    assert client_B.graph() is True
    client_B.commit()
    print "Test swap 2 passed"


# ------------------------------------------------------------------
# Standalone entrypoint


@click.command()
@click.option('--inproc', is_flag=True, help="Use in-process chain")
@click.argument('endpoint', required=False)
def tests(inproc, endpoint=None):
    """
    Connect to Ion RPC server and perform tests.

    :param inproc: Use in-process API
    :param endpoint: http URL for JSON-RPC endpoint
    """
    if inproc:
        endpoint = IonRpcServer()

    test_pay(endpoint)
    test_swap(endpoint)
    test_swap2(endpoint)
    sys.exit()


@click.command()
@click.option('--ion-rpc', default='127.0.0.1:8545', help='Ethereum JSON-RPC HTTP endpoint', callback=arg_ethrpc)
@click.option('--ion-account', help='Ethereum account address')
@click.option('--ion-contract', help='IonLink contract address')
@click.option('--secret', help="Secret to use in transaction")
@click.argument('endpoint', required=False)
def main(ion_rpc, ion_account, ion_contract, secret, endpoint=None):
    """
    Connect to Ion RPC server and perform payments on chain.

    :param inproc: Use in-process API
    :param endpoint: http URL for JSON-RPC endpoint
    :param secret: secret to be used in transaction
    """

    print(ion_rpc.net_version())
    if not ion_contract or not ion_account:
        print("IonLink disabled")
        ionlink = None
    else:
        if not ion_rpc:
            ion_rpc = EthJsonRpc('127.0.0.1', 8545)
        # TODO: load ABI from package resources
        ionlink = ion_rpc.proxy("abi/IonLink.abi", ion_contract, ion_account)

    server = IonRpcServer(ionlink)

    make_payment(server, os.urandom(32))

@click.command()
@click.option('--ion-rpc', default='127.0.0.1:8545', help='Ethereum JSON-RPC HTTP endpoint', callback=arg_ethrpc)
@click.option('--ion-account', help='Ethereum account address')
@click.option('--ion-contract', help='IonLink contract address')
@click.argument('endpoint', required=False)
def tree(ion_rpc, ion_account, ion_contract, endpoint=None):
    """
    Connect to Ion RPC server and perform payments on chain.

    :param inproc: Use in-process API
    :param endpoint: http URL for JSON-RPC endpoint
    """

    print(ion_rpc.net_version())
    if not ion_contract or not ion_account:
        print("IonLink disabled")
        ionlink = None
    else:
        if not ion_rpc:
            ion_rpc = EthJsonRpc('127.0.0.1', 8545)
        # TODO: load ABI from package resources
        ionlink = ion_rpc.proxy("abi/IonLink.abi", ion_contract, ion_account)

    server = IonRpcServer(ionlink)
    print_ion_history_tree(server)


commands = click.Group('commands', help="RPC Client")
commands.add_command(tests, "test")
commands.add_command(main, "main")
commands.add_command(tree, "get_tree")

if __name__ == "__main__":
    commands.main()
