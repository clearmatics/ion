## Copyright (c) 2018 Harry Roberts. All Rights Reserved.
## SPDX-License-Identifier: LGPL-3.0+

import sys
import os
import json
import time

from flask import Flask, Blueprint, request, abort, jsonify, make_response
from werkzeug.routing import BaseConverter

from ..args import arg_bytes32, arg_bytes20, arg_uint256
from ..utils import scan_bin

from .common import MINIMUM_EXPIRY_DURATION


def api_abort(message, code=400):
    return abort(make_response(jsonify(dict(_error=message)), code))


def param(the_dict, key):
    if key not in the_dict:
        return api_abort("Parameter required: " + key)
    return the_dict[key]


def param_filter_arg(the_dict, key, filter_fn):
    """
    Applies a click argument filter from `..args` to a value from a dictionary
    Does the things with HTTP errors etc. upon failure.
    """
    value = param(the_dict, key)
    try:
        value = filter_fn(None, None, value)
    except Exception as ex:
        return api_abort("Invalid parameter '%s' - %s" % (key, str(ex)))
    return value


def param_bytes32(the_dict, key):
    return param_filter_arg(the_dict, key, arg_bytes32).encode('hex')


def param_bytes20(the_dict, key):
    return param_filter_arg(the_dict, key, arg_bytes20).encode('hex')


def param_uint256(the_dict, key):
    return param_filter_arg(the_dict, key, arg_uint256)


class BytesConverter(BaseConverter):
    """
    Accepts hex encoded bytes as an argument
    Provides raw bytes to Python
    Marshals between raw bytes and hex encoded
    """
    BYTES_LEN = None
    def __init__(self, url_map, *items):
        assert self.BYTES_LEN is not None
        super(BytesConverter, self).__init__(url_map)
        self.regex = '(0x)?[a-fA-F0-9]{' + str(self.BYTES_LEN * 2) + '}'

    def to_python(self, value):
        # Normalise hex encoding...
        return scan_bin(value).encode('hex')


class Bytes32Converter(BytesConverter):
    """Accept 32 hex-encoded bytes as URL param"""
    BYTES_LEN = 32


class Bytes20Converter(BytesConverter):
    """Accept 20 hex-encoded bytes as URL param"""
    BYTES_LEN = 20


