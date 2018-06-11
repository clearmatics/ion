# Copyright (c) 2018 Harry Roberts. All Rights Reserved.
# SPDX-License-Identifier: LGPL-3.0+

import os
from hashlib import sha256

from ethereum.utils import keccak

from ..restclient import RestClient

from .common import get_default_expiry, make_htlc_proxy


DEPOSIT_SIGNATURE = (keccak.new(digest_bits=256)
                     .update('Deposit(address,bytes32,uint256)')
                     .hexdigest()[:8])


class Proposal(object):
    def __init__(self, client, resource, data=None):
        self._client = client
        self._resource = resource
        self._data = data
        if data is None:
            self.refresh()

    def refresh(self):
        self._data = self._resource.GET()

    def confirm(self):
        pass

    def release(self):
        pass

    def finish(self):
        pass


class Exchange(object):
    def __init__(self, client, resource, data_dict=None):
        self._client = client
        self._resource = resource
        self._data = None
        if data_dict is None:
            self.refresh()
        else:
            self.data = data_dict

    def _make_proposal(self, key, propdata=None):
        return Proposal(self._client, self._resource(key), propdata)

    @property
    def data(self):
        return self._data

    @data.setter
    def set_data(self, value):
        value['proposals'] = [self._make_proposal(key, propdata)
                              for key, propdata in value['proposals'].items()]
        self._data = value

    def refresh(self):
        self.data = self._resource.GET()

    @property
    def proposals(self):
        return self._data['proposals']

    def propose(self):
        """
        Submit a proposal for the exchange by depositing your tokens
        into a HTLC contract.
        """
        # Create a random secret
        secret = os.urandom(32)
        secret_hashed = sha256(secret).digest()

        # Proposal parameters
        prop_receiver = self._data['offer_address']
        prop_value = self._data['want_amount']
        prop_expiry = get_default_expiry()


        # Perform deposit
        ethrpc = self._client.ethrpc
        my_address = self._client.ethrpc

        htlc_address = self._data['want_htlc_contract']
        htlc_contract = make_htlc_proxy(ethrpc, htlc_address, my_address)
        htlc_contract.Deposit(prop_receiver, secret_hashed, prop_expiry, value=prop_value)

        # TODO: wait for deposit to go through?
        #       or provide some kind of receipt...

        # Notify coordinator of proposal
        proposal_resource = self._resource(secret_hashed.encode('hex'))
        response = proposal_resource.POST(
            expiry=prop_expiry,
            depositor=my_address.encode('hex')
        )
        require(response['ok'] == 1, "Proposal coordinator API error")

        proposal = Proposal(self._client, proposal_resource)
        self.proposals.append(proposal)
        return proposal


class CoordinatorClient(object):
    """
    Uses the coordinator API and the HTLC Ethereum contract to
    perform a cross-chain atomic swap as easily as possible.

    essentially the glue between the coordinator and the contract.
    """
    def __init__(self, my_address, ethrpc, api_url=None, resource=None):
        self._my_address = my_address
        self._resource = RestClient(api_url) if resource is None else resource
        self._ethrpc = ethrpc

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
