## Copyright (c) 2018 Harry Roberts. All Rights Reserved.
## SPDX-License-Identifier: LGPL-3.0+

import sys

from flask import Flask, Blueprint, request, jsonify

from ..ethrpc import EthJsonRpc
from ..utils import normalise_address
from ..webutils import (Bytes32Converter, Bytes20Converter, params_parse, api_abort,
                        param_bytes20, param_bytes32, param_uint256)

from .manager import ExchangeManager, ExchangeError


class CoordinatorBlueprint(Blueprint):
    """
    Provides a web API for coordinating cross-chain HTLC exchanges
    """

    def __init__(self, htlc_address, rpc, **kwa):
        Blueprint.__init__(self, 'htlc', __name__, **kwa)

        self._manager = ExchangeManager(htlc_address, rpc)

        # XXX: This sure looks hacky... assignment not allowed in lambda, callback etc.
        self.record(lambda s: s.app.url_map.converters.__setitem__('bytes32', Bytes32Converter))
        self.record(lambda s: s.app.url_map.converters.__setitem__('bytes20', Bytes20Converter))

        self.add_url_rule("/", 'index', self.index)
        self.add_url_rule("/list", 'list', self.index)
        self.add_url_rule("/advertise", 'advertise', self.exch_advertise, methods=['POST'])
        self.add_url_rule("/<bytes20:exch_id>", 'get', self.exch_get, methods=['GET'])
        self.add_url_rule("/<bytes20:exch_id>/<bytes32:secret_hashed>", 'proposal_get',
                          self.exch_proposal_get, methods=['GET'])
        self.add_url_rule("/<bytes20:exch_id>/<bytes32:secret_hashed>", 'propose',
                          self.exch_propose, methods=['POST'])
        self.add_url_rule("/<bytes20:exch_id>/<bytes32:secret_hashed>/confirm", 'confirm',
                          self.exch_confirm, methods=['POST'])
        self.add_url_rule("/<bytes20:exch_id>/<bytes32:secret_hashed>/release", 'release',
                          self.exch_release, methods=['POST'])
        self.add_url_rule("/<bytes20:exch_id>/<bytes32:secret_hashed>/finish", 'finish',
                          self.exch_finish, methods=['POST'])

    def index(self):
        """
        Display list of all exchanges, and their details
        """
        return jsonify(self._manager.exchanges)

    def exch_advertise(self):
        """
        Advertise a potential exchange, advertiser offers N of X, wants M of Y
        The address of the 'offer' side is advertised, as only that address can withdraw.

        This is performed by Alice
        """

        params = params_parse(request.form, dict(
            offer_address=param_bytes20,
            offer_amount=param_uint256,
            want_amount=param_uint256
        ))

        try:
            exch_guid = self._manager.advertise(**params)
        except ExchangeError as ex:
            return api_abort(str(ex))

        return jsonify(dict(
            id=exch_guid,
            ok=1
        ))

    def exch_get(self, exch_id):
        """
        Retrieve details of exchange
        """
        try:
            exch = self._manager.get_exchange(exch_id)
        except ExchangeError as ex:
            return api_abort(str(ex))
        return jsonify(exch)

    def exch_propose(self, exch_id, secret_hashed):
        """
        Somebody who has what the offerer wants proves they have deposited it.
        They post proof as a proposal for exchange, this includes a hash of the secret they chose.
        They post details of this deposit as the proposal.

        This is performed by Bob
        """
        params = params_parse(request.form, dict(
            expiry=param_uint256,
            depositor=param_bytes20,
            txid=param_bytes32,
        ))

        try:
            _, proposal = self._manager.propose(exch_id, secret_hashed, **params)
        except ExchangeError as ex:
            return api_abort(str(ex))

        # TODO: redirect to proposal URL? - or avoid another GET request...
        return jsonify(proposal)

    def exch_proposal_get(self, exch_id, secret_hashed):
        """
        Retrieve details for a specific exchange proposal
        """
        try:
            _, proposal = self._manager.get_proposal(exch_id, secret_hashed)
        except ExchangeError as ex:
            return api_abort(str(ex))
        return jsonify(proposal)

    def exch_confirm(self, exch_id, secret_hashed):
        """
        The initial offerer (A), who wants M of Y and offers N of X
        Upon seeing the proposal by B for M of Y, posts their confirmation of the specific deal
        The offerer (A) must have deposited N of X, using the same hash image that (B) proposed

        This is performed by Alice
        """
        params = params_parse(request.form, dict(
            txid=param_bytes32,
        ))

        try:
            self._manager.confirm(exch_id, secret_hashed, **params)
        except ExchangeError as ex:
            return api_abort(str(ex))

        # TODO: return updated proposal object
        return jsonify(dict(
            ok=1
        ))

    def exch_release(self, exch_id, secret_hashed):
        """
        Bob releases the secret to withdraw funds that Alice deposited.

        This is performed by Bob
        """
        params = params_parse(request.form, dict(
            secret=param_bytes32,
            txid=param_bytes32,
        ))

        try:
            self._manager.release(exch_id, secret_hashed, **params)
        except ExchangeError as ex:
            return api_abort(str(ex))

        return jsonify(dict(
            ok=1
        ))

    def exch_finish(self, exch_id, secret_hashed):
        """
        Alice uses the secret revealed by Bob to withdraw the funds he deposited.

        This is performed by Alice.

        This completes the exchange.
        """
        params = params_parse(request.form, dict(
            txid=param_bytes32,
        ))

        try:
            self._manager.finish(exch_id, secret_hashed, **params)
        except ExchangeError as ex:
            return api_abort(str(ex))

        return jsonify(dict(
            ok=1
        ))


def main(htlc_address, rpc=None):
    """
    Development server for coordinator

    NOTE: not suitable for 'production'
    SEE: http://flask.pocoo.org/docs/1.0/deploying/#deployment
    """
    htlc_address = normalise_address(htlc_address)
    if rpc is None:
        rpc = EthJsonRpc()

    coordinator = CoordinatorBlueprint(htlc_address, rpc)

    app = Flask(__name__)
    app.register_blueprint(coordinator, url_prefix='/htlc')

    # NOTE: Flask reloader doesn't work well with packages *shakes-fists*
    app.run(use_reloader=False)
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1]))
