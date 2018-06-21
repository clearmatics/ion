## Copyright (c) 2018 Harry Roberts. All Rights Reserved.
## SPDX-License-Identifier: LGPL-3.0+

import time
from binascii import hexlify, unhexlify

from ..utils import normalise_address, require

from .common import MINIMUM_EXPIRY_DURATION, make_htlc_proxy


def verify_deposit(side, rpc, exch, proposal, txid):
    """
    Verifies that the contract deposit matches the exchange and proposal
    """
    require(side in ['proposer', 'confirmer'], "Side must be 'proposer' or 'confirmer'")

    expiry = proposal['expiry']
    secret_hashed = proposal['secret_hashed']

    # Verification logic is the same, but uses different parameters
    # depending on which side, the proposer or the confirmer
    if side == 'proposer':
        deposit_guid = unhexlify(proposal['offer_guid'])
        htlc_address = exch['want_htlc_address']
        expected_amount = exch['want_amount']
        expected_receiver = exch['offer_address']
        expected_sender = proposal['depositor']
    else:
        deposit_guid = unhexlify(proposal['taker_guid'])
        htlc_address = exch['offer_htlc_address']
        expected_amount = exch['offer_amount']
        expected_receiver = proposal['depositor']
        expected_sender = exch['offer_address']

    rpc.receipt_wait(txid)

    contract = make_htlc_proxy(rpc, htlc_address)

    # Verify expiry time is acceptable
    # XXX: should minimum expiry be left to the contract, or the coordinator?
    now = int(time.time())
    min_expiry = now + MINIMUM_EXPIRY_DURATION
    if expiry < min_expiry:
        raise ExchangeError("Expiry too short, got %d expected >= %d" % (
            expiry, min_expiry))

    # Verify on-chain expiry matches
    onchain_expiry = contract.GetExpiry(deposit_guid)
    if expiry != onchain_expiry:
        raise ExchangeError("Expiry doesn't match contract, got %d expected %d" % (
            expiry, onchain_expiry))

    # Verify on-chain hashed secret
    onchain_sechash = hexlify(contract.GetSecretHashed(deposit_guid)).decode('ascii')
    if onchain_sechash != secret_hashed:
        raise ExchangeError("Hashed secret doesn't match contract, got %s expected %s" % (
            onchain_sechash, secret_hashed))

    # 1 = Deposited
    onchain_state = contract.GetState(deposit_guid)
    if onchain_state != 1:
        raise ExchangeError("Exchange is in wrong state, got %d expected %d" % (
            onchain_state, 1))

    # Verify receiver
    onchain_receiver = normalise_address(contract.GetReceiver(deposit_guid))
    if onchain_receiver != expected_receiver:
        raise ExchangeError("Wrong receiver address, got %s expected %s" % (
            onchain_receiver, expected_receiver))

    # Verify sender
    onchain_sender = normalise_address(contract.GetSender(deposit_guid))
    if onchain_sender != expected_sender:
        raise ExchangeError("Wrong sender address, got %s expected %s" % (
            onchain_sender, expected_sender))

    # Ensure deposited amount is more or greater than what was wanted
    onchain_amount = contract.GetAmount(deposit_guid)
    if onchain_amount < expected_amount:
        raise ExchangeError("Propose amount differs from want amount, got %d expected %d" % (
            onchain_amount, expected_amount))

    return True