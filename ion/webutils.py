## Copyright (c) 2018 Harry Roberts. All Rights Reserved.
## SPDX-License-Identifier: LGPL-3.0+

from binascii import hexlify

from flask import jsonify, abort, make_response
from werkzeug.routing import BaseConverter

from .args import arg_bytes32, arg_bytes20, arg_uint256
from .utils import scan_bin


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
    return hexlify(param_filter_arg(the_dict, key, arg_bytes32)).decode('ascii')


def param_bytes20(the_dict, key):
    return hexlify(param_filter_arg(the_dict, key, arg_bytes20)).decode('ascii')


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
        return hexlify(scan_bin(value)).decode('ascii')


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
