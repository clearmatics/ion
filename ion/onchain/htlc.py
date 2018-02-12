from __future__ import print_function
import sys
import argparse

from ..ethrpc import EthJsonRpc
from ..args import PosInt256, EthRpc, Bytes20, Bytes32


def HTL_Contract(rpc, contract, account):
    return rpc.proxy("abi/HTLC.abi", contract, account)


def htlc_options(args):
    parser = argparse.ArgumentParser(description="HTLC utility")

    parser.add_argument('-r', '--rpc', metavar="ip:port", dest='rpc', action=EthRpc,
                        help='Ethereum RPC address', default='127.0.0.1:8545')

    parser.add_argument('-c', '--contract', metavar="0x...20", dest='contract', action=Bytes20,
                        help='ERC-223 contract address', required=True)

    parser.add_argument('-a', '--account', metavar="0x...20", dest='from_account', action=Bytes20,
                        help='Ethereum account address', required=True)

    subparsers = parser.add_subparsers()

    deposit_group = subparsers.add_parser('deposit')
    deposit_group.add_argument('hash', action=Bytes32)
    deposit_group.add_argument('recipient', action=Bytes20)
    deposit_group.set_defaults(action="deposit")

    claim_group = subparsers.add_parser('claim')
    claim_group.add_argument('-k', '--key', action=Bytes32)
    claim_group.add_argument('-s', '--signature', action=Bytes32)
    claim_group.add_argument('lock_id', action=PosInt256)
    claim_group.add_argument('preimage', action=Bytes20)
    claim_group.set_defaults(action="claim")

    refund_group = subparsers.add_parser('refund')
    refund_group.add_argument('lock_id', action=PosInt256)
    refund_group.add_argument('signature', action=PosInt256)
    refund_group.set_defaults(action="refund")

    opts = parser.parse_args(args or sys.argv[1:])

    if isinstance(opts.rpc, str):
        opts.rpc = EthJsonRpc(*opts.rpc.split(':'))

    return opts


def main(args=None):
    opts = htlc_options(args)

    lock_contract = HTL_Contract(opts.rpc, opts.contract, opts.from_account)

    print("RPC Server:", opts.rpc)
    print("Contract:", opts.contract.encode('hex'))
    print("Account:", opts.from_account.encode('hex'))
    print("")

    if opts.action == "deposit":
        lock_contract.mint()

if __name__ == "__main__":
    main()
