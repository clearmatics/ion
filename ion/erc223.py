from __future__ import print_function
import sys
import argparse

from ethjsonrpc import EthJsonRpc

from .args import Bytes20Action, EthRpcAction, PosInt256Action
from .solproxy import solproxy


def Token(rpc, contract, account):
    return solproxy(rpc, "abi/Token.abi", contract.encode('hex'), account.encode('hex'))


def erc223_options(args):
    parser = argparse.ArgumentParser(description="Plasma Chain")

    parser.add_argument('-r', '--rpc', metavar="ip:port", dest='rpc', action=EthRpcAction,
                        help='Ethereum RPC address', default='127.0.0.1:8545')

    parser.add_argument('-c', '--contract', metavar="0x...20", dest='contract', action=Bytes20Action,
                        help='ERC-223 contract address', required=True)

    parser.add_argument('-a', '--account', metavar="0x...20", dest='from_account', action=Bytes20Action,
                        help='Ethereum account address', required=True)

    parser.add_argument('-v', '--value', dest='value', action=PosInt256Action,
                        help='Amount of tokens to transfer')

    parser.add_argument('--transfer', action=Bytes20Action)

    parser.add_argument('--balance', action=Bytes20Action)

    parser.add_argument('--mint', action=PosInt256Action)

    opts = parser.parse_args(args or sys.argv[1:])

    if isinstance(opts.rpc, str):
        opts.rpc = EthJsonRpc(*opts.rpc.split(':'))

    if opts.transfer:
        if not opts.value:
            parser.error("--account and --value required for --transfer")

    return opts


def main(args=None):
    opts = erc223_options(args)

    token = Token(opts.rpc, opts.contract, opts.from_account)

    print("RPC Server:", opts.rpc)
    print("Contract:", opts.contract.encode('hex'))
    print("Account:", opts.from_account.encode('hex'))
    print("")

    if opts.transfer:
        print("Transfer %r to %r" % (opts.value, opts.transfer.encode('hex')))
        result = token.transfer_a9059cbb(opts.transfer.encode('hex'), opts.value)

    if opts.mint:
        token.mint(opts.mint)

    #print("From acct", opts.from_account.encode('hex'), len(opts.from_account.encode('hex')))
    print("Balance = ", token.balanceOf(opts.balance or opts.from_account.encode('hex')))



if __name__ == "__main__":
    main()
