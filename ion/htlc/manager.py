## Copyright (c) 2018 Harry Roberts. All Rights Reserved.
## SPDX-License-Identifier: LGPL-3.0+

import os
import time
from binascii import hexlify, unhexlify
from hashlib import sha256

from ..utils import normalise_address
from ..ethrpc import EthJsonRpc

from .common import MINIMUM_EXPIRY_DURATION, make_htlc_proxy, ExchangeError
from .verify import verify_deposit


class ExchangeManager(object):
    def __init__(self, htlc_address, ethrpc):
        self._exchanges = dict()
        self._htlc_address = normalise_address(htlc_address)
        self._rpc = ethrpc
        assert isinstance(ethrpc, EthJsonRpc)

    @property
    def exchanges(self):
        return self._exchanges

    def get_exchange(self, exch_guid):
        return self._exchanges.get(exch_guid)

    def get_proposal(self, exch_guid, secret_hashed):
        exch = self.get_exchange(exch_guid)
        proposal = exch['proposals'].get(secret_hashed, None)
        if not proposal:
            raise ExchangeError("Unknown proposal")
        return exch, proposal

    def advertise(self, **kwa):
        exch_guid = hexlify(os.urandom(20)).decode('ascii')

        exch = dict(
            guid=exch_guid,
            offer_address=kwa['offer_address'],
            offer_amount=kwa['offer_amount'],
            want_amount=kwa['want_amount'],
            proposals=dict(),
            chosen_proposal=None,

            # Temporary placeholders
            # TODO: replace with correct contracts depending on network
            offer_htlc_address=self._htlc_address,
            want_htlc_address=self._htlc_address
        )

        self._exchanges[exch_guid] = exch

        return exch_guid

    def propose(self, exch_guid, secret_hashed, **kwa):
        exch = self.get_exchange(exch_guid)

        expiry = kwa['expiry']
        depositor = kwa['depositor']

        if exch['chosen_proposal']:
            raise ExchangeError("Proposal has already been chosen")

        # Hashed secret is the 'image', pre-image can be supplied to prove knowledge of secret
        if secret_hashed in exch['proposals']:
            raise ExchangeError("Duplicate proposal secret")

        # GUID used for the exchanges
        # offer_guid = Deposit() by B (the proposer)
        offer_guid = sha256(unhexlify(exch['offer_address']) + unhexlify(secret_hashed)).digest()
        # taker_guid = Deposit() by A (the initial offerer)
        taker_guid = sha256(unhexlify(depositor) + unhexlify(secret_hashed)).digest()

        # Wait for transaction success
        txid = kwa['txid']

        proposal = dict(
            secret_hashed=secret_hashed,
            expiry=expiry,
            depositor=depositor,
            offer_guid=hexlify(offer_guid).decode('ascii'),
            taker_guid=hexlify(taker_guid).decode('ascii'),
            propose_txid=txid,
        )

        verify_deposit('proposer', self._rpc, exch, proposal, txid)

        exch['proposals'][secret_hashed] = proposal
        return exch, proposal

    def confirm(self, exch_guid, secret_hashed, **kwa):
        exch, proposal = self.get_proposal(exch_guid, secret_hashed)

        txid = kwa['txid']

        verify_deposit('proposer', self._rpc, exch, proposal, txid)

        proposal['confirm_txid'] = txid
        exch['chosen_proposal'] = secret_hashed


    def release(self, exch_guid, secret_hashed, **kwa):
        exch, proposal = self.get_proposal(exch_guid, secret_hashed)

        secret_hex = kwa['secret']
        secret = unhexlify(secret_hex)
        secret_hashed_check = sha256(secret).digest()
        secret_hashed_check_hex = hexlify(secret_hashed_check).decode('ascii')
        if secret_hashed_check_hex != secret_hashed:
            raise ExchangeError(' '.join(["Secret doesn't match! Got", secret_hashed_check_hex, 'expected', secret_hashed]))

        # Wait for transaction success
        txid = kwa['txid']
        self._rpc.receipt_wait(txid)

        contract = make_htlc_proxy(self._rpc, exch['want_htlc_address'])
        # XXX: if the server errors out here... then proposal won't get updated, this is bad!

        # 2 = Withdrawn
        offer_guid = unhexlify(proposal['offer_guid'])
        onchain_state = contract.GetState(offer_guid)
        print("After release, State is ", onchain_state)
        """
        # XXX: even though we've waited for a successful receipt, the state is still `1`
        #      but in the 'finish' call, the state with the same params is `2`
        if onchain_state != 2:
            raise ExchangeError("Exchange is in wrong state")
        """

        proposal['secret'] = secret_hex
        proposal['release_txid'] = txid

    def finish(self, exch_guid, secret_hashed, **kwa):
        exch, proposal = self.get_proposal(exch_guid, secret_hashed)

        contract = make_htlc_proxy(self._rpc, exch['offer_htlc_address'])

        taker_guid = unhexlify(proposal['taker_guid'])

        txid = kwa['txid']
        self._rpc.receipt_wait(txid)

        # 2 = Withdrawn
        onchain_state = contract.GetState(taker_guid)
        if onchain_state != 2:
            raise ExchangeError("Exchange is in wrong state")

        proposal['finish_txid'] = txid
