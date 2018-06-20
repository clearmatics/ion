## Copyright (c) 2018 Harry Roberts. All Rights Reserved.
## SPDX-License-Identifier: LGPL-3.0+

import sys

from flask import Flask, Blueprint, request, abort, jsonify, make_response
from werkzeug.routing import BaseConverter

from ..args import arg_bytes32, arg_bytes20, arg_uint256
from ..utils import scan_bin

from .manager import ExchangeManager, ExchangeError


#######################################################################
# TODO: move to ..webutils or something

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


def params_parse(data, params):
    """
    Performs param unmarshalling operations to return a dictionary
    params_parse({...}, {name: unmarshal})
    """
    out = dict()
    for key, val in params.items():
        out[key] = val(data, key)
    return out


#######################################################################


class CoordinatorBlueprint(Blueprint):
    """
    Provides a web API for coordinating cross-chain HTLC exchanges
    """

    def __init__(self, htlc_address, **kwa):
        Blueprint.__init__(self, 'htlc', __name__, **kwa)

        self._manager = ExchangeManager(htlc_address)

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
            depositor=param_bytes20
        ))

        try:
            exch, proposal = self._manager.propose(exch_id, secret_hashed, **params)
        except ExchangeError as ex:
            return api_abort(str(ex))

        # TODO: redirect to proposal URL? - or avoid another GET request...
        return jsonify(dict(
            ok=1
        ))

    def exch_proposal_get(self, exch_id, secret_hashed):
        """
        Retrieve details for a specific exchange proposal
        """
        try:
            exch, proposal = self._manager.get_proposal(exch_id, secret_hashed)
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
        try:
            self._manager.confirm(exch_id, secret_hashed)
        except ExchangeError as ex:
            return api_abort(str(ex))

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
        try:
            self._manager.finish(exch_id, secret_hashed)
        except ExchangeError as ex:
            return api_abort(str(ex))

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
