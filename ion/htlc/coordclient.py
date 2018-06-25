# Copyright (c) 2018 Harry Roberts. All Rights Reserved.
# SPDX-License-Identifier: LGPL-3.0+

import os
from hashlib import sha256
from binascii import hexlify, unhexlify

from ..utils import require, normalise_address
from ..ethrpc import EthJsonRpc
from ..restclient import RestClient

from .common import get_default_expiry, make_htlc_proxy


class Proposal(object):
    def __init__(self, coordapi, exch_obj, resource, data=None):
        assert isinstance(coordapi, CoordinatorClient)
        assert isinstance(exch_obj, Exchange)
        assert isinstance(resource, RestClient)
        self._coordapi = coordapi
        self._exch_obj = exch_obj
        self._resource = resource
        self._data = data
        if data is None:
            self.refresh()

    def refresh(self):
        self._data = self._resource.GET()

    @property
    def data(self):
        return self._data

    @property
    def secret(self):
        return self._data.get('secret')

    @property
    def secret_hashed(self):
        return self._data['secret_hashed']

    def confirm(self, wait=True):
        """
        Confirm the exchange by depositing your side of the deal
        You must be the original deal offerer to confirm the exchange
        It must be locked for the same parameters as the proposal
        """
        ethrpc = self._coordapi.ethrpc
        my_address = self._coordapi.my_address
        exch_data = self._exch_obj.data

        require(my_address == exch_data['offer_address'],
                ' '.join(["Only offerer can confirm", my_address, '!=', exch_data['offer_address']]))

        # XXX: one side of the expiry must be twice as long as the other to handle failure case
        conf_expiry = self._data['expiry']
        conf_receiver = self._data['depositor']
        conf_secret_hashed = self._data['secret_hashed']
        conf_value = exch_data['offer_amount']

        # Offerer deposits their side of the deal, locked to same hashed secret 
        htlc_address = exch_data['offer_htlc_address']
        htlc_contract = make_htlc_proxy(ethrpc, htlc_address, my_address)
        txn = htlc_contract.Deposit(conf_receiver, unhexlify(conf_secret_hashed), conf_expiry, value=conf_value)

        receipt = txn.receipt(wait=wait)
        if receipt and int(receipt['status'], 16) == 0:
            raise RuntimeError("Confirm failed, txn: " + txn.txid)

        self._resource.confirm.POST(
            txid=txn.txid
        )

        return txn

    def release(self, secret, wait=True):
        """
        Reveal the secret by withdrawing from your side of the exchange
        This must be performed by party B (the proposer)
        """
        assert len(secret) == 32
        # TODO: make this callable by either side, determine which one we are?
        ethrpc = self._coordapi.ethrpc
        my_address = self._coordapi.my_address
        exch_data = self._exch_obj.data

        secret_hex = hexlify(secret).decode('ascii')
        secret_hashed = sha256(secret).digest()
        secret_hashed_hex = hexlify(secret_hashed).decode('ascii')

        require(my_address == self._data['depositor'], "Only proposer can release")
        require(self._data['secret_hashed'] == secret_hashed_hex, "Secrets don't match!")

        exch_guid = unhexlify(self._data['taker_guid'])

        htlc_address = exch_data['offer_htlc_address']
        htlc_contract = make_htlc_proxy(ethrpc, htlc_address, my_address)
        txn = htlc_contract.Withdraw(exch_guid, secret)
        txn.wait()

        # Reveal secret, posting back to server
        self._resource.release.POST(
            secret=secret_hex,
            txid=txn.txid
        )

    def finish(self, wait=True):
        """
        After the secret has been revealed by the proposer the offerer
        can withdraws funds using the same secret.
        """
        ethrpc = self._coordapi.ethrpc
        my_address = self._coordapi.my_address
        exch_data = self._exch_obj.data

        secret_hex = self._data['secret']
        secret = unhexlify(secret_hex)

        # TODO: verify secret hashes to hashed image

        exch_guid = unhexlify(self._data['offer_guid'])

        htlc_address = exch_data['want_htlc_address']
        htlc_contract = make_htlc_proxy(ethrpc, htlc_address, my_address)
        txn = htlc_contract.Withdraw(exch_guid, secret)

        receipt = txn.wait()
        if receipt and int(receipt['status'], 16) == 0:
            raise RuntimeError("Finish failed, txn: " + txn.txid)

        self._resource.finish.POST(
            txid=txn.txid,
        )

    def refund(self):
        """
        After the timeout has expired the funds can be withdrawn by the original depositor
        """
        ethrpc = self._coordapi.ethrpc
        my_address = self._coordapi.my_address
        exch_data = self._exch_obj.data

        secret_hashed_hex = self._data['secret_hashed']
        secret_hashed = secret_hashed_hex.decode('hex')

        # TODO: detertmine which side we're on, automagically call correct one
        if False:
            htlc_address = self._data['depositor']
        else:
            htlc_address = exch_data['want_htlc_address']
    
        htlc_contract = make_htlc_proxy(ethrpc, htlc_address, my_address)
        htlc_contract.Refund(secret_hashed)


