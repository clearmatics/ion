import os
import time
from hashlib import sha256

from ..utils import normalise_address

from .common import MINIMUM_EXPIRY_DURATION


class ExchangeError(Exception):
    pass


class ExchangeManager(object):
    def __init__(self, htlc_address):
        self._exchanges = dict()
        self._htlc_address = normalise_address(htlc_address)

    @property
    def exchanges(self):
        return self._exchanges

    def get_exchange(self, exch_guid):
        return self._exchanges.get(exch_guid)

    def get_proposal(self, exch_guid, secret_hashed):
        exch = self.get_exchange(exch_guid)
        proposal = exch['proposals'].get(secret_hashed)
        if not proposal:
            raise ExchangeError("Unknown proposal")
        return exch, proposal

    def advertise(self, **kwa):
        exch_guid = os.urandom(20).encode('hex')

        # Save exchange details
        # TODO: replace with class instance, `Exchange` ?

        exch = dict(
            guid=exch_guid,
            offer_address=kwa['offer_address'],
            offer_amount=kwa['offer_amount'],
            want_amount=kwa['want_amount'],
            proposals=dict(),
            chosen_proposal=None,

            # Temporary placeholders
            # TODO: replace with correct contracts
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

        # TODO: verify details on-chain, expiry, depositor and secret must match

        # Verify expiry time is acceptable
        # XXX: should minimum expiry be left to the contract, or the coordinator?
        now = int(time.time())
        min_expiry = now + MINIMUM_EXPIRY_DURATION
        if expiry < min_expiry:
            raise ExchangeError("Expiry too short")

        # GUID used for the exchanges
        # XXX: these are swapped
        # offer_guid = Deposit() by B (the proposer)
        offer_guid = sha256(exch['offer_address'].decode('hex') + secret_hashed.decode('hex')).digest()
        # taker_guid = Deposit() by A (the initial offerer)
        taker_guid = sha256(depositor.decode('hex') + secret_hashed.decode('hex')).digest()

        # Store proposal
        proposal = dict(
            secret_hashed=secret_hashed,
            expiry=expiry,
            depositor=depositor,
            offer_guid=offer_guid.encode('hex'),
            taker_guid=taker_guid.encode('hex'),
        )
        exch['proposals'][secret_hashed] = proposal

        return exch, proposal

    def confirm(self, exch_guid, secret_hashed, **kwa):
        exch, proposal = self.get_proposal(exch_guid, secret_hashed)

        exch['chosen_proposal'] = secret_hashed

    def release(self, exch_guid, secret_hashed, **kwa):
        exch, proposal = self.get_proposal(exch_guid, secret_hashed)

        secret_hex = kwa['secret']
        secret = secret_hex.decode('hex')
        secret_hashed_check = sha256(secret).digest()
        secret_hashed_check_hex = secret_hashed_check.encode('hex')

        if secret_hashed_check_hex != secret_hashed:
            raise ExchangeError(' '.join(["Secret doesn't match! Got", secret_hashed_check_hex, 'expected', secret_hashed]))

        proposal['secret'] = secret_hex

    def finish(self, exch_guid, secret_hashed, **kwa):
        exch, proposal = self.get_proposal(exch_guid, secret_hashed)
