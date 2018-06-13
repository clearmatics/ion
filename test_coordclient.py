#!/usr/bin/env python
## Copyright (c) 2018 Harry Roberts.
## SPDX-License-Identifier: LGPL-3.0+

import sys

from ion.ethrpc import EthJsonRpc
from ion.htlc.coordclient import CoordinatorClient


def main(api_url='http://127.0.0.1:5000/htlc'):
    ethrpc = EthJsonRpc('127.0.0.1', 8545)
    accounts = ethrpc.eth_accounts()
    assert len(accounts) > 1

    addr_A = accounts[0]
    addr_B = accounts[1]

    client_A = CoordinatorClient(addr_A, ethrpc, api_url)
    client_B = CoordinatorClient(addr_B, ethrpc, api_url)

    balance_A = ethrpc.eth_getBalance(addr_A)
    balance_B = ethrpc.eth_getBalance(addr_B)

    # A advertises exchange, offers 1000 for 500
    client_A.advertise(1000, 500)

    # B retrieves exchanges
    exchanges_B = client_B.list()
    print("B list", exchanges_B)

    exch = exchanges_B[0]
    exch_data = exch.data
    print("Data", exch.data)

    proposal = exch.propose()
    print("Proposal is", proposal)


if __name__ == "__main__":
    main(*sys.argv[1:])
