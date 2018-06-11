# Copyright (c) 2018 Harry Roberts. All Rights Reserved.
# SPDX-License-Identifier: LGPL-3.0+

from ..restclient import RestClient


class Proposal(object):
    def __init__(self, restapi, data):
        self._restapi = restapi
        self._data = data

    def confirm(self):
        pass

    def release(self):
        pass

    def finish(self):
        pass


class Exchange(object):
    def __init__(self, restapi, data_dict):
        self._restapi = restapi
        self._data = None
        self.data = data_dict

    @property
    def data(self):
        return self._data

    @data.setter
    def set_data(self, value):
        # TODO: wrap proposals
        self._data = value

    def refresh(self):
        self._data = self._restapi.GET()

    @property
    def proposals(self):
        return self._data['proposals']

    def propose(self):
        # TODO: perform POST request to create proposal
        #       add proposal to list
        #       then return Proposal object
        pass


class CoordinatorClient(object):
    """
    Uses the coordinator API and the HTLC Ethereum contract to
    perform a cross-chain atomic swap as easily as possible.

    essentially the glue between the coordinator and the contract.
    """
    def __init__(self, ethrpc, api_url):
        # TODO: needs parameter for HTLC contract address
        self._restapi = RestClient(api_url)
        self._ethrpc = ethrpc
        #self._htlc = rpc.proxy('abi/HTLC.abi', contract, account)

    def list(self):
        results = self._restapi.list.GET()
        retval = list()
        for exch_id, data in results.items():
            resource = self._restapi(exch_id)
            exch = Exchange(resource, data)
            retval.append(exch)
        return retval

    def advertise(self):
        # TODO: return Exchange object
        pass
