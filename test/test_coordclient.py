#!/usr/bin/env python3
## Copyright (c) 2018 Harry Roberts.
## SPDX-License-Identifier: LGPL-3.0+

import sys

from ion.ethrpc import EthJsonRpc
from ion.htlc.coordclient import CoordinatorClient


def main(api_url='http://127.0.0.1:5000/htlc'):
    # Connect to Ethereum instance
    ethrpc = EthJsonRpc('127.0.0.1', 8545)
    accounts = ethrpc.eth_accounts()
    assert len(accounts) > 1
    addr_A = accounts[0]
    addr_B = accounts[1]
    balance_A_start = ethrpc.eth_getBalance(addr_A)
    balance_B_start = ethrpc.eth_getBalance(addr_B)

    offer = 1000
    want = 500

    # Setup a coordinator client for each A and B
    client_A = CoordinatorClient(addr_A, ethrpc, api_url)
    client_B = CoordinatorClient(addr_B, ethrpc, api_url)

    # A advertises exchange, offers 1000 for 500
    exch_A = client_A.advertise(offer, want)

    # B retrieves exchanges
    exchanges_B = client_B.list()
    print("B list", exchanges_B)

    # Verify details of the exchange
    exch_B = [_ for _ in exchanges_B if _.guid == exch_A.guid][0]
    exch_B_data = exch_B.data
    print("Data", exch_B.data)

    # B decides to participate in the swap
    # Proposing deposits the funds
    secret, prop_B = exch_B.propose()
    print("Proposal is", prop_B)
    balance_B_deposited = ethrpc.eth_getBalance(addr_B)

    # A then confirms the proposal, depositing their side of the deal
    # This 'locks-in' the trade between A and B, denying any further proposals
    exch_A.refresh()
    print("Updated data:", exch_A.data)
    prop_A = exch_A.proposal(prop_B.secret_hashed)
    prop_A.confirm()

    # Verify proposal has been chosen
    exch_A.refresh()
    exch_B.refresh()
    print("Chosen proposal is", exch_A.chosen_proposal)

    # B side then releases the secret, withdrawing the funds A deposited
    prop_B.release(secret)

    balance_B_released = ethrpc.eth_getBalance(addr_B)
    print("B Balance start", balance_B_start)
    print("B Balance deposited", balance_B_deposited)
    print("B Balance released", balance_B_released)
    print("B Difference",  balance_B_deposited - balance_B_released)

    # Then A finishes, using the released secret to withdraw the funds B deposited
    prop_A.refresh()
    print("Prop A data", prop_A.data)
    prop_A.finish()


if __name__ == "__main__":
    main(*sys.argv[1:])