class Exchange(object):
    def __init__(self, coordapi, resource, data_dict=None):
        assert isinstance(coordapi, CoordinatorClient)
        assert isinstance(resource, RestClient)
        self._coordapi = coordapi
        self._resource = resource
        self._data = None
        if data_dict is None:
            self.refresh()
        else:
            self.set_data(data_dict)

    def _make_proposal(self, secret_hashed_hex, propdata=None):
        return Proposal(self._coordapi, self, self._resource(secret_hashed_hex), propdata)

    @property
    def guid(self):
        return self._data['guid']

    @property
    def data(self):
        return self._data

    def set_data(self, value):
        value['proposals'] = {key: self._make_proposal(key, propdata)
                              for key, propdata in value['proposals'].items()}
        self._data = value

    def refresh(self):
        self.set_data(self._resource.GET())

    @property
    def chosen_proposal(self):
        prop_id = self._data['chosen_proposal']
        if prop_id:
            return self.proposal(prop_id)
        return None

    @property
    def proposals(self):
        return self._data['proposals']

    def proposal(self, secret_hashed_hex):
        return self.proposals[secret_hashed_hex]

    def propose(self):
        """
        Submit a proposal for the exchange by depositing your tokens
        into a HTLC contract.
        """
        # Create a random secret
        secret = os.urandom(32)
        secret_hashed = sha256(secret).digest()
        secret_hashed_hex = hexlify(secret_hashed).decode('ascii')

        # Proposal parameters
        prop_receiver = self._data['offer_address']
        prop_value = self._data['want_amount']
        prop_expiry = get_default_expiry()

        # Perform deposit
        ethrpc = self._coordapi.ethrpc
        my_address = self._coordapi.my_address

        # TODO: verify adequate balance to cover the deposit

        require(my_address != prop_receiver, "Cannot be both sides of exchange")

        htlc_address = self._data['want_htlc_address']
        htlc_contract = make_htlc_proxy(ethrpc, htlc_address, my_address)
        txn = htlc_contract.Deposit(prop_receiver, secret_hashed, prop_expiry, value=prop_value)

        receipt = txn.receipt(wait=True)
        if receipt and int(receipt['status'], 16) == 0:
            raise RuntimeError("Propose deposit failed, txn: " + txn.txid)

        # Notify coordinator of proposal
        proposal_resource = self._resource(secret_hashed_hex)
        propdata = proposal_resource.POST(
            expiry=prop_expiry,
            depositor=my_address,
            txid=txn.txid
        )

        # Add proposal to list, then return it
        proposal = self._make_proposal(secret_hashed_hex, propdata)
        self.proposals[secret_hashed_hex] = proposal
        return secret, proposal


class CoordinatorClient(object):
    """
    Uses the coordinator API and the HTLC Ethereum contract to
    perform a cross-chain atomic swap as easily as possible.

    essentially the glue between the coordinator and the contract.
    """
    def __init__(self, my_address, ethrpc, api_url=None, resource=None):
        self._my_address = normalise_address(my_address)
        self._resource = RestClient(api_url) if resource is None else resource
        self._ethrpc = ethrpc
        assert isinstance(self._resource, RestClient)
        assert isinstance(ethrpc, EthJsonRpc)

    @property
    def ethrpc(self):
        return self._ethrpc

    @property
    def my_address(self):
        return self._my_address

    def list(self):
        results = self._resource.list.GET()
        retval = list()
        for exch_id, data in results.items():
            resource = self._resource(exch_id)
            exch = Exchange(self, resource, data)
            retval.append(exch)
        return retval

    def get_exchange(self, exch_id, exch_data=None):
        assert len(exch_id) == 40
        resource = self._resource(exch_id)
        return Exchange(self, resource, exch_data)

    def advertise(self, offer_amount, want_amount):
        resp = self._resource.advertise.POST(
            offer_address=self._my_address,
            offer_amount=offer_amount,
            want_amount=want_amount)
        return self.get_exchange(resp['id'])