class CoordinatorBlueprint(Blueprint):
    """
    Provides a web API for coordinating cross-chain HTLC exchanges
    """

    def __init__(self, htlc_address, **kwa):
        Blueprint.__init__(self, 'htlc', __name__, **kwa)

        self._htlc_address = htlc_address
        self._exchanges = dict()

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

    def _get_exch(self, exch_id):
        if exch_id not in self._exchanges:
            return api_abort("Unknown exchange", code=404)
        return self._exchanges[exch_id]

    def _get_proposal(self, exch_id, secret_hashed):
        exch = self._get_exch(exch_id)
        proposal = exch['proposals'].get(secret_hashed)
        if not proposal:
            return api_abort("Unknown proposal", code=404)
        return exch, proposal

    def index(self):
        """
        Display list of all exchanges, and their details
        """
        return jsonify(self._exchanges)

    def exch_advertise(self):
        """
        Advertise a potential exchange, advertiser offers N of X, wants M of Y
        The address of the 'offer' side is advertised, as only that address can withdraw.

        This is performed by Alice
        """
        # Parse and validate input parameters
        offer_address = param_bytes20(request.form, 'offer_address')
        offer_amount = param_uint256(request.form, 'offer_amount')
        want_amount = param_uint256(request.form, 'want_amount')

        # TODO: validate contract addresses etc. and verify on-chain stuff

        exch_id = os.urandom(20).encode('hex')

        # Save exchange details
        # TODO: replace with class instance, `Exchange` ?
        self._exchanges[exch_id] = dict(
            offer_address=offer_address,
            offer_amount=offer_amount,
            want_amount=want_amount,
            proposals=dict(),
            chosen_proposal=None,

            # Temporary placeholders
            # TODO: replace with correct contracts
            offer_htlc_address=self._htlc_address.encode('hex'),
            want_htlc_address=self._htlc_address.encode('hex')
        )

        return jsonify(dict(
            id=exch_id,
            ok=1
        ))

    def exch_get(self, exch_id):
        """
        Retrieve details of exchange
        """
        exch = self._get_exch(exch_id)
        return jsonify(exch)

    def exch_propose(self, exch_id, secret_hashed):
        """
        Somebody who has what the offerer wants proves they have deposited it.
        They post proof as a proposal for exchange, this includes a hash of the secret they chose.
        They post details of this deposit as the proposal.

        This is performed by Bob
        """
        exch = self._get_exch(exch_id)

        if exch['chosen_proposal']:
            return api_abort("Proposal has already been chosen", code=409)

        # Hashed secret is the 'image', pre-image can be supplied to prove knowledge of secret
        if secret_hashed in exch['proposals']:
            return api_abort("Duplicate proposal secret", code=409)

        # TODO: verify either side of the exchange aren't the same

        expiry = param_uint256(request.form, 'expiry')
        depositor = param_bytes20(request.form, 'depositor')

        # TODO: verify details on-chain, expiry, depositor and secret must match

        # Verify expiry time is acceptable
        # XXX: should minimum expiry be left to the contract, or the coordinator?
        now = int(time.time())
        min_expiry = now + MINIMUM_EXPIRY_DURATION
        if expiry < min_expiry:
            return api_abort("Expiry too short")

        # Store proposal
        exch['proposals'][secret_hashed] = dict(
            secret_hashed=secret_hashed,
            expiry=expiry,
            depositor=depositor
        )

        return jsonify(dict(
            ok=1
        ))

    def exch_proposal_get(self, exch_id, secret_hashed):
        """
        Retrieve details for a specific exchange proposal
        """
        exch, proposal = self._get_proposal(exch_id, secret_hashed)
        return jsonify(proposal)

    def exch_confirm(self, exch_id, secret_hashed):
        """
        The initial offerer (A), who wants M of Y and offers N of X
        Upon seeing the proposal by B for M of Y, posts their confirmation of the specific deal
        The offerer (A) must have deposited N of X, using the same hash image that (B) proposed

        This is performed by Alice
        """
        exch, proposal = self._get_proposal(exch_id, secret_hashed)
        
        # XXX: one side of the expiry must be twice as long as the other to handle failure case
        # TODO: verify on-chain details match the proposal

        exch['chosen_proposal'] = secret_hashed

        return jsonify(dict(
            ok=1
        ))

    def exch_release(self, exch_id, secret_hashed):
        """
        Bob releases the secret to withdraw funds that Alice deposited.

        This is performed by Bob
        """
        exch, proposal = self._get_proposal(exch_id, secret_hashed)
        secret = param_bytes32(request.form, 'secret')
        # TODO: verify secret matches secret_hashed

        return jsonify(dict(
            ok=1
        ))

    def exch_finish(self, exch_id, secret_hashed):
        """
        Alice uses the secret revealed by Bob to withdraw the funds he deposited.

        This is performed by Alice.

        This completes the exchange.
        """
        exch, proposal = self._get_proposal(exch_id, secret_hashed)

        # XXX: technically a web API call isn't necessary for this step
        #      the API should monitor the state of both sides of the exchange
        #      and update the status / information automagically

        return jsonify(dict(
            ok=1
        ))


def main(htlc_address):
    """
    Development server for coordinator

    NOTE: not suitable for 'production'
    SEE: http://flask.pocoo.org/docs/1.0/deploying/#deployment
    """
    if len(htlc_address) != 20:
        htlc_address = scan_bin(htlc_address)
    print("HTLC address:", htlc_address.encode('hex'))

    coordinator = CoordinatorBlueprint(htlc_address)
    app = Flask(__name__)
    # app.debug = 1
    app.register_blueprint(coordinator, url_prefix='/htlc')

    # NOTE: Flask reloader is DAF, doesn't work well with packages *shakes-fists*
    app.run(use_reloader=False)
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1]))
