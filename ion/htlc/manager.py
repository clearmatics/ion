## Copyright (c) 2018 Harry Roberts. All Rights Reserved.
## SPDX-License-Identifier: LGPL-3.0+

import os
import time
from binascii import hexlify, unhexlify
from hashlib import sha256

from ..utils import normalise_address
from ..ethrpc import EthJsonRpc

from .common import MINIMUM_EXPIRY_DURATION, make_htlc_proxy


class ExchangeError(Exception):
    pass


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

        # Verify expiry time is acceptable
        # XXX: should minimum expiry be left to the contract, or the coordinator?
        now = int(time.time())
        min_expiry = now + MINIMUM_EXPIRY_DURATION
        if expiry < min_expiry:
            raise ExchangeError("Expiry too short")

        # GUID used for the exchanges
        # offer_guid = Deposit() by B (the proposer)
        offer_guid = sha256(unhexlify(exch['offer_address']) + unhexlify(secret_hashed)).digest()
        # taker_guid = Deposit() by A (the initial offerer)
        taker_guid = sha256(unhexlify(depositor) + unhexlify(secret_hashed)).digest()

        # TODO: make sure we have the right contract address
        contract = make_htlc_proxy(self._rpc, exch['want_htlc_address'])

        # Verify on-chain expiry matches
        onchain_expiry = contract.GetExpiry(offer_guid)
        if expiry != onchain_expiry:
            raise ExchangeError("Submitted expiry doesn't match on-chain data")

        # Verify on-chain hashed secret
        onchain_sechash = contract.GetSecretHashed(offer_guid)
        if hexlify(onchain_sechash).decode('ascii') != secret_hashed:
            raise ExchangeError("Submitted hashed secret doesn't match on-chain data")

        # 1 = Deposited
        if contract.GetState(offer_guid) != 1:
            raise ExchangeError("Exchange is in wrong state")

        # Verify receiver
        onchain_receiver = normalise_address(contract.GetReceiver(offer_guid))
        if onchain_receiver != exch['offer_address']:
            raise ExchangeError("Offer address differs from receiver")

        # Verify sender
        onchain_sender = normalise_address(contract.GetSender(offer_guid))
        if onchain_sender != depositor:
            raise ExchangeError("Sender address differs from depositor")

        # Ensure deposited amount is more or greater than what was wanted
        onchain_amount = contract.GetAmount(offer_guid)
        if onchain_amount < exch['want_amount']:
            raise ExchangeError("Propose amount differs from want amount")

        # Store proposal
        proposal = dict(
            secret_hashed=secret_hashed,
            expiry=expiry,
            depositor=depositor,
            offer_guid=hexlify(offer_guid).decode('ascii'),
            taker_guid=hexlify(taker_guid).decode('ascii'),
        )
        exch['proposals'][secret_hashed] = proposal

        return exch, proposal

    def confirm(self, exch_guid, secret_hashed, **kwa):
        exch, proposal = self.get_proposal(exch_guid, secret_hashed)

        exch['chosen_proposal'] = secret_hashed

        contract = make_htlc_proxy(self._rpc, exch['offer_htlc_address'])

        taker_guid = unhexlify(proposal['taker_guid'])

        # 1 = Deposited
        if contract.GetState(taker_guid) != 1:
            raise ExchangeError("Exchange is in wrong state")

        # Ensure deposited amount is more or greater than what was wanted
        onchain_amount = contract.GetAmount(taker_guid)
        if onchain_amount < exch['offer_amount']:
            raise ExchangeError("Confirm amount differs from offer amount")

        # Verify on-chain hashed secret
        onchain_sechash = contract.GetSecretHashed(taker_guid)
        if hexlify(onchain_sechash).decode('ascii') != secret_hashed:
            raise ExchangeError("Submitted hashed secret doesn't match on-chain data")

        # Verify on-chain expiry matches
        expiry = proposal['expiry']
        onchain_expiry = contract.GetExpiry(taker_guid)
        if expiry != onchain_expiry:
            raise ExchangeError("Submitted expiry doesn't match on-chain data")

        # Verify receiver
        onchain_receiver = normalise_address(contract.GetReceiver(taker_guid))
        if onchain_receiver != proposal['depositor']:
            raise ExchangeError("Offer address differs from receiver")

        # Verify sender
        deposited_for = exch['offer_address']
        onchain_sender = normalise_address(contract.GetSender(taker_guid))
        if onchain_sender != deposited_for:
            raise ExchangeError("Sender address differs from depositor")
        

    def release(self, exch_guid, secret_hashed, **kwa):
        exch, proposal = self.get_proposal(exch_guid, secret_hashed)

        secret_hex = kwa['secret']
        secret = unhexlify(secret_hex)
        secret_hashed_check = sha256(secret).digest()
        secret_hashed_check_hex = hexlify(secret_hashed_check).decode('ascii')

        if secret_hashed_check_hex != secret_hashed:
            raise ExchangeError(' '.join(["Secret doesn't match! Got", secret_hashed_check_hex, 'expected', secret_hashed]))

        contract = make_htlc_proxy(self._rpc, exch['want_htlc_address'])

        proposal['secret'] = secret_hex

    def finish(self, exch_guid, secret_hashed, **kwa):
        exch, proposal = self.get_proposal(exch_guid, secret_hashed)

        contract = make_htlc_proxy(self._rpc, exch['offer_htlc_address'])